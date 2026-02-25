// =============================================================================
// FILE: internal/prompts/keybindings.go
// PURPOSE: Keybinding definitions for interactive prompts and TUI.
//          Maps keyboard shortcuts to prompt actions.
//          Ports Python prompts/keybindings.py.
// =============================================================================

package prompts

// ---------------------------------------------------------------------------
// Keybinding definitions
// ---------------------------------------------------------------------------

// KeyBinding represents a keyboard shortcut mapping.
type KeyBinding struct {
	Key         string // Key name (e.g. "enter", "ctrl+c", "tab")
	Description string // Human-readable description
	Action      string // Action identifier
}

// DefaultKeybindings returns the default keybinding set for interactive prompts.
func DefaultKeybindings() []KeyBinding {
	return []KeyBinding{
		{Key: "enter", Description: "Confirm selection", Action: "confirm"},
		{Key: "space", Description: "Toggle item", Action: "toggle"},
		{Key: "up", Description: "Move up", Action: "up"},
		{Key: "down", Description: "Move down", Action: "down"},
		{Key: "tab", Description: "Next field", Action: "next"},
		{Key: "shift+tab", Description: "Previous field", Action: "prev"},
		{Key: "ctrl+a", Description: "Select all", Action: "select_all"},
		{Key: "ctrl+n", Description: "Select none", Action: "select_none"},
		{Key: "ctrl+c", Description: "Cancel/Exit", Action: "cancel"},
		{Key: "esc", Description: "Back/Cancel", Action: "back"},
		{Key: "/", Description: "Search/Filter", Action: "search"},
	}
}

// KeybindingHelp returns a formatted help string for all keybindings.
func KeybindingHelp() string {
	bindings := DefaultKeybindings()
	maxKeyLen := 0
	for _, b := range bindings {
		if len(b.Key) > maxKeyLen {
			maxKeyLen = len(b.Key)
		}
	}

	var help string
	for _, b := range bindings {
		padding := maxKeyLen - len(b.Key)
		help += b.Key
		for i := 0; i < padding+2; i++ {
			help += " "
		}
		help += b.Description + "\n"
	}
	return help
}
