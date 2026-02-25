// =============================================================================
// FILE: internal/http/ratelimit.go
// PURPOSE: Token bucket rate limiter for API requests. Prevents exceeding
//          the OF API rate limits by spacing out requests. Ports the rate
//          limit logic from Python sessionmanager.
// =============================================================================

package http

import (
	"context"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Rate limiter
// ---------------------------------------------------------------------------

// RateLimiter implements a simple token bucket for request throttling.
type RateLimiter struct {
	enabled  bool
	interval time.Duration
	lastReq  time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a rate limiter.
//
// Parameters:
//   - enabled: Whether rate limiting is active.
//
// Returns:
//   - A configured RateLimiter.
func NewRateLimiter(enabled bool) *RateLimiter {
	return &RateLimiter{
		enabled:  enabled,
		interval: 200 * time.Millisecond, // ~5 req/s default
	}
}

// Wait blocks until a request slot is available.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - Error if context is cancelled while waiting.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	if !rl.enabled {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastReq)

	if elapsed < rl.interval {
		wait := rl.interval - elapsed
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	rl.lastReq = time.Now()
	return nil
}

// SetInterval changes the minimum interval between requests.
//
// Parameters:
//   - d: The new minimum interval.
func (rl *RateLimiter) SetInterval(d time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.interval = d
}
