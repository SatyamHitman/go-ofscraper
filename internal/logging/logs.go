// =============================================================================
// FILE: internal/logging/logs.go
// PURPOSE: High-level log operations used throughout the application. Provides
//          convenience functions for structured logging with consistent key
//          names and context propagation. Ports Python utils/logs/logs.py.
// =============================================================================

package logging

import (
	"context"
	"log/slog"
)

// ---------------------------------------------------------------------------
// Context-aware loggers
// ---------------------------------------------------------------------------

// FromCtx returns a logger enriched with any context-scoped attributes
// (e.g. model name, request ID). Falls back to the root logger when the
// context has no attached logger.
//
// Parameters:
//   - ctx: The request/operation context.
//
// Returns:
//   - A *slog.Logger with context attributes.
func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger); ok {
		return l
	}
	return Logger()
}

// WithLogger stores a logger in the context for downstream retrieval.
//
// Parameters:
//   - ctx: The parent context.
//   - l: The logger to attach.
//
// Returns:
//   - A derived context containing the logger.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// loggerCtxKey is the unexported context key for the slog.Logger.
type loggerCtxKey struct{}

// ---------------------------------------------------------------------------
// Convenience wrappers
// ---------------------------------------------------------------------------

// WithModel returns a logger annotated with the model/username being processed.
//
// Parameters:
//   - username: The OF model username.
//
// Returns:
//   - An annotated *slog.Logger.
func WithModel(username string) *slog.Logger {
	return Logger().With(slog.String("model", username))
}

// WithComponent returns a logger annotated with the subsystem name.
//
// Parameters:
//   - name: Component identifier (e.g. "download", "api", "auth").
//
// Returns:
//   - An annotated *slog.Logger.
func WithComponent(name string) *slog.Logger {
	return Logger().With(slog.String("component", name))
}

// WithRequestID returns a logger annotated with a request/trace identifier.
//
// Parameters:
//   - id: The request or correlation ID.
//
// Returns:
//   - An annotated *slog.Logger.
func WithRequestID(id string) *slog.Logger {
	return Logger().With(slog.String("request_id", id))
}

// ---------------------------------------------------------------------------
// Standard attribute keys
// ---------------------------------------------------------------------------

// Common attribute key constants for consistent log field naming.
const (
	KeyModel     = "model"
	KeyMediaID   = "media_id"
	KeyPostID    = "post_id"
	KeyMediaType = "media_type"
	KeyAction    = "action"
	KeyPath      = "path"
	KeyURL       = "url"
	KeyStatus    = "status"
	KeyDuration  = "duration"
	KeyBytes     = "bytes"
	KeyError     = "error"
	KeyRetry     = "retry"
	KeyWorker    = "worker"
)
