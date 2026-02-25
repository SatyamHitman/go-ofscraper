// =============================================================================
// FILE: internal/commands/scraper/schedule.go
// PURPOSE: Schedule configuration and helpers for daemon mode. Defines
//          scheduling intervals, next-run computation, and schedule validation.
//          Ports Python runner/manager/daemon.py scheduling logic.
// =============================================================================

package scraper

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// ScheduleConfig
// ---------------------------------------------------------------------------

// ScheduleConfig holds the scheduling parameters for daemon mode runs.
type ScheduleConfig struct {
	// Interval is the duration between scrape runs.
	Interval time.Duration

	// StartTime optionally specifies a fixed time-of-day for the first run
	// (e.g., "03:00" for 3 AM). Empty means start immediately.
	StartTime string

	// MaxRuns limits the total number of runs (0 = unlimited).
	MaxRuns int
}

// DefaultSchedule returns a ScheduleConfig with sensible defaults.
//
// Returns:
//   - A ScheduleConfig with a 6-hour interval and no run limit.
func DefaultSchedule() ScheduleConfig {
	return ScheduleConfig{
		Interval: 6 * time.Hour,
		MaxRuns:  0,
	}
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

// Validate checks that the schedule configuration is valid.
//
// Returns:
//   - Error if the configuration is invalid.
func (sc ScheduleConfig) Validate() error {
	if sc.Interval < 1*time.Minute {
		return fmt.Errorf("schedule interval must be at least 1 minute, got %s", sc.Interval)
	}
	if sc.StartTime != "" {
		if _, err := parseTimeOfDay(sc.StartTime); err != nil {
			return fmt.Errorf("invalid start_time %q: %w", sc.StartTime, err)
		}
	}
	if sc.MaxRuns < 0 {
		return fmt.Errorf("max_runs must be non-negative, got %d", sc.MaxRuns)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Next run computation
// ---------------------------------------------------------------------------

// NextRun calculates the time of the next scheduled run.
//
// Parameters:
//   - lastRun: The time the last run started. Zero value means no previous run.
//
// Returns:
//   - The time.Time of the next run.
func (sc ScheduleConfig) NextRun(lastRun time.Time) time.Time {
	if lastRun.IsZero() && sc.StartTime != "" {
		return nextTimeOfDay(sc.StartTime)
	}
	if lastRun.IsZero() {
		return time.Now()
	}
	return lastRun.Add(sc.Interval)
}

// SleepDuration calculates how long to sleep before the next run.
//
// Parameters:
//   - lastRun: The time the last run started.
//
// Returns:
//   - The duration to sleep.
func (sc ScheduleConfig) SleepDuration(lastRun time.Time) time.Duration {
	next := sc.NextRun(lastRun)
	d := time.Until(next)
	if d < 0 {
		return 0
	}
	return d
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseTimeOfDay parses a "HH:MM" string into hour and minute.
func parseTimeOfDay(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected HH:MM format: %w", err)
	}
	return t, nil
}

// nextTimeOfDay returns the next occurrence of the given time-of-day.
func nextTimeOfDay(s string) time.Time {
	parsed, err := parseTimeOfDay(s)
	if err != nil {
		return time.Now()
	}
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
