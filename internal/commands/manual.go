// =============================================================================
// FILE: internal/commands/manual.go
// PURPOSE: Manual download command. Handles URL-based downloads where the user
//          provides direct media URLs as CLI arguments. Ports Python
//          runner/manual/manual.py.
// =============================================================================

package commands

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"gofscraper/internal/app"
	cmdutils "gofscraper/internal/commands/utils"
)

// ---------------------------------------------------------------------------
// ManualCommand
// ---------------------------------------------------------------------------

// ManualCommand handles direct URL-based media downloads.
type ManualCommand struct {
	cmdutils.CommandBase
}

// NewManualCommand creates a new ManualCommand.
//
// Parameters:
//   - logger: Structured logger for output.
//
// Returns:
//   - A configured ManualCommand.
func NewManualCommand(logger *slog.Logger) *ManualCommand {
	return &ManualCommand{
		CommandBase: cmdutils.NewCommandBase(logger),
	}
}

// Name returns the command name.
func (m *ManualCommand) Name() string {
	return "manual"
}

// Run executes the manual download command.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance providing session and config.
//   - urls: The URLs to download.
//
// Returns:
//   - Error if the download process fails.
func (m *ManualCommand) Run(ctx context.Context, a *app.App, urls []string) error {
	m.LogStart(m.Name(), urls)
	defer m.LogDone(m.Name())

	if len(urls) == 0 {
		return fmt.Errorf("no URLs provided for manual download")
	}

	// Validate and parse URLs.
	parsed, err := m.validateURLs(urls)
	if err != nil {
		return err
	}

	m.Logger.Info("starting manual download", "url_count", len(parsed))

	// Process each URL.
	var succeeded, failed int
	for i, u := range parsed {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		m.Logger.Info(fmt.Sprintf("downloading %d/%d: %s", i+1, len(parsed), u.String()))

		if err := m.downloadURL(ctx, a, u); err != nil {
			m.Logger.Error("download failed", "url", u.String(), "error", err)
			failed++
			continue
		}
		succeeded++
	}

	m.Logger.Info("manual download complete",
		"succeeded", succeeded,
		"failed", failed,
		"total", len(parsed),
	)

	return nil
}

// validateURLs parses and validates the input URL strings.
func (m *ManualCommand) validateURLs(urls []string) ([]*url.URL, error) {
	var parsed []*url.URL
	var invalid []string

	for _, raw := range urls {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		u, err := url.Parse(raw)
		if err != nil || u.Scheme == "" || u.Host == "" {
			invalid = append(invalid, raw)
			continue
		}
		parsed = append(parsed, u)
	}

	if len(invalid) > 0 {
		m.Logger.Warn("skipping invalid URLs", "urls", invalid)
	}

	if len(parsed) == 0 {
		return nil, fmt.Errorf("no valid URLs provided")
	}

	return parsed, nil
}

// downloadURL downloads a single URL.
func (m *ManualCommand) downloadURL(ctx context.Context, a *app.App, u *url.URL) error {
	_ = a.Session() // Will be used when download is wired.
	_ = ctx

	// TODO: Wire to download handler - resolve filename, download content.
	m.Logger.Debug("download placeholder", "url", u.String())
	return nil
}
