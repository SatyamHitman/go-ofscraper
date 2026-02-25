// =============================================================================
// FILE: internal/utils/manager.go
// PURPOSE: Manager utility functions. Provides shared helpers used by the
//          various manager types (ModelManager, ProfileManager, etc.) for
//          common operations like progress tracking and batch processing.
//          Ports Python utils/manager.py.
// =============================================================================

package utils

import (
	"context"
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// Progress tracker
// ---------------------------------------------------------------------------

// ProgressTracker tracks progress of a batch operation. Thread-safe.
type ProgressTracker struct {
	total     int64
	completed atomic.Int64
	failed    atomic.Int64
	skipped   atomic.Int64
}

// NewProgressTracker creates a tracker for the given total item count.
//
// Parameters:
//   - total: Expected total items to process.
//
// Returns:
//   - A new ProgressTracker.
func NewProgressTracker(total int64) *ProgressTracker {
	return &ProgressTracker{total: total}
}

// IncrCompleted atomically increments the completed count.
func (pt *ProgressTracker) IncrCompleted() { pt.completed.Add(1) }

// IncrFailed atomically increments the failed count.
func (pt *ProgressTracker) IncrFailed() { pt.failed.Add(1) }

// IncrSkipped atomically increments the skipped count.
func (pt *ProgressTracker) IncrSkipped() { pt.skipped.Add(1) }

// Total returns the expected total.
func (pt *ProgressTracker) Total() int64 { return pt.total }

// Completed returns the completed count.
func (pt *ProgressTracker) Completed() int64 { return pt.completed.Load() }

// Failed returns the failed count.
func (pt *ProgressTracker) Failed() int64 { return pt.failed.Load() }

// Skipped returns the skipped count.
func (pt *ProgressTracker) Skipped() int64 { return pt.skipped.Load() }

// Processed returns total processed (completed + failed + skipped).
func (pt *ProgressTracker) Processed() int64 {
	return pt.completed.Load() + pt.failed.Load() + pt.skipped.Load()
}

// Percent returns completion percentage (0.0 to 100.0).
func (pt *ProgressTracker) Percent() float64 {
	if pt.total == 0 {
		return 100.0
	}
	return float64(pt.Processed()) / float64(pt.total) * 100.0
}

// ---------------------------------------------------------------------------
// Batch processor
// ---------------------------------------------------------------------------

// BatchProcessor runs a function on each item concurrently with a bounded
// number of workers. Collects errors for items that fail.
type BatchProcessor[T any] struct {
	Workers int // Max concurrent workers (0 = sequential)
}

// ProcessResult holds the outcome for a single batch item.
type ProcessResult[T any] struct {
	Item  T
	Err   error
	Index int
}

// Process runs fn on each item with the configured concurrency.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - items: Items to process.
//   - fn: Processing function for each item.
//
// Returns:
//   - A slice of results, one per item, in order.
func (bp *BatchProcessor[T]) Process(ctx context.Context, items []T, fn func(context.Context, T) error) []ProcessResult[T] {
	results := make([]ProcessResult[T], len(items))

	if bp.Workers <= 1 {
		// Sequential processing.
		for i, item := range items {
			results[i] = ProcessResult[T]{Item: item, Index: i}
			if ctx.Err() != nil {
				results[i].Err = ctx.Err()
				continue
			}
			results[i].Err = fn(ctx, item)
		}
		return results
	}

	// Concurrent processing with bounded workers.
	sem := make(chan struct{}, bp.Workers)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		sem <- struct{}{} // Acquire worker slot.

		go func(idx int, it T) {
			defer wg.Done()
			defer func() { <-sem }() // Release worker slot.

			results[idx] = ProcessResult[T]{Item: it, Index: idx}
			if ctx.Err() != nil {
				results[idx].Err = ctx.Err()
				return
			}
			results[idx].Err = fn(ctx, it)
		}(i, item)
	}

	wg.Wait()
	return results
}
