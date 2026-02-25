// =============================================================================
// FILE: internal/tui/utils/lock.go
// PURPOSE: TUI render lock. Prevents concurrent display updates by providing
//          a mutex-based lock for terminal rendering operations.
//          Ports Python utils/live/utils.py render lock.
// =============================================================================

package tuiutils

import (
	"sync"
)

// ---------------------------------------------------------------------------
// RenderLock
// ---------------------------------------------------------------------------

// RenderLock provides mutual exclusion for terminal rendering operations.
// It prevents concurrent display updates that can cause garbled output.
type RenderLock struct {
	mu sync.Mutex
}

// globalRenderLock is the package-level render lock instance.
var globalRenderLock RenderLock

// Lock acquires the render lock.
func (rl *RenderLock) Lock() {
	rl.mu.Lock()
}

// Unlock releases the render lock.
func (rl *RenderLock) Unlock() {
	rl.mu.Unlock()
}

// WithLock executes the given function while holding the render lock.
func (rl *RenderLock) WithLock(fn func()) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	fn()
}

// GlobalRenderLock returns the package-level render lock.
func GlobalRenderLock() *RenderLock {
	return &globalRenderLock
}

// AcquireRenderLock acquires the global render lock.
func AcquireRenderLock() {
	globalRenderLock.Lock()
}

// ReleaseRenderLock releases the global render lock.
func ReleaseRenderLock() {
	globalRenderLock.Unlock()
}

// WithRenderLock executes the given function while holding the global render lock.
func WithRenderLock(fn func()) {
	globalRenderLock.WithLock(fn)
}
