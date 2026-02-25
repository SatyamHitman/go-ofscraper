// =============================================================================
// FILE: internal/logging/text_handler.go
// PURPOSE: Custom text handler for slog that produces compact, single-line
//          output suitable for progress-bar-style TUI overlays where
//          vertical space is limited. Ports Python
//          utils/logs/classes/handlers/text.py.
// =============================================================================

package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Text handler
// ---------------------------------------------------------------------------

// textHandler writes minimal single-line log entries to a writer. Designed
// for embedding in TUI panels where multi-line JSON or text output would
// disrupt the layout.
type textHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
	group string
	mu    sync.Mutex
}

// NewTextHandler creates a compact text log handler.
//
// Parameters:
//   - w: The output writer.
//   - level: Minimum log level.
//
// Returns:
//   - A slog.Handler that writes compact text lines.
func NewTextHandler(w io.Writer, level slog.Level) slog.Handler {
	return &textHandler{
		w:     w,
		level: level,
	}
}

// Enabled implements slog.Handler.
func (th *textHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= th.level
}

// Handle implements slog.Handler. Writes a single line in the format:
// HH:MM:SS LEVEL message key=value ...
func (th *textHandler) Handle(_ context.Context, record slog.Record) error {
	th.mu.Lock()
	defer th.mu.Unlock()

	var sb strings.Builder
	sb.WriteString(record.Time.Format(time.TimeOnly))
	sb.WriteByte(' ')

	// Short level tag.
	switch {
	case record.Level >= slog.LevelError:
		sb.WriteString("ERR ")
	case record.Level >= slog.LevelWarn:
		sb.WriteString("WRN ")
	case record.Level >= slog.LevelInfo:
		sb.WriteString("INF ")
	default:
		sb.WriteString("DBG ")
	}

	sb.WriteString(record.Message)

	// Append static attrs with optional group prefix.
	prefix := ""
	if th.group != "" {
		prefix = th.group + "."
	}
	for _, a := range th.attrs {
		sb.WriteString(fmt.Sprintf(" %s%s=%v", prefix, a.Key, a.Value.Any()))
	}

	// Append record attrs.
	record.Attrs(func(a slog.Attr) bool {
		sb.WriteString(fmt.Sprintf(" %s%s=%v", prefix, a.Key, a.Value.Any()))
		return true
	})

	sb.WriteByte('\n')

	_, err := io.WriteString(th.w, sb.String())
	return err
}

// WithAttrs implements slog.Handler.
func (th *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &textHandler{
		w:     th.w,
		level: th.level,
		attrs: append(th.attrs, attrs...),
		group: th.group,
	}
}

// WithGroup implements slog.Handler.
func (th *textHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if th.group != "" {
		newGroup = th.group + "." + name
	}
	return &textHandler{
		w:     th.w,
		level: th.level,
		attrs: th.attrs,
		group: newGroup,
	}
}
