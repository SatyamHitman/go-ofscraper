// =============================================================================
// FILE: internal/logging/other.go
// PURPOSE: Multi-handler and miscellaneous logging utilities. The multiHandler
//          fans out log records to multiple slog.Handlers (stdout + file +
//          discord simultaneously). Ports Python utils/logs/other.py.
// =============================================================================

package logging

import (
	"context"
	"log/slog"
)

// ---------------------------------------------------------------------------
// Multi-handler
// ---------------------------------------------------------------------------

// multiHandler fans out each log record to multiple underlying handlers.
type multiHandler struct {
	handlers []slog.Handler
}

// newMultiHandler creates a handler that forwards records to all given handlers.
//
// Parameters:
//   - handlers: The handlers to fan out to.
//
// Returns:
//   - A combined slog.Handler.
func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

// Enabled returns true if any underlying handler is enabled for the level.
func (mh *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle sends the record to all enabled handlers. Returns the first error.
func (mh *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, h := range mh.handlers {
		if h.Enabled(ctx, record.Level) {
			if err := h.Handle(ctx, record); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// WithAttrs returns a new multi-handler where each underlying handler has
// the additional attributes.
func (mh *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

// WithGroup returns a new multi-handler where each underlying handler has
// the group applied.
func (mh *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
