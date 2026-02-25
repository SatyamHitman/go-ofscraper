// =============================================================================
// FILE: internal/tui/fields/date.go
// PURPOSE: Date range filter field. Allows filtering by after/before dates.
//          Ports Python classes/table/fields/date.py.
// =============================================================================

package fields

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// DateField
// ---------------------------------------------------------------------------

// DateField is a filter field for specifying a date range (after/before).
type DateField struct {
	name   string
	after  string // Date string in YYYY-MM-DD format, or empty.
	before string // Date string in YYYY-MM-DD format, or empty.
	focus  DateFocus
}

// DateFocus indicates which sub-input is focused.
type DateFocus int

const (
	DateFocusAfter DateFocus = iota
	DateFocusBefore
)

// DateLayout is the expected date format.
const DateLayout = "2006-01-02"

// NewDateField creates a new DateField with the given display name.
func NewDateField(name string) *DateField {
	return &DateField{
		name:  name,
		focus: DateFocusAfter,
	}
}

// Name returns the display name of the field.
func (f *DateField) Name() string {
	return f.name
}

// Value returns the current value as "after=YYYY-MM-DD;before=YYYY-MM-DD".
// Empty components are omitted.
func (f *DateField) Value() string {
	var parts []string
	if f.after != "" {
		parts = append(parts, "after="+f.after)
	}
	if f.before != "" {
		parts = append(parts, "before="+f.before)
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

// Reset clears both date bounds.
func (f *DateField) Reset() {
	f.after = ""
	f.before = ""
	f.focus = DateFocusAfter
}

// After returns the after date string.
func (f *DateField) After() string {
	return f.after
}

// Before returns the before date string.
func (f *DateField) Before() string {
	return f.before
}

// SetAfter sets the after date, validating the format.
func (f *DateField) SetAfter(s string) error {
	if s == "" {
		f.after = ""
		return nil
	}
	_, err := time.Parse(DateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	f.after = s
	return nil
}

// SetBefore sets the before date, validating the format.
func (f *DateField) SetBefore(s string) error {
	if s == "" {
		f.before = ""
		return nil
	}
	_, err := time.Parse(DateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	f.before = s
	return nil
}

// Focus returns the current sub-input focus.
func (f *DateField) Focus() DateFocus {
	return f.focus
}

// SetFocus sets which sub-input is focused.
func (f *DateField) SetFocus(focus DateFocus) {
	f.focus = focus
}

// ToggleFocus switches between after and before focus.
func (f *DateField) ToggleFocus() {
	if f.focus == DateFocusAfter {
		f.focus = DateFocusBefore
	} else {
		f.focus = DateFocusAfter
	}
}

// View renders the field for display in the sidebar.
func (f *DateField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	afterVal := f.after
	if afterVal == "" {
		afterVal = "any"
	}
	beforeVal := f.before
	if beforeVal == "" {
		beforeVal = "any"
	}

	afterStr := muted.Render(fmt.Sprintf("After: %s", afterVal))
	beforeStr := muted.Render(fmt.Sprintf("Before: %s", beforeVal))

	return fmt.Sprintf("%s\n  %s\n  %s", label, afterStr, beforeStr)
}
