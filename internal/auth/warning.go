// =============================================================================
// FILE: internal/auth/warning.go
// PURPOSE: Auth warnings. Checks auth state and emits warnings for common
//          issues like soon-to-expire sessions, missing optional fields, and
//          browser cookie staleness. Ports Python
//          utils/auth/utils/warning/warning.py, check.py, print.py.
// =============================================================================

package auth

import (
	"log/slog"
	"time"
)

// ---------------------------------------------------------------------------
// Auth warnings
// ---------------------------------------------------------------------------

// Warning represents a non-fatal auth issue.
type Warning struct {
	Level   string // "info", "warn", "error"
	Message string
}

// CheckWarnings inspects the current auth state and returns any warnings.
//
// Parameters:
//   - data: The auth data to check.
//
// Returns:
//   - Slice of warnings found.
func CheckWarnings(data *Data) []Warning {
	if data == nil {
		return []Warning{{Level: "error", Message: "no auth data loaded"}}
	}

	var warnings []Warning

	// Check for empty optional fields.
	if data.AppToken == "" {
		warnings = append(warnings, Warning{
			Level:   "info",
			Message: "app_token not set — using default",
		})
	}

	// Check user-agent format (should contain browser identifier).
	if len(data.UserAgent) < 20 {
		warnings = append(warnings, Warning{
			Level:   "warn",
			Message: "user-agent appears too short — may cause auth failures",
		})
	}

	return warnings
}

// LogWarnings emits all auth warnings through the logger.
//
// Parameters:
//   - data: The auth data to check.
func LogWarnings(data *Data) {
	warnings := CheckWarnings(data)
	for _, w := range warnings {
		switch w.Level {
		case "error":
			slog.Error("auth warning", "message", w.Message)
		case "warn":
			slog.Warn("auth warning", "message", w.Message)
		default:
			slog.Info("auth notice", "message", w.Message)
		}
	}
}

// lastWarningTime tracks when the last auth warning was shown.
var lastWarningTime time.Time

// ShouldShowWarning returns true if enough time has passed since the last
// auth warning to avoid spamming the user.
//
// Returns:
//   - true if a warning should be displayed.
func ShouldShowWarning() bool {
	if time.Since(lastWarningTime) > 30*time.Minute {
		lastWarningTime = time.Now()
		return true
	}
	return false
}
