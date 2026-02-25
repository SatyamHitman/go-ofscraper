// =============================================================================
// FILE: internal/logging/sensitive.go
// PURPOSE: Sensitive-data redaction handler. Wraps another slog.Handler and
//          scrubs known sensitive fields (tokens, cookies, passwords) from log
//          records before forwarding them. Ports Python
//          utils/logs/utils/sensitive.py.
// =============================================================================

package logging

import (
	"context"
	"log/slog"
	"strings"
)

// ---------------------------------------------------------------------------
// Sensitive field redaction
// ---------------------------------------------------------------------------

// sensitiveKeys lists attribute key substrings that indicate sensitive data
// which must be redacted in logs.
var sensitiveKeys = []string{
	"token",
	"cookie",
	"password",
	"secret",
	"auth",
	"x-bc",
	"sess",
	"key",
	"sign",
	"credential",
}

// redactedValue is the replacement string for sensitive data.
const redactedValue = "[REDACTED]"

// sensitiveHandler wraps a slog.Handler and redacts sensitive attributes.
type sensitiveHandler struct {
	inner slog.Handler
}

// newSensitiveHandler creates a handler that redacts sensitive fields before
// forwarding records to the inner handler.
//
// Parameters:
//   - inner: The handler to forward sanitised records to.
//
// Returns:
//   - A wrapping slog.Handler that performs redaction.
func newSensitiveHandler(inner slog.Handler) slog.Handler {
	return &sensitiveHandler{inner: inner}
}

// Enabled delegates to the inner handler.
func (sh *sensitiveHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return sh.inner.Enabled(ctx, level)
}

// Handle redacts sensitive attrs in the record before forwarding.
func (sh *sensitiveHandler) Handle(ctx context.Context, record slog.Record) error {
	// Clone the record so we don't modify the original.
	clean := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)

	record.Attrs(func(a slog.Attr) bool {
		clean.AddAttrs(redactAttr(a))
		return true
	})

	return sh.inner.Handle(ctx, clean)
}

// WithAttrs redacts static attrs and forwards.
func (sh *sensitiveHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	redacted := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		redacted[i] = redactAttr(a)
	}
	return &sensitiveHandler{inner: sh.inner.WithAttrs(redacted)}
}

// WithGroup delegates to the inner handler.
func (sh *sensitiveHandler) WithGroup(name string) slog.Handler {
	return &sensitiveHandler{inner: sh.inner.WithGroup(name)}
}

// ---------------------------------------------------------------------------
// Redaction logic
// ---------------------------------------------------------------------------

// redactAttr returns the attribute with its value replaced if the key matches
// any sensitive pattern.
//
// Parameters:
//   - a: The slog attribute to check.
//
// Returns:
//   - The original attribute, or one with a redacted value.
func redactAttr(a slog.Attr) slog.Attr {
	keyLower := strings.ToLower(a.Key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(keyLower, sensitive) {
			return slog.String(a.Key, redactedValue)
		}
	}
	return a
}

// isSensitive checks if a key name matches any sensitive pattern.
//
// Parameters:
//   - key: The attribute key to test.
//
// Returns:
//   - true if the key appears sensitive.
func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
