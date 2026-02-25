// =============================================================================
// FILE: internal/tui/fields/numeric.go
// PURPOSE: Numeric range filter field. Allows filtering by integer min/max
//          range. Ports Python classes/table/fields/numeric.py.
// =============================================================================

package fields

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// NumericField
// ---------------------------------------------------------------------------

// NumericField is a filter field for specifying an integer min/max range.
type NumericField struct {
	name  string
	min   *int64
	max   *int64
	focus NumericFocus
}

// NumericFocus indicates which sub-input is focused.
type NumericFocus int

const (
	NumericFocusMin NumericFocus = iota
	NumericFocusMax
)

// NewNumericField creates a new NumericField with the given display name.
func NewNumericField(name string) *NumericField {
	return &NumericField{
		name:  name,
		focus: NumericFocusMin,
	}
}

// Name returns the display name of the field.
func (f *NumericField) Name() string {
	return f.name
}

// Value returns the current range as "min=N;max=N". Empty components are omitted.
func (f *NumericField) Value() string {
	var parts []string
	if f.min != nil {
		parts = append(parts, fmt.Sprintf("min=%d", *f.min))
	}
	if f.max != nil {
		parts = append(parts, fmt.Sprintf("max=%d", *f.max))
	}
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ";" + parts[i]
	}
	return result
}

// Reset clears both bounds.
func (f *NumericField) Reset() {
	f.min = nil
	f.max = nil
	f.focus = NumericFocusMin
}

// Min returns the minimum value, or nil if unset.
func (f *NumericField) Min() *int64 {
	return f.min
}

// Max returns the maximum value, or nil if unset.
func (f *NumericField) Max() *int64 {
	return f.max
}

// SetMin sets the minimum value from a string. Empty string clears the value.
func (f *NumericField) SetMin(s string) error {
	if s == "" {
		f.min = nil
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %w", err)
	}
	f.min = &v
	return nil
}

// SetMax sets the maximum value from a string. Empty string clears the value.
func (f *NumericField) SetMax(s string) error {
	if s == "" {
		f.max = nil
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %w", err)
	}
	f.max = &v
	return nil
}

// Focus returns the current sub-input focus.
func (f *NumericField) Focus() NumericFocus {
	return f.focus
}

// ToggleFocus switches between min and max focus.
func (f *NumericField) ToggleFocus() {
	if f.focus == NumericFocusMin {
		f.focus = NumericFocusMax
	} else {
		f.focus = NumericFocusMin
	}
}

// View renders the field for display in the sidebar.
func (f *NumericField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	minStr := "any"
	if f.min != nil {
		minStr = fmt.Sprintf("%d", *f.min)
	}
	maxStr := "any"
	if f.max != nil {
		maxStr = fmt.Sprintf("%d", *f.max)
	}

	rangeStr := muted.Render(fmt.Sprintf("[%s .. %s]", minStr, maxStr))
	return fmt.Sprintf("%s %s", label, rangeStr)
}
