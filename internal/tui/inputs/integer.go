// =============================================================================
// FILE: internal/tui/inputs/integer.go
// PURPOSE: Numeric-only text input widget with validation. Only accepts
//          digit characters and optional minus sign.
//          Ports Python classes/table/inputs/integer.py.
// =============================================================================

package inputs

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// IntegerInput
// ---------------------------------------------------------------------------

// IntegerInput is a text input that only accepts integer values.
type IntegerInput struct {
	input       textinput.Model
	focused     bool
	label       string
	allowNeg    bool
	err         error
}

// NewIntegerInput creates a new IntegerInput with the given label.
func NewIntegerInput(label string, allowNegative bool) IntegerInput {
	ti := textinput.New()
	if allowNegative {
		ti.Placeholder = "e.g. -10 or 42"
	} else {
		ti.Placeholder = "e.g. 42"
	}
	ti.CharLimit = 20
	ti.Width = 16
	ti.Validate = func(s string) error {
		if s == "" || s == "-" {
			return nil
		}
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	return IntegerInput{
		input:    ti,
		label:    label,
		allowNeg: allowNegative,
	}
}

// Focus gives focus to the input.
func (i *IntegerInput) Focus() tea.Cmd {
	i.focused = true
	return i.input.Focus()
}

// Blur removes focus from the input.
func (i *IntegerInput) Blur() {
	i.focused = false
	i.input.Blur()
}

// Focused returns whether the input is focused.
func (i *IntegerInput) Focused() bool {
	return i.focused
}

// Value returns the current text value.
func (i *IntegerInput) Value() string {
	return i.input.Value()
}

// IntValue returns the parsed integer value and any error.
func (i *IntegerInput) IntValue() (int64, error) {
	s := i.input.Value()
	if s == "" || s == "-" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// SetValue sets the input value.
func (i *IntegerInput) SetValue(s string) {
	i.input.SetValue(s)
}

// Reset clears the input.
func (i *IntegerInput) Reset() {
	i.input.SetValue("")
	i.err = nil
}

// Err returns the last validation error.
func (i *IntegerInput) Err() error {
	return i.err
}

// Update handles input events.
func (i *IntegerInput) Update(msg tea.Msg) tea.Cmd {
	if !i.focused {
		return nil
	}
	var cmd tea.Cmd
	i.input, cmd = i.input.Update(msg)

	// Validate the current value.
	s := i.input.Value()
	if s != "" && s != "-" {
		_, i.err = strconv.ParseInt(s, 10, 64)
	} else {
		i.err = nil
	}

	return cmd
}

// View renders the integer input.
func (i *IntegerInput) View() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginRight(1)

	result := labelStyle.Render(i.label+":") + " " + i.input.View()

	if i.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#DC2626"))
		result += " " + errStyle.Render("invalid")
	}

	return result
}
