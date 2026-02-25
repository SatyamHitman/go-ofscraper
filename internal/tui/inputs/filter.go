// =============================================================================
// FILE: internal/tui/inputs/filter.go
// PURPOSE: Filter text input widget. Used for filtering/searching within lists
//          (e.g., narrowing down items in a select field).
//          Ports Python classes/table/inputs/filter.py.
// =============================================================================

package inputs

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// FilterInput
// ---------------------------------------------------------------------------

// FilterInput is a text input widget for filtering list items.
type FilterInput struct {
	input    textinput.Model
	focused  bool
	label    string
}

// NewFilterInput creates a new FilterInput with the given label.
func NewFilterInput(label string) FilterInput {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.CharLimit = 64
	ti.Width = 24

	return FilterInput{
		input: ti,
		label: label,
	}
}

// Focus gives focus to the filter input.
func (f *FilterInput) Focus() tea.Cmd {
	f.focused = true
	return f.input.Focus()
}

// Blur removes focus from the filter input.
func (f *FilterInput) Blur() {
	f.focused = false
	f.input.Blur()
}

// Focused returns whether the input is focused.
func (f *FilterInput) Focused() bool {
	return f.focused
}

// Value returns the current filter text.
func (f *FilterInput) Value() string {
	return f.input.Value()
}

// SetValue sets the filter text.
func (f *FilterInput) SetValue(s string) {
	f.input.SetValue(s)
}

// Reset clears the filter text.
func (f *FilterInput) Reset() {
	f.input.SetValue("")
}

// Matches returns true if the given item matches the current filter text.
// Matching is case-insensitive substring search.
func (f *FilterInput) Matches(item string) bool {
	query := f.input.Value()
	if query == "" {
		return true
	}
	return strings.Contains(
		strings.ToLower(item),
		strings.ToLower(query),
	)
}

// Update handles input events.
func (f *FilterInput) Update(msg tea.Msg) tea.Cmd {
	if !f.focused {
		return nil
	}
	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return cmd
}

// View renders the filter input.
func (f *FilterInput) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginRight(1)

	return labelStyle.Render(f.label+":") + " " + f.input.View()
}
