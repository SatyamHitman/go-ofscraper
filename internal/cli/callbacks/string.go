// =============================================================================
// FILE: internal/cli/callbacks/string.go
// PURPOSE: String transformation callbacks for flag values.
// =============================================================================

package callbacks

import (
	"strings"
)

// TrimAndLower trims whitespace and lowercases the string.
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// SplitCSV splits a comma-separated string into a trimmed slice.
// Returns nil for an empty input.
func SplitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Truncate truncates a string to maxLen characters, appending an ellipsis if
// it was shortened.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
