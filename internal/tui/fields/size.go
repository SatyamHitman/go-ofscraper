// =============================================================================
// FILE: internal/tui/fields/size.go
// PURPOSE: File size range filter field. Allows filtering by min/max file size
//          in bytes, with human-readable display.
//          Ports Python classes/table/fields/size.py.
// =============================================================================

package fields

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// SizeField
// ---------------------------------------------------------------------------

// SizeUnit represents the unit for size display/input.
type SizeUnit int

const (
	SizeBytes SizeUnit = iota
	SizeKB
	SizeMB
	SizeGB
)

// sizeUnitNames maps SizeUnit to display suffix.
var sizeUnitNames = []string{"B", "KB", "MB", "GB"}

// sizeUnitMultipliers maps SizeUnit to bytes multiplier.
var sizeUnitMultipliers = []int64{1, 1024, 1024 * 1024, 1024 * 1024 * 1024}

// SizeField is a filter field for specifying a file size range.
type SizeField struct {
	name  string
	min   *int64 // Size in bytes.
	max   *int64 // Size in bytes.
	unit  SizeUnit
	focus SizeFocus
}

// SizeFocus indicates which sub-input is focused.
type SizeFocus int

const (
	SizeFocusMin SizeFocus = iota
	SizeFocusMax
)

// NewSizeField creates a new SizeField with the given display name.
func NewSizeField(name string) *SizeField {
	return &SizeField{
		name: name,
		unit: SizeMB,
	}
}

// Name returns the display name of the field.
func (f *SizeField) Name() string {
	return f.name
}

// Value returns the current range as "min=N;max=N" in bytes.
func (f *SizeField) Value() string {
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
func (f *SizeField) Reset() {
	f.min = nil
	f.max = nil
	f.focus = SizeFocusMin
}

// Min returns the minimum size in bytes, or nil if unset.
func (f *SizeField) Min() *int64 {
	return f.min
}

// Max returns the maximum size in bytes, or nil if unset.
func (f *SizeField) Max() *int64 {
	return f.max
}

// SetMin sets the minimum size from a string value in the current unit.
func (f *SizeField) SetMin(s string) error {
	if s == "" {
		f.min = nil
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid size: %w", err)
	}
	if v < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	bytes := int64(v * float64(sizeUnitMultipliers[f.unit]))
	f.min = &bytes
	return nil
}

// SetMax sets the maximum size from a string value in the current unit.
func (f *SizeField) SetMax(s string) error {
	if s == "" {
		f.max = nil
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid size: %w", err)
	}
	if v < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	bytes := int64(v * float64(sizeUnitMultipliers[f.unit]))
	f.max = &bytes
	return nil
}

// Unit returns the current display unit.
func (f *SizeField) Unit() SizeUnit {
	return f.unit
}

// SetUnit sets the display unit.
func (f *SizeField) SetUnit(u SizeUnit) {
	f.unit = u
}

// CycleUnit advances to the next size unit.
func (f *SizeField) CycleUnit() {
	f.unit = (f.unit + 1) % SizeUnit(len(sizeUnitNames))
}

// Focus returns the current sub-input focus.
func (f *SizeField) Focus() SizeFocus {
	return f.focus
}

// ToggleFocus switches between min and max focus.
func (f *SizeField) ToggleFocus() {
	if f.focus == SizeFocusMin {
		f.focus = SizeFocusMax
	} else {
		f.focus = SizeFocusMin
	}
}

// formatSize formats a byte count for display in the current unit.
func (f *SizeField) formatSize(bytes int64) string {
	val := float64(bytes) / float64(sizeUnitMultipliers[f.unit])
	return fmt.Sprintf("%.1f %s", val, sizeUnitNames[f.unit])
}

// View renders the field for display in the sidebar.
func (f *SizeField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	minStr := "any"
	if f.min != nil {
		minStr = f.formatSize(*f.min)
	}
	maxStr := "any"
	if f.max != nil {
		maxStr = f.formatSize(*f.max)
	}

	rangeStr := muted.Render(fmt.Sprintf("[%s .. %s]", minStr, maxStr))
	return fmt.Sprintf("%s %s", label, rangeStr)
}
