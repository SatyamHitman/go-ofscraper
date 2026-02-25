// =============================================================================
// FILE: internal/app/daemon.go
// PURPOSE: Daemon mode for running the scraper on a schedule using a ticker.
//          Manages DaemonConfig with configurable intervals and runs the
//          scrape loop until the context is cancelled. Ports Python
//          runner/manager/daemon.py.
// =============================================================================

package app

import (
	"context"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// DaemonConfig
// ---------------------------------------------------------------------------

// DaemonConfig holds the configuration for daemon mode operation.
type DaemonConfig struct {
	// Interval is the duration between scrape runs.
	Interval time.Duration

	// Actions are the actions to perform on each run.
	Actions []string

	// Areas are the content areas to scrape on each run.
	Areas []string
}

// DefaultDaemonConfig returns a DaemonConfig with sensible defaults.
//
// Returns:
//   - A DaemonConfig with a 6-hour interval.
func DefaultDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Interval: 6 * time.Hour,
		Actions:  []string{"download"},
		Areas:    []string{"timeline", "messages", "stories"},
	}
}

// Validate checks that the daemon configuration is valid.
//
// Returns:
//   - Error if the configuration is invalid.
func (dc DaemonConfig) Validate() error {
	if dc.Interval < 1*time.Minute {
		return fmt.Errorf("daemon interval must be at least 1 minute, got %s", dc.Interval)
	}
	if len(dc.Actions) == 0 {
		return fmt.Errorf("daemon requires at least one action")
	}
	if len(dc.Areas) == 0 {
		return fmt.Errorf("daemon requires at least one content area")
	}
	return nil
}

// ---------------------------------------------------------------------------
// RunDaemon
// ---------------------------------------------------------------------------

// RunDaemon starts the daemon loop that runs scrape operations at the
// configured interval. Blocks until the context is cancelled.
//
// Parameters:
//   - ctx: Context for cancellation (e.g., from signal handler).
//   - dc: The daemon configuration.
//
// Returns:
//   - Error if the daemon fails to start or encounters a fatal error.
func (a *App) RunDaemon(ctx context.Context, dc DaemonConfig) error {
	if err := dc.Validate(); err != nil {
		return fmt.Errorf("invalid daemon config: %w", err)
	}

	a.logger.Info("daemon mode starting",
		"interval", dc.Interval,
		"actions", dc.Actions,
		"areas", dc.Areas,
	)

	// Run immediately on first iteration.
	if err := a.runDaemonCycle(ctx, dc); err != nil {
		a.logger.Error("daemon cycle failed", "error", err)
	}

	ticker := time.NewTicker(dc.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("daemon shutting down")
			return nil
		case <-ticker.C:
			a.logger.Info("daemon cycle starting")
			if err := a.runDaemonCycle(ctx, dc); err != nil {
				a.logger.Error("daemon cycle failed", "error", err)
				// Continue running; do not exit on non-fatal errors.
			}
			a.logger.Info(fmt.Sprintf("daemon sleeping for %s until next run", dc.Interval))
		}
	}
}

// runDaemonCycle executes a single scrape cycle within the daemon loop.
func (a *App) runDaemonCycle(ctx context.Context, dc DaemonConfig) error {
	// Dispatch to the action router for each configured action.
	for _, action := range dc.Actions {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := a.RunAction(ctx, action, dc.Areas, nil); err != nil {
			return fmt.Errorf("daemon action %s: %w", action, err)
		}
	}
	return nil
}
