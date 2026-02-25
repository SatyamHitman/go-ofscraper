// =============================================================================
// FILE: internal/tui/sections/sidebar.go
// PURPOSE: Sidebar section. Displays the filter fields list in the TUI,
//          handles selection and editing of individual filter fields.
//          Ports Python classes/table/sections/sidebar.py.
// =============================================================================

package sections

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Field interface
// ---------------------------------------------------------------------------

// Field is the interface that all sidebar filter fields must implement.
type Field interface {
	Name() string
	Value() string
	Reset()
	View() string
}

// ---------------------------------------------------------------------------
// Sidebar
// ---------------------------------------------------------------------------

// Sidebar displays a list of filter fields and handles selection/editing.
type Sidebar struct {
	fields  []Field
	cursor  int
	editing bool
	width   int
	height  int
}

// NewSidebar creates a new Sidebar with the given fields.
func NewSidebar(fields []Field) *Sidebar {
	return &Sidebar{
		fields: fields,
	}
}

// SetSize sets the available dimensions for the sidebar.
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// Fields returns the list of fields.
func (s *Sidebar) Fields() []Field {
	return s.fields
}

// Cursor returns the current cursor position.
func (s *Sidebar) Cursor() int {
	return s.cursor
}

// Editing returns whether the sidebar is in editing mode.
func (s *Sidebar) Editing() bool {
	return s.editing
}

// SetEditing sets the editing mode.
func (s *Sidebar) SetEditing(editing bool) {
	s.editing = editing
}

// CursorUp moves the cursor up one position.
func (s *Sidebar) CursorUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

// CursorDown moves the cursor down one position.
func (s *Sidebar) CursorDown() {
	if s.cursor < len(s.fields)-1 {
		s.cursor++
	}
}

// SelectedField returns the currently selected field, or nil.
func (s *Sidebar) SelectedField() Field {
	if s.cursor < 0 || s.cursor >= len(s.fields) {
		return nil
	}
	return s.fields[s.cursor]
}

// ResetAll resets all fields to their default values.
func (s *Sidebar) ResetAll() {
	for _, f := range s.fields {
		f.Reset()
	}
}

// ActiveFilterCount returns the number of fields with a non-empty value.
func (s *Sidebar) ActiveFilterCount() int {
	count := 0
	for _, f := range s.fields {
		if f.Value() != "" {
			count++
		}
	}
	return count
}

// Update handles input events for the sidebar.
func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.editing {
			return nil // Editing is handled by the parent.
		}
		switch msg.String() {
		case "up", "k":
			s.CursorUp()
		case "down", "j":
			s.CursorDown()
		case "enter":
			s.editing = true
		case "r":
			if f := s.SelectedField(); f != nil {
				f.Reset()
			}
		case "R":
			s.ResetAll()
		}
	}
	return nil
}

// View renders the sidebar.
func (s *Sidebar) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginBottom(1)

	activeCount := s.ActiveFilterCount()
	title := "Filters"
	if activeCount > 0 {
		title = fmt.Sprintf("Filters (%d active)", activeCount)
	}

	var content string
	content += titleStyle.Render(title) + "\n\n"

	for i, f := range s.fields {
		cursor := "  "
		if i == s.cursor {
			cursor = "> "
			if s.editing {
				cursor = "* "
			}
		}

		line := f.View()
		if i == s.cursor {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C3AED")).
				Bold(true).
				Render(line)
		}

		content += cursor + line + "\n"
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginTop(1)

	content += "\n" + helpStyle.Render("j/k: navigate  enter: edit  r: reset  R: reset all")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#374151")).
		Padding(1, 1).
		Width(s.width).
		Height(s.height)

	return boxStyle.Render(content)
}
