// =============================================================================
// FILE: internal/http/sleeper.go
// PURPOSE: AdaptiveSleeper for handling 429 and 403 responses. Increases
//          sleep duration on rate limits and slowly decays back to zero over
//          time. Ports Python managers/sessionmanager/sleepers.py.
// =============================================================================

package http

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Adaptive sleeper
// ---------------------------------------------------------------------------

// AdaptiveSleeper dynamically adjusts sleep time based on API responses.
// Increases on rate limits/forbidden, decays over time back to zero.
type AdaptiveSleeper struct {
	currentSleep time.Duration
	maxSleep     time.Duration
	minSleep     time.Duration
	increment    time.Duration
	lastBump     time.Time
	decayAfter   time.Duration
	mu           sync.Mutex
}

// NewAdaptiveSleeper creates a sleeper with sensible defaults.
//
// Returns:
//   - A configured AdaptiveSleeper.
func NewAdaptiveSleeper() *AdaptiveSleeper {
	return &AdaptiveSleeper{
		currentSleep: 0,
		maxSleep:     30 * time.Second,
		minSleep:     0,
		increment:    2 * time.Second,
		decayAfter:   60 * time.Second,
	}
}

// Sleep blocks for the current adaptive duration. If enough time has passed
// since the last bump, the sleep duration decays.
//
// Parameters:
//   - ctx: Context for cancellation.
func (as *AdaptiveSleeper) Sleep(ctx context.Context) {
	as.mu.Lock()
	defer as.mu.Unlock()

	// Decay sleep if no recent bumps.
	as.decay()

	if as.currentSleep <= 0 {
		return
	}

	slog.Debug("adaptive sleep", "duration", as.currentSleep)

	select {
	case <-time.After(as.currentSleep):
	case <-ctx.Done():
	}
}

// OnRateLimit increases sleep duration in response to a 429/504.
func (as *AdaptiveSleeper) OnRateLimit() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.currentSleep += as.increment
	if as.currentSleep > as.maxSleep {
		as.currentSleep = as.maxSleep
	}
	as.lastBump = time.Now()

	slog.Debug("sleeper bump (rate limit)", "new_sleep", as.currentSleep)
}

// OnForbidden increases sleep duration in response to a 403.
func (as *AdaptiveSleeper) OnForbidden() {
	as.mu.Lock()
	defer as.mu.Unlock()

	// Smaller increment for 403 than 429.
	as.currentSleep += as.increment / 2
	if as.currentSleep > as.maxSleep {
		as.currentSleep = as.maxSleep
	}
	as.lastBump = time.Now()

	slog.Debug("sleeper bump (forbidden)", "new_sleep", as.currentSleep)
}

// Reset clears the sleep duration.
func (as *AdaptiveSleeper) Reset() {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.currentSleep = 0
}

// decay reduces the sleep duration if enough time has passed since the
// last rate limit bump. Called while holding the lock.
func (as *AdaptiveSleeper) decay() {
	if as.currentSleep <= 0 || as.lastBump.IsZero() {
		return
	}

	elapsed := time.Since(as.lastBump)
	if elapsed > as.decayAfter {
		// Halve the sleep duration for each decay period elapsed.
		periods := int(elapsed / as.decayAfter)
		for i := 0; i < periods; i++ {
			as.currentSleep /= 2
		}
		if as.currentSleep < as.minSleep {
			as.currentSleep = 0
		}
	}
}
