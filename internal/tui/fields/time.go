// =============================================================================
// FILE: internal/tui/fields/time.go
// PURPOSE: Time/duration range filter field. Allows filtering by min/max
//          duration (e.g., video length). Ports Python classes/table/fields/time.py.
// =============================================================================

package fields

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// TimeField
// ---------------------------------------------------------------------------

// TimeField is a filter field for specifying a duration range (min/max).
type TimeField struct {
	name  string
	min   *time.Duration
	max   *time.Duration
	focus TimeFocus
}

// TimeFocus indicates which sub-input is focused.
type TimeFocus int

const (
	TimeFocusMin TimeFocus = iota
	TimeFocusMax
)

// NewTimeField creates a new TimeField with the given display name.
func NewTimeField(name string) *TimeField {
	return &TimeField{
		name:  name,
		focus: TimeFocusMin,
	}
}

// Name returns the display name of the field.
func (f *TimeField) Name() string {
	return f.name
}

// Value returns the current range as "min=Ns;max=Ns" in seconds.
func (f *TimeField) Value() string {
	var parts []string
	if f.min != nil {
		parts = append(parts, fmt.Sprintf("min=%d", int64(f.min.Seconds())))
	}
	if f.max != nil {
		parts = append(parts, fmt.Sprintf("max=%d", int64(f.max.Seconds())))
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
func (f *TimeField) Reset() {
	f.min = nil
	f.max = nil
	f.focus = TimeFocusMin
}

// Min returns the minimum duration, or nil if unset.
func (f *TimeField) Min() *time.Duration {
	return f.min
}

// Max returns the maximum duration, or nil if unset.
func (f *TimeField) Max() *time.Duration {
	return f.max
}

// SetMin sets the minimum duration from a duration string (e.g., "1m30s", "90s").
// Empty string clears the value.
func (f *TimeField) SetMin(s string) error {
	if s == "" {
		f.min = nil
		return nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration (use e.g. 1m30s, 90s, 2h): %w", err)
	}
	if d < 0 {
		return fmt.Errorf("duration cannot be negative")
	}
	f.min = &d
	return nil
}

// SetMax sets the maximum duration from a duration string.
// Empty string clears the value.
func (f *TimeField) SetMax(s string) error {
	if s == "" {
		f.max = nil
		return nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration (use e.g. 1m30s, 90s, 2h): %w", err)
	}
	if d < 0 {
		return fmt.Errorf("duration cannot be negative")
	}
	f.max = &d
	return nil
}

// Focus returns the current sub-input focus.
func (f *TimeField) Focus() TimeFocus {
	return f.focus
}

// ToggleFocus switches between min and max focus.
func (f *TimeField) ToggleFocus() {
	if f.focus == TimeFocusMin {
		f.focus = TimeFocusMax
	} else {
		f.focus = TimeFocusMin
	}
}

// formatDuration formats a duration for display.
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// View renders the field for display in the sidebar.
func (f *TimeField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	minStr := "any"
	if f.min != nil {
		minStr = formatDuration(*f.min)
	}
	maxStr := "any"
	if f.max != nil {
		maxStr = formatDuration(*f.max)
	}

	rangeStr := muted.Render(fmt.Sprintf("[%s .. %s]", minStr, maxStr))
	return fmt.Sprintf("%s %s", label, rangeStr)
}
