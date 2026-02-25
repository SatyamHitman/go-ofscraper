// =============================================================================
// FILE: internal/utils/sems.go
// PURPOSE: Semaphore utilities. Provides a weighted semaphore implementation
//          using channels for bounding concurrent operations (downloads, API
//          calls, DB writes, etc.). Ports Python utils/sems.py.
// =============================================================================

package utils

import (
	"context"
)

// ---------------------------------------------------------------------------
// Semaphore
// ---------------------------------------------------------------------------

// Semaphore is a counting semaphore implemented with a buffered channel.
// It bounds the number of concurrent operations to a fixed limit.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a semaphore with the given capacity.
//
// Parameters:
//   - n: Maximum concurrent holders. Must be > 0.
//
// Returns:
//   - A new Semaphore.
func NewSemaphore(n int) *Semaphore {
	if n <= 0 {
		n = 1
	}
	return &Semaphore{ch: make(chan struct{}, n)}
}

// Acquire blocks until a slot is available or the context is cancelled.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - Error if the context was cancelled before acquiring.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire attempts to acquire without blocking.
//
// Returns:
//   - true if a slot was acquired, false if none available.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release frees a previously acquired slot. Must be called exactly once per
// successful Acquire.
func (s *Semaphore) Release() {
	<-s.ch
}

// Available returns the number of currently available slots.
//
// Returns:
//   - Number of free slots.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}

// Cap returns the semaphore's total capacity.
//
// Returns:
//   - The maximum concurrent holders.
func (s *Semaphore) Cap() int {
	return cap(s.ch)
}
