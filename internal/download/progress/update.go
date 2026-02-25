// =============================================================================
// FILE: internal/download/progress/update.go
// PURPOSE: Progress update tracking. Manages download progress state and
//          emits updates for display. Ports Python progress/update.py.
// =============================================================================

package progress

import (
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// Progress tracker
// ---------------------------------------------------------------------------

// Tracker tracks download progress for a batch of media items.
type Tracker struct {
	total      int64
	completed  int64
	failed     int64
	skipped    int64
	bytesTotal int64
	bytesDone  int64

	startTime time.Time
	mu        sync.RWMutex
	callbacks []func(Update)
}

// Update represents a progress update event.
type Update struct {
	Total      int64
	Completed  int64
	Failed     int64
	Skipped    int64
	BytesTotal int64
	BytesDone  int64
	Elapsed    time.Duration
	Speed      float64 // bytes per second
}

// NewTracker creates a progress tracker for the given total items.
//
// Parameters:
//   - total: Total number of items to track.
//
// Returns:
//   - A new Tracker.
func NewTracker(total int64) *Tracker {
	return &Tracker{
		total:     total,
		startTime: time.Now(),
	}
}

// OnUpdate registers a callback for progress updates.
func (t *Tracker) OnUpdate(fn func(Update)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.callbacks = append(t.callbacks, fn)
}

// AddCompleted increments the completed counter and emits an update.
func (t *Tracker) AddCompleted(bytes int64) {
	atomic.AddInt64(&t.completed, 1)
	atomic.AddInt64(&t.bytesDone, bytes)
	t.emit()
}

// AddFailed increments the failed counter and emits an update.
func (t *Tracker) AddFailed() {
	atomic.AddInt64(&t.failed, 1)
	t.emit()
}

// AddSkipped increments the skipped counter and emits an update.
func (t *Tracker) AddSkipped() {
	atomic.AddInt64(&t.skipped, 1)
	t.emit()
}

// SetBytesTotal sets the total expected bytes.
func (t *Tracker) SetBytesTotal(bytes int64) {
	atomic.StoreInt64(&t.bytesTotal, bytes)
}

// Current returns the current progress state.
func (t *Tracker) Current() Update {
	elapsed := time.Since(t.startTime)
	done := atomic.LoadInt64(&t.bytesDone)
	var speed float64
	if elapsed > 0 {
		speed = float64(done) / elapsed.Seconds()
	}

	return Update{
		Total:      atomic.LoadInt64(&t.total),
		Completed:  atomic.LoadInt64(&t.completed),
		Failed:     atomic.LoadInt64(&t.failed),
		Skipped:    atomic.LoadInt64(&t.skipped),
		BytesTotal: atomic.LoadInt64(&t.bytesTotal),
		BytesDone:  done,
		Elapsed:    elapsed,
		Speed:      speed,
	}
}

// Percent returns the completion percentage (0-100).
func (t *Tracker) Percent() float64 {
	total := atomic.LoadInt64(&t.total)
	if total == 0 {
		return 0
	}
	completed := atomic.LoadInt64(&t.completed) + atomic.LoadInt64(&t.failed) + atomic.LoadInt64(&t.skipped)
	return float64(completed) / float64(total) * 100
}

// emit sends an update to all registered callbacks.
func (t *Tracker) emit() {
	update := t.Current()
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, fn := range t.callbacks {
		fn(update)
	}
}
