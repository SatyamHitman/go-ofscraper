// =============================================================================
// FILE: internal/logging/discord.go
// PURPOSE: Discord webhook log handler. Sends log messages to a Discord
//          channel via webhook. Buffers messages to avoid rate limits and
//          batches them on a timer. Ports Python
//          utils/logs/classes/handlers/discord.py.
// =============================================================================

package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Discord handler
// ---------------------------------------------------------------------------

// discordHandler sends log records to a Discord webhook URL. Messages are
// buffered and flushed periodically to respect Discord rate limits.
type discordHandler struct {
	webhookURL string
	level      slog.Level
	attrs      []slog.Attr
	group      string

	// Buffer for batching messages.
	mu      sync.Mutex
	buf     []string
	timer   *time.Timer
	stopped bool
}

// flushInterval is how often buffered messages are sent to Discord.
const flushInterval = 5 * time.Second

// maxDiscordMessageLen is the maximum content length Discord accepts.
const maxDiscordMessageLen = 2000

// newDiscordHandler creates a handler that posts log messages to the given
// Discord webhook URL.
//
// Parameters:
//   - webhookURL: The Discord webhook endpoint.
//   - level: Minimum log level to send.
//
// Returns:
//   - A slog.Handler that buffers and sends to Discord.
func newDiscordHandler(webhookURL string, level slog.Level) slog.Handler {
	dh := &discordHandler{
		webhookURL: webhookURL,
		level:      level,
	}

	// Start the flush timer.
	dh.timer = time.AfterFunc(flushInterval, dh.flush)

	// Register for shutdown cleanup.
	registerCloseable(dh)

	return dh
}

// Enabled reports whether the handler handles records at the given level.
func (dh *discordHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= dh.level
}

// Handle buffers the log record for later dispatch to Discord.
func (dh *discordHandler) Handle(_ context.Context, record slog.Record) error {
	dh.mu.Lock()
	defer dh.mu.Unlock()

	if dh.stopped {
		return nil
	}

	// Format: [LEVEL] message key=value ...
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", record.Level.String(), record.Message))

	// Append static attrs.
	for _, a := range dh.attrs {
		sb.WriteString(fmt.Sprintf(" %s=%v", a.Key, a.Value.Any()))
	}

	// Append record attrs.
	record.Attrs(func(a slog.Attr) bool {
		sb.WriteString(fmt.Sprintf(" %s=%v", a.Key, a.Value.Any()))
		return true
	})

	dh.buf = append(dh.buf, sb.String())
	return nil
}

// WithAttrs returns a new handler with additional attributes.
func (dh *discordHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &discordHandler{
		webhookURL: dh.webhookURL,
		level:      dh.level,
		attrs:      append(dh.attrs, attrs...),
		group:      dh.group,
		buf:        dh.buf,
		timer:      dh.timer,
	}
}

// WithGroup returns a new handler with a group name prefix.
func (dh *discordHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if dh.group != "" {
		newGroup = dh.group + "." + name
	}
	return &discordHandler{
		webhookURL: dh.webhookURL,
		level:      dh.level,
		attrs:      dh.attrs,
		group:      newGroup,
		buf:        dh.buf,
		timer:      dh.timer,
	}
}

// Close flushes remaining messages and stops the timer.
//
// Returns:
//   - Error from the final flush, if any.
func (dh *discordHandler) Close() error {
	dh.mu.Lock()
	dh.stopped = true
	dh.mu.Unlock()

	if dh.timer != nil {
		dh.timer.Stop()
	}

	dh.flush()
	return nil
}

// flush sends all buffered messages to Discord. Called by the timer and
// during shutdown.
func (dh *discordHandler) flush() {
	dh.mu.Lock()
	if len(dh.buf) == 0 {
		// Reset timer if not stopped.
		if !dh.stopped {
			dh.timer.Reset(flushInterval)
		}
		dh.mu.Unlock()
		return
	}

	// Drain buffer.
	messages := dh.buf
	dh.buf = nil

	if !dh.stopped {
		dh.timer.Reset(flushInterval)
	}
	dh.mu.Unlock()

	// Join messages and chunk to respect Discord limits.
	combined := strings.Join(messages, "\n")
	chunks := chunkString(combined, maxDiscordMessageLen)

	for _, chunk := range chunks {
		_ = postToDiscord(dh.webhookURL, chunk)
	}
}

// ---------------------------------------------------------------------------
// Discord HTTP helpers
// ---------------------------------------------------------------------------

// discordPayload is the JSON body for a Discord webhook message.
type discordPayload struct {
	Content string `json:"content"`
}

// postToDiscord sends a single message to the webhook URL.
//
// Parameters:
//   - url: The Discord webhook URL.
//   - content: The message text.
//
// Returns:
//   - Error if the HTTP request fails.
func postToDiscord(url, content string) error {
	payload := discordPayload{Content: "```\n" + content + "\n```"}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}
	return nil
}

// chunkString splits a string into pieces of at most maxLen bytes.
//
// Parameters:
//   - s: The string to split.
//   - maxLen: Maximum length of each chunk.
//
// Returns:
//   - Slice of string chunks.
func chunkString(s string, maxLen int) []string {
	if len(s) <= maxLen {
		return []string{s}
	}

	var chunks []string
	for len(s) > 0 {
		end := maxLen
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[:end])
		s = s[end:]
	}
	return chunks
}
