// =============================================================================
// FILE: internal/logging/trace.go
// PURPOSE: Trace-level logging helpers. Provides convenience functions for
//          emitting TRACE-level messages that are only recorded when the
//          log level is set sufficiently low. Ports Python
//          utils/logs/utils/trace.py.
// =============================================================================

package logging

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// ---------------------------------------------------------------------------
// Trace logging
// ---------------------------------------------------------------------------

// Trace emits a TRACE-level log message through the root logger.
// No-ops if the current level is above TRACE.
//
// Parameters:
//   - msg: The log message.
//   - args: Optional key-value pairs for structured attributes.
func Trace(msg string, args ...any) {
	if !IsTrace() {
		return
	}
	logAtLevel(context.Background(), LevelTrace, msg, args...)
}

// TraceCtx emits a TRACE-level log message using the context's logger.
//
// Parameters:
//   - ctx: The context (may contain a logger).
//   - msg: The log message.
//   - args: Optional key-value pairs.
func TraceCtx(ctx context.Context, msg string, args ...any) {
	l := FromCtx(ctx)
	if !l.Enabled(ctx, LevelTrace) {
		return
	}
	logAtLevel(ctx, LevelTrace, msg, args...)
}

// TraceFunc emits a TRACE-level message produced by a function. The function
// is only called if TRACE is enabled, avoiding expensive formatting when the
// level is suppressed.
//
// Parameters:
//   - fn: A function that returns the message and attrs.
func TraceFunc(fn func() (string, []any)) {
	if !IsTrace() {
		return
	}
	msg, args := fn()
	logAtLevel(context.Background(), LevelTrace, msg, args...)
}

// ---------------------------------------------------------------------------
// Internal helper
// ---------------------------------------------------------------------------

// logAtLevel emits a log record at the specified level, capturing the correct
// caller location (two frames up from the public Trace* function).
func logAtLevel(ctx context.Context, level slog.Level, msg string, args ...any) {
	l := Logger()

	// Build a record with the correct source location.
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip: Callers, logAtLevel, Trace*

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])

	// Add args as attrs.
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		record.AddAttrs(slog.Any(key, args[i+1]))
	}

	_ = l.Handler().Handle(ctx, record)
}
