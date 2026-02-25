// =============================================================================
// FILE: internal/tui/live/live.go
// PURPOSE: Live display manager. Manages real-time terminal output for
//          download progress, task status, and log display.
//          Ports Python utils/live/live.py.
// =============================================================================

package live

import (
	"fmt"
	"sync"
)

// ---------------------------------------------------------------------------
// Display manager
// ---------------------------------------------------------------------------

// Display manages live terminal updates.
type Display struct {
	mu      sync.Mutex
	enabled bool
	lines   []string
}

// New creates a new live display.
func New(enabled bool) *Display {
	return &Display{
		enabled: enabled,
	}
}

// SetEnabled toggles live display mode.
func (d *Display) SetEnabled(enabled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = enabled
}

// Update replaces all displayed lines with the given content.
func (d *Display) Update(lines []string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.enabled {
		return
	}
	d.lines = lines
}

// Render returns the current display content.
func (d *Display) Render() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	var out string
	for _, line := range d.lines {
		out += line + "\n"
	}
	return out
}

// Print outputs the current display to stdout.
func (d *Display) Print() {
	content := d.Render()
	if content != "" {
		fmt.Print(content)
	}
}

// Clear removes all displayed content.
func (d *Display) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lines = nil
}
