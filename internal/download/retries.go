// =============================================================================
// FILE: internal/download/retries.go
// PURPOSE: Download retry policy. Defines retry behavior specific to download
//          operations with configurable attempts and backoff.
//          Ports Python utils/retries.py.
// =============================================================================

package download

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// ---------------------------------------------------------------------------
// Retry config
// ---------------------------------------------------------------------------

// RetryPolicy configures download retry behavior.
type RetryPolicy struct {
	MaxAttempts int           // Maximum number of attempts
	InitialWait time.Duration // Initial wait between retries
	MaxWait     time.Duration // Maximum wait between retries
	Multiplier  float64       // Backoff multiplier
}

// DefaultRetryPolicy returns the default download retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		InitialWait: 2 * time.Second,
		MaxWait:     30 * time.Second,
		Multiplier:  2.0,
	}
}

// ---------------------------------------------------------------------------
// Retry execution
// ---------------------------------------------------------------------------

// WithRetry executes a function with retry logic.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - policy: The retry policy to use.
//   - logger: Optional logger for retry events.
//   - fn: The function to execute.
//
// Returns:
//   - Error from the last attempt, or nil on success.
func WithRetry(ctx context.Context, policy RetryPolicy, logger *slog.Logger, fn func() error) error {
	wait := policy.InitialWait

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt == policy.MaxAttempts {
			return fmt.Errorf("failed after %d attempts: %w", policy.MaxAttempts, err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if logger != nil {
			logger.Warn("download retry",
				"attempt", attempt,
				"max", policy.MaxAttempts,
				"wait", wait,
				"error", err,
			)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}

		// Exponential backoff.
		wait = time.Duration(float64(wait) * policy.Multiplier)
		if wait > policy.MaxWait {
			wait = policy.MaxWait
		}
	}

	return fmt.Errorf("retry exhausted")
}
