// =============================================================================
// FILE: internal/tui/fields/response.go
// PURPOSE: Response type filter field. Filters by post response type such as
//          timeline, archived, pinned, messages, etc.
//          Ports Python classes/table/fields/response.py.
// =============================================================================

package fields

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// ResponseField
// ---------------------------------------------------------------------------

// ResponseOption represents an available response type for filtering.
var ResponseOptions = []string{
	"timeline",
	"archived",
	"pinned",
	"streams",
	"stories",
	"highlights",
	"profile",
	"messages",
}

// ResponseField is a filter field for response type selection.
type ResponseField struct {
	name     string
	selected map[string]bool
}

// NewResponseField creates a new ResponseField with the given display name.
func NewResponseField(name string) *ResponseField {
	return &ResponseField{
		name:     name,
		selected: make(map[string]bool),
	}
}

// Name returns the display name of the field.
func (f *ResponseField) Name() string {
	return f.name
}

// Value returns the selected response types as a comma-separated string.
func (f *ResponseField) Value() string {
	var active []string
	for _, opt := range ResponseOptions {
		if f.selected[opt] {
			active = append(active, opt)
		}
	}
	return strings.Join(active, ",")
}

// Reset clears all selections.
func (f *ResponseField) Reset() {
	f.selected = make(map[string]bool)
}

// Toggle flips the selection state of the given response type.
func (f *ResponseField) Toggle(option string) {
	if f.selected[option] {
		delete(f.selected, option)
	} else {
		f.selected[option] = true
	}
}

// IsSelected returns whether the given option is currently selected.
func (f *ResponseField) IsSelected(option string) bool {
	return f.selected[option]
}

// Selected returns a slice of all selected response types.
func (f *ResponseField) Selected() []string {
	var result []string
	for _, opt := range ResponseOptions {
		if f.selected[opt] {
			result = append(result, opt)
		}
	}
	return result
}

// View renders the field for display in the sidebar.
func (f *ResponseField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	active := f.Value()
	if active == "" {
		active = "All"
	}

	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	if len(f.selected) > 0 {
		valStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	}

	indicator := valStyle.Render(fmt.Sprintf("[%s]", active))
	return fmt.Sprintf("%s %s", label, indicator)
}
