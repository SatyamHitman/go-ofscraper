// =============================================================================
// FILE: internal/tui/fields/bool.go
// PURPOSE: Boolean toggle filter field. Cycles through true/false/unset states
//          for binary filter conditions. Ports Python classes/table/fields/bool.py.
// =============================================================================

package fields

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// BoolField
// ---------------------------------------------------------------------------

// BoolState represents the three-state boolean: unset, true, or false.
type BoolState int

const (
	BoolUnset BoolState = iota
	BoolTrue
	BoolFalse
)

// BoolField is a filter field that toggles between true, false, and unset.
type BoolField struct {
	name  string
	state BoolState
}

// NewBoolField creates a new BoolField with the given display name.
func NewBoolField(name string) *BoolField {
	return &BoolField{
		name:  name,
		state: BoolUnset,
	}
}

// Name returns the display name of the field.
func (f *BoolField) Name() string {
	return f.name
}

// Value returns the current value as a string: "", "true", or "false".
func (f *BoolField) Value() string {
	switch f.state {
	case BoolTrue:
		return "true"
	case BoolFalse:
		return "false"
	default:
		return ""
	}
}

// Reset clears the field back to unset.
func (f *BoolField) Reset() {
	f.state = BoolUnset
}

// Toggle advances the state: unset -> true -> false -> unset.
func (f *BoolField) Toggle() {
	switch f.state {
	case BoolUnset:
		f.state = BoolTrue
	case BoolTrue:
		f.state = BoolFalse
	case BoolFalse:
		f.state = BoolUnset
	}
}

// State returns the current BoolState.
func (f *BoolField) State() BoolState {
	return f.state
}

// View renders the field for display in the sidebar.
func (f *BoolField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	var indicator string
	switch f.state {
	case BoolTrue:
		indicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#059669")).Render("[Yes]")
	case BoolFalse:
		indicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC2626")).Render("[No]")
	default:
		indicator = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("[Any]")
	}

	return fmt.Sprintf("%s %s", label, indicator)
}
