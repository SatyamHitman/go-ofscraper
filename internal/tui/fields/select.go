// =============================================================================
// FILE: internal/tui/fields/select.go
// PURPOSE: Generic multi-select filter field. Allows selecting multiple items
//          from a predefined list of choices.
//          Ports Python classes/table/fields/select.py.
// =============================================================================

package fields

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// SelectField
// ---------------------------------------------------------------------------

// SelectField is a generic filter field that allows selecting multiple items
// from a list of choices.
type SelectField struct {
	name     string
	choices  []string
	selected map[string]bool
	cursor   int
}

// NewSelectField creates a new SelectField with the given name and choices.
func NewSelectField(name string, choices []string) *SelectField {
	return &SelectField{
		name:     name,
		choices:  choices,
		selected: make(map[string]bool),
	}
}

// Name returns the display name of the field.
func (f *SelectField) Name() string {
	return f.name
}

// Value returns the selected choices as a comma-separated string.
func (f *SelectField) Value() string {
	var active []string
	for _, c := range f.choices {
		if f.selected[c] {
			active = append(active, c)
		}
	}
	return strings.Join(active, ",")
}

// Reset clears all selections and resets the cursor.
func (f *SelectField) Reset() {
	f.selected = make(map[string]bool)
	f.cursor = 0
}

// Choices returns the available choices.
func (f *SelectField) Choices() []string {
	return f.choices
}

// Toggle flips the selection state of the choice at the current cursor.
func (f *SelectField) Toggle() {
	if f.cursor < 0 || f.cursor >= len(f.choices) {
		return
	}
	choice := f.choices[f.cursor]
	if f.selected[choice] {
		delete(f.selected, choice)
	} else {
		f.selected[choice] = true
	}
}

// ToggleChoice flips the selection state of a specific choice by name.
func (f *SelectField) ToggleChoice(choice string) {
	if f.selected[choice] {
		delete(f.selected, choice)
	} else {
		f.selected[choice] = true
	}
}

// IsSelected returns whether the given choice is selected.
func (f *SelectField) IsSelected(choice string) bool {
	return f.selected[choice]
}

// Selected returns a slice of all selected choices.
func (f *SelectField) Selected() []string {
	var result []string
	for _, c := range f.choices {
		if f.selected[c] {
			result = append(result, c)
		}
	}
	return result
}

// Cursor returns the current cursor position.
func (f *SelectField) Cursor() int {
	return f.cursor
}

// CursorUp moves the cursor up one position.
func (f *SelectField) CursorUp() {
	if f.cursor > 0 {
		f.cursor--
	}
}

// CursorDown moves the cursor down one position.
func (f *SelectField) CursorDown() {
	if f.cursor < len(f.choices)-1 {
		f.cursor++
	}
}

// View renders the field for display in the sidebar.
func (f *SelectField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	active := f.Value()
	if active == "" {
		active = "None"
	}

	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	if len(f.selected) > 0 {
		valStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	}

	indicator := valStyle.Render(fmt.Sprintf("[%s]", active))
	return fmt.Sprintf("%s %s", label, indicator)
}
