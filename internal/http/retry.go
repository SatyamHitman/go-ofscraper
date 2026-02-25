// =============================================================================
// FILE: internal/http/retry.go
// PURPOSE: HTTP retry logic with exponential backoff. Automatically retries
//          failed requests on transient errors (timeouts, 5xx, connection
//          resets). Ports retry logic from Python sessionmanager.
// =============================================================================

package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Retry configuration
// ---------------------------------------------------------------------------

// RetryConfig controls retry behaviour for HTTP requests.
type RetryConfig struct {
	MaxAttempts int           // Maximum total attempts (including first try)
	MinWait     time.Duration // Minimum wait between retries
	MaxWait     time.Duration // Maximum wait between retries
	Multiplier  float64       // Backoff multiplier (e.g. 2.0 for doubling)
}

// DefaultRetryConfig returns sensible defaults for API requests.
//
// Returns:
//   - A RetryConfig with standard settings.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 5,
		MinWait:     1 * time.Second,
		MaxWait:     30 * time.Second,
		Multiplier:  2.0,
	}
}

// ---------------------------------------------------------------------------
// Retry execution
// ---------------------------------------------------------------------------

// DoWithRetry executes an HTTP request with automatic retry on transient errors.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - sm: The session manager to use.
//   - req: The request to execute.
//   - cfg: Retry configuration.
//
// Returns:
//   - The successful Response, and any error after all retries exhausted.
func DoWithRetry(ctx context.Context, sm *SessionManager, req *Request, cfg RetryConfig) (*Response, error) {
	var lastErr error
	wait := cfg.MinWait

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		resp, err := sm.Do(ctx, req)

		// Success.
		if err == nil && resp.IsOK() {
			return resp, nil
		}

		// Non-retryable errors.
		if err == nil && (resp.IsAuthError() || resp.IsNotFound()) {
			return resp, nil // Return as-is for caller to handle.
		}

		// Determine if retryable.
		if err == nil && resp.IsRateLimit() {
			lastErr = fmt.Errorf("rate limited (status %d)", resp.StatusCode)
			resp.Close()
		} else if err != nil && isRetryableError(err) {
			lastErr = err
		} else if err == nil && resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error (status %d)", resp.StatusCode)
			resp.Close()
		} else {
			// Not retryable.
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		if attempt < cfg.MaxAttempts {
			slog.Debug("retrying request",
				"url", req.URL,
				"attempt", attempt,
				"wait", wait,
				"error", lastErr,
			)

			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return nil, ctx.Err()
			}

			// Exponential backoff.
			wait = time.Duration(float64(wait) * cfg.Multiplier)
			if wait > cfg.MaxWait {
				wait = cfg.MaxWait
			}
		}
	}

	return nil, fmt.Errorf("max retries (%d) exceeded: %w", cfg.MaxAttempts, lastErr)
}

// isRetryableError checks if a network error is transient.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Context cancellation is not retryable.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Network errors are generally retryable.
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	// Connection reset/refused.
	msg := err.Error()
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "EOF")
}
