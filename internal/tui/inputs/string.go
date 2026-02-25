// =============================================================================
// FILE: internal/tui/inputs/string.go
// PURPOSE: General-purpose text input widget. A thin wrapper around the
//          bubbletea textinput component with label support.
//          Ports Python classes/table/inputs/string.py.
// =============================================================================

package inputs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// StringInput
// ---------------------------------------------------------------------------

// StringInput is a general-purpose text input widget.
type StringInput struct {
	input       textinput.Model
	focused     bool
	label       string
}

// NewStringInput creates a new StringInput with the given label and placeholder.
func NewStringInput(label, placeholder string) StringInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 24

	return StringInput{
		input: ti,
		label: label,
	}
}

// Focus gives focus to the input.
func (s *StringInput) Focus() tea.Cmd {
	s.focused = true
	return s.input.Focus()
}

// Blur removes focus from the input.
func (s *StringInput) Blur() {
	s.focused = false
	s.input.Blur()
}

// Focused returns whether the input is focused.
func (s *StringInput) Focused() bool {
	return s.focused
}

// Value returns the current text value.
func (s *StringInput) Value() string {
	return s.input.Value()
}

// SetValue sets the input value.
func (s *StringInput) SetValue(v string) {
	s.input.SetValue(v)
}

// Reset clears the input.
func (s *StringInput) Reset() {
	s.input.SetValue("")
}

// SetWidth sets the input width.
func (s *StringInput) SetWidth(w int) {
	s.input.Width = w
}

// SetCharLimit sets the maximum character count.
func (s *StringInput) SetCharLimit(n int) {
	s.input.CharLimit = n
}

// Update handles input events.
func (s *StringInput) Update(msg tea.Msg) tea.Cmd {
	if !s.focused {
		return nil
	}
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return cmd
}

// View renders the string input.
func (s *StringInput) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginRight(1)

	return labelStyle.Render(s.label+":") + " " + s.input.View()
}
