// =============================================================================
// FILE: internal/tui/fields/price.go
// PURPOSE: Price range filter field. Allows filtering by min/max float price.
//          Ports Python classes/table/fields/price.py.
// =============================================================================

package fields

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// PriceField
// ---------------------------------------------------------------------------

// PriceField is a filter field for specifying a price range (min/max float).
type PriceField struct {
	name  string
	min   *float64
	max   *float64
	focus PriceFocus
}

// PriceFocus indicates which sub-input is focused.
type PriceFocus int

const (
	PriceFocusMin PriceFocus = iota
	PriceFocusMax
)

// NewPriceField creates a new PriceField with the given display name.
func NewPriceField(name string) *PriceField {
	return &PriceField{
		name:  name,
		focus: PriceFocusMin,
	}
}

// Name returns the display name of the field.
func (f *PriceField) Name() string {
	return f.name
}

// Value returns the current range as "min=N.NN;max=N.NN".
func (f *PriceField) Value() string {
	var parts []string
	if f.min != nil {
		parts = append(parts, fmt.Sprintf("min=%.2f", *f.min))
	}
	if f.max != nil {
		parts = append(parts, fmt.Sprintf("max=%.2f", *f.max))
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
func (f *PriceField) Reset() {
	f.min = nil
	f.max = nil
	f.focus = PriceFocusMin
}

// Min returns the minimum price, or nil if unset.
func (f *PriceField) Min() *float64 {
	return f.min
}

// Max returns the maximum price, or nil if unset.
func (f *PriceField) Max() *float64 {
	return f.max
}

// SetMin sets the minimum price from a string. Empty string clears the value.
func (f *PriceField) SetMin(s string) error {
	if s == "" {
		f.min = nil
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid price: %w", err)
	}
	if v < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	f.min = &v
	return nil
}

// SetMax sets the maximum price from a string. Empty string clears the value.
func (f *PriceField) SetMax(s string) error {
	if s == "" {
		f.max = nil
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid price: %w", err)
	}
	if v < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	f.max = &v
	return nil
}

// Focus returns the current sub-input focus.
func (f *PriceField) Focus() PriceFocus {
	return f.focus
}

// ToggleFocus switches between min and max focus.
func (f *PriceField) ToggleFocus() {
	if f.focus == PriceFocusMin {
		f.focus = PriceFocusMax
	} else {
		f.focus = PriceFocusMin
	}
}

// View renders the field for display in the sidebar.
func (f *PriceField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	minStr := "any"
	if f.min != nil {
		minStr = fmt.Sprintf("$%.2f", *f.min)
	}
	maxStr := "any"
	if f.max != nil {
		maxStr = fmt.Sprintf("$%.2f", *f.max)
	}

	rangeStr := muted.Render(fmt.Sprintf("[%s .. %s]", minStr, maxStr))
	return fmt.Sprintf("%s %s", label, rangeStr)
}
