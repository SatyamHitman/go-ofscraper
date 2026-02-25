// =============================================================================
// FILE: internal/tui/live/updater.go
// PURPOSE: Periodic refresh timer for live display. Triggers display updates
//          at a configurable interval. Ports Python utils/live/updater.py.
// =============================================================================

package live

import (
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Updater
// ---------------------------------------------------------------------------

// UpdateFunc is called on each refresh tick.
type UpdateFunc func()

// Updater periodically triggers display refreshes.
type Updater struct {
	mu       sync.Mutex
	interval time.Duration
	ticker   *time.Ticker
	stopCh   chan struct{}
	running  bool
	onUpdate UpdateFunc
}

// NewUpdater creates a new Updater with the given refresh interval.
func NewUpdater(interval time.Duration) *Updater {
	return &Updater{
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// SetInterval changes the refresh interval. Takes effect on next Start.
func (u *Updater) SetInterval(d time.Duration) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.interval = d
}

// SetOnUpdate sets the callback function invoked on each tick.
func (u *Updater) SetOnUpdate(fn UpdateFunc) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.onUpdate = fn
}

// Start begins periodic updates. Does nothing if already running.
func (u *Updater) Start() {
	u.mu.Lock()
	if u.running {
		u.mu.Unlock()
		return
	}
	u.running = true
	u.ticker = time.NewTicker(u.interval)
	u.stopCh = make(chan struct{})
	u.mu.Unlock()

	go u.run()
}

// Stop stops the periodic updates.
func (u *Updater) Stop() {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.running {
		return
	}
	u.running = false
	close(u.stopCh)
	u.ticker.Stop()
}

// Running returns whether the updater is currently active.
func (u *Updater) Running() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.running
}

// run is the internal goroutine loop.
func (u *Updater) run() {
	for {
		select {
		case <-u.stopCh:
			return
		case <-u.ticker.C:
			u.mu.Lock()
			fn := u.onUpdate
			u.mu.Unlock()
			if fn != nil {
				fn()
			}
		}
	}
}
