// =============================================================================
// FILE: internal/tui/live/screens.go
// PURPOSE: Screen manager. Manages switching between different display screens
//          in the live terminal output (progress, tasks, logs).
//          Ports Python utils/live/screens.py.
// =============================================================================

package live

import (
	"sync"
)

// ---------------------------------------------------------------------------
// Screen
// ---------------------------------------------------------------------------

// Screen represents a named display screen with render content.
type Screen struct {
	Name    string
	Content string
}

// ---------------------------------------------------------------------------
// ScreenManager
// ---------------------------------------------------------------------------

// ScreenManager manages multiple named screens and tracks the active one.
type ScreenManager struct {
	mu      sync.Mutex
	screens map[string]*Screen
	active  string
	order   []string // Ordered screen names for cycling.
}

// NewScreenManager creates a new ScreenManager.
func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		screens: make(map[string]*Screen),
	}
}

// AddScreen registers a new screen. The first screen added becomes active.
func (sm *ScreenManager) AddScreen(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.screens[name]; exists {
		return
	}

	sm.screens[name] = &Screen{Name: name}
	sm.order = append(sm.order, name)

	if sm.active == "" {
		sm.active = name
	}
}

// RemoveScreen removes a screen by name.
func (sm *ScreenManager) RemoveScreen(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.screens, name)

	for i, n := range sm.order {
		if n == name {
			sm.order = append(sm.order[:i], sm.order[i+1:]...)
			break
		}
	}

	if sm.active == name {
		if len(sm.order) > 0 {
			sm.active = sm.order[0]
		} else {
			sm.active = ""
		}
	}
}

// SetActive switches to the named screen.
func (sm *ScreenManager) SetActive(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, exists := sm.screens[name]; exists {
		sm.active = name
	}
}

// Active returns the name of the active screen.
func (sm *ScreenManager) Active() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.active
}

// CycleNext switches to the next screen in order.
func (sm *ScreenManager) CycleNext() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(sm.order) <= 1 {
		return
	}

	for i, name := range sm.order {
		if name == sm.active {
			sm.active = sm.order[(i+1)%len(sm.order)]
			return
		}
	}
}

// UpdateContent sets the content for the named screen.
func (sm *ScreenManager) UpdateContent(name, content string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, exists := sm.screens[name]; exists {
		s.Content = content
	}
}

// Render returns the content of the active screen.
func (sm *ScreenManager) Render() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.active == "" {
		return ""
	}
	if s, exists := sm.screens[sm.active]; exists {
		return s.Content
	}
	return ""
}

// ScreenNames returns the ordered list of screen names.
func (sm *ScreenManager) ScreenNames() []string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	result := make([]string, len(sm.order))
	copy(result, sm.order)
	return result
}

// ScreenCount returns the number of registered screens.
func (sm *ScreenManager) ScreenCount() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return len(sm.screens)
}
