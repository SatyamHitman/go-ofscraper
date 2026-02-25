// =============================================================================
// FILE: internal/download/globals.go
// PURPOSE: Download global state. Tracks global download counters and provides
//          session-level statistics. Ports Python actions/utils/globals.py.
// =============================================================================

package download

import (
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// Global counters
// ---------------------------------------------------------------------------

var (
	globalTotal     int64
	globalSucceeded int64
	globalFailed    int64
	globalSkipped   int64
)

// IncrTotal increments the global total counter.
func IncrTotal(n int64) { atomic.AddInt64(&globalTotal, n) }

// IncrSucceeded increments the global success counter.
func IncrSucceeded() { atomic.AddInt64(&globalSucceeded, 1) }

// IncrFailed increments the global failure counter.
func IncrFailed() { atomic.AddInt64(&globalFailed, 1) }

// IncrSkipped increments the global skip counter.
func IncrSkipped() { atomic.AddInt64(&globalSkipped, 1) }

// GlobalStats returns the current global download statistics.
func GlobalStats() (total, succeeded, failed, skipped int64) {
	return atomic.LoadInt64(&globalTotal),
		atomic.LoadInt64(&globalSucceeded),
		atomic.LoadInt64(&globalFailed),
		atomic.LoadInt64(&globalSkipped)
}

// ResetGlobalStats resets all global counters to zero.
func ResetGlobalStats() {
	atomic.StoreInt64(&globalTotal, 0)
	atomic.StoreInt64(&globalSucceeded, 0)
	atomic.StoreInt64(&globalFailed, 0)
	atomic.StoreInt64(&globalSkipped, 0)
}
