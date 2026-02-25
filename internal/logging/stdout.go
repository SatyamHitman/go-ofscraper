// =============================================================================
// FILE: internal/logging/stdout.go
// PURPOSE: Stdout log handler with optional coloured output using the tint
//          library for slog. Renders human-friendly log lines to the terminal.
//          Ports Python utils/logs/stdout.py.
// =============================================================================

package logging

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// ---------------------------------------------------------------------------
// Stdout handler
// ---------------------------------------------------------------------------

// newStdoutHandler creates a slog.Handler that writes human-readable,
// optionally coloured log lines to the given writer (typically os.Stdout).
//
// Parameters:
//   - w: The output writer.
//   - level: Minimum log level to emit.
//   - color: Whether to enable ANSI colour codes.
//
// Returns:
//   - A configured slog.Handler.
func newStdoutHandler(w io.Writer, level slog.Level, color bool) slog.Handler {
	// When colour is enabled, use tint for pretty terminal output.
	if color && isTerminal(w) {
		return tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
			NoColor:    false,
		})
	}

	// Fallback to plain text handler.
	return slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})
}

// isTerminal checks if the writer is connected to a terminal (for colour
// support detection). Only os.Stdout and os.Stderr are considered terminals.
//
// Parameters:
//   - w: The writer to check.
//
// Returns:
//   - true if the writer is a known terminal file descriptor.
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		// Character device = terminal on most platforms.
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}
