// =============================================================================
// FILE: internal/tui/live/groups.go
// PURPOSE: Display group. Groups related display elements together for
//          organized rendering in the live terminal output.
//          Ports Python utils/live/groups.py.
// =============================================================================

package live

import (
	"strings"
	"sync"
)

// ---------------------------------------------------------------------------
// DisplayGroup
// ---------------------------------------------------------------------------

// DisplayElement is something that can render itself as a string.
type DisplayElement interface {
	Render() string
}

// DisplayGroup groups related display elements for rendering together.
type DisplayGroup struct {
	mu       sync.Mutex
	name     string
	elements []DisplayElement
	visible  bool
}

// NewDisplayGroup creates a new DisplayGroup with the given name.
func NewDisplayGroup(name string) *DisplayGroup {
	return &DisplayGroup{
		name:    name,
		visible: true,
	}
}

// Name returns the group name.
func (g *DisplayGroup) Name() string {
	return g.name
}

// AddElement adds a display element to the group.
func (g *DisplayGroup) AddElement(elem DisplayElement) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.elements = append(g.elements, elem)
}

// ClearElements removes all elements from the group.
func (g *DisplayGroup) ClearElements() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.elements = nil
}

// ElementCount returns the number of elements in the group.
func (g *DisplayGroup) ElementCount() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.elements)
}

// SetVisible sets whether the group is visible.
func (g *DisplayGroup) SetVisible(visible bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.visible = visible
}

// Visible returns whether the group is visible.
func (g *DisplayGroup) Visible() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.visible
}

// Render returns the combined output of all elements in the group.
// Returns empty string if the group is not visible.
func (g *DisplayGroup) Render() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.visible || len(g.elements) == 0 {
		return ""
	}

	var parts []string
	for _, elem := range g.elements {
		rendered := elem.Render()
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	return strings.Join(parts, "\n")
}
