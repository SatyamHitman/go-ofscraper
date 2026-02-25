// =============================================================================
// FILE: internal/logging/level.go
// PURPOSE: Log level utilities and custom level definitions. Adds a TRACE
//          level below DEBUG for very verbose diagnostic output that is
//          normally suppressed. Ports Python utils/logs/utils/level.py.
// =============================================================================

package logging

import (
	"log/slog"
)

// ---------------------------------------------------------------------------
// Custom log levels
// ---------------------------------------------------------------------------

// LevelTrace is a custom level below DEBUG for extremely verbose output.
// slog uses integers for levels: DEBUG=-4, INFO=0, WARN=4, ERROR=8.
// TRACE sits at -8 (below DEBUG).
const LevelTrace slog.Level = -8

// LevelNames maps custom level values to human-readable names.
var LevelNames = map[slog.Level]string{
	LevelTrace: "TRACE",
}

// ---------------------------------------------------------------------------
// Level helpers
// ---------------------------------------------------------------------------

// LevelFromString parses a level name and returns the slog.Level. Supports
// the standard names (DEBUG, INFO, WARN, ERROR) plus our custom TRACE level.
//
// Parameters:
//   - name: The level name (case-insensitive).
//
// Returns:
//   - The corresponding slog.Level. Defaults to INFO for unrecognised names.
func LevelFromString(name string) slog.Level {
	// Check custom levels first.
	switch name {
	case "TRACE", "trace":
		return LevelTrace
	}

	// Fall back to slog's built-in parsing.
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(name)); err != nil {
		return slog.LevelInfo
	}
	return lvl
}

// LevelToString returns a human-readable name for the given level.
//
// Parameters:
//   - level: The slog.Level to convert.
//
// Returns:
//   - A string name (e.g. "TRACE", "DEBUG", "INFO", "WARN", "ERROR").
func LevelToString(level slog.Level) string {
	if name, ok := LevelNames[level]; ok {
		return name
	}
	return level.String()
}

// IsTrace reports whether the current logging level permits TRACE output.
// This is useful for guarding expensive debug-formatting calls.
//
// Returns:
//   - true if TRACE-level messages would be emitted.
func IsTrace() bool {
	return Logger().Enabled(nil, LevelTrace)
}

// IsDebug reports whether the current logging level permits DEBUG output.
//
// Returns:
//   - true if DEBUG-level messages would be emitted.
func IsDebug() bool {
	return Logger().Enabled(nil, slog.LevelDebug)
}
