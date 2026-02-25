// =============================================================================
// FILE: internal/tui/sections/table.go
// PURPOSE: Table section. Displays data in a sortable, paginated table with
//          columns and rows. Ports Python classes/table/sections/table.py.
// =============================================================================

package sections

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Table types
// ---------------------------------------------------------------------------

// Column defines a table column.
type Column struct {
	Key   string
	Title string
	Width int
}

// Row is a table row mapping column keys to display values.
type Row map[string]string

// SortOrder represents the sort direction.
type SortOrder int

const (
	SortNone SortOrder = iota
	SortAsc
	SortDesc
)

// ---------------------------------------------------------------------------
// TableSection
// ---------------------------------------------------------------------------

// TableSection displays data in a paginated, sortable table.
type TableSection struct {
	columns   []Column
	rows      []Row
	cursor    int
	page      int
	pageSize  int
	sortCol   string
	sortOrder SortOrder
	width     int
	height    int
}

// NewTableSection creates a new TableSection with the given columns.
func NewTableSection(columns []Column) *TableSection {
	return &TableSection{
		columns:  columns,
		pageSize: 20,
	}
}

// SetSize sets the available dimensions for the table.
func (t *TableSection) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// SetRows replaces the table data. Resets cursor and page.
func (t *TableSection) SetRows(rows []Row) {
	t.rows = rows
	t.cursor = 0
	t.page = 0
}

// SetPageSize sets the number of rows per page.
func (t *TableSection) SetPageSize(size int) {
	if size > 0 {
		t.pageSize = size
	}
}

// Rows returns the current row data.
func (t *TableSection) Rows() []Row {
	return t.rows
}

// SelectedRow returns the currently selected row, or nil.
func (t *TableSection) SelectedRow() Row {
	visible := t.visibleRows()
	if t.cursor < 0 || t.cursor >= len(visible) {
		return nil
	}
	return visible[t.cursor]
}

// TotalPages returns the total number of pages.
func (t *TableSection) TotalPages() int {
	if len(t.rows) == 0 {
		return 1
	}
	pages := len(t.rows) / t.pageSize
	if len(t.rows)%t.pageSize != 0 {
		pages++
	}
	return pages
}

// Page returns the current page (0-indexed).
func (t *TableSection) Page() int {
	return t.page
}

// visibleRows returns the rows for the current page.
func (t *TableSection) visibleRows() []Row {
	start := t.page * t.pageSize
	if start >= len(t.rows) {
		return nil
	}
	end := start + t.pageSize
	if end > len(t.rows) {
		end = len(t.rows)
	}
	return t.rows[start:end]
}

// Sort sorts the rows by the given column key.
func (t *TableSection) Sort(colKey string) {
	if t.sortCol == colKey {
		// Toggle sort order.
		switch t.sortOrder {
		case SortAsc:
			t.sortOrder = SortDesc
		case SortDesc:
			t.sortOrder = SortNone
			t.sortCol = ""
			return
		default:
			t.sortOrder = SortAsc
		}
	} else {
		t.sortCol = colKey
		t.sortOrder = SortAsc
	}

	order := t.sortOrder
	sort.SliceStable(t.rows, func(i, j int) bool {
		a := t.rows[i][colKey]
		b := t.rows[j][colKey]
		if order == SortDesc {
			return a > b
		}
		return a < b
	})

	t.cursor = 0
	t.page = 0
}

// Update handles input events for the table.
func (t *TableSection) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}
		case "down", "j":
			visible := t.visibleRows()
			if t.cursor < len(visible)-1 {
				t.cursor++
			}
		case "left", "h":
			if t.page > 0 {
				t.page--
				t.cursor = 0
			}
		case "right", "l":
			if t.page < t.TotalPages()-1 {
				t.page++
				t.cursor = 0
			}
		case "home":
			t.page = 0
			t.cursor = 0
		case "end":
			t.page = t.TotalPages() - 1
			t.cursor = 0
		}
	}
	return nil
}

// View renders the table.
func (t *TableSection) View() string {
	if len(t.columns) == 0 {
		return "No columns defined"
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))
	rowStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB"))
	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1F2937"))

	// Render header.
	var headerCells []string
	for _, col := range t.columns {
		title := col.Title
		if col.Key == t.sortCol {
			switch t.sortOrder {
			case SortAsc:
				title += " ^"
			case SortDesc:
				title += " v"
			}
		}
		headerCells = append(headerCells, padRight(title, col.Width))
	}
	header := headerStyle.Render(strings.Join(headerCells, " | "))

	// Separator.
	sepWidth := 0
	for _, col := range t.columns {
		sepWidth += col.Width
	}
	sepWidth += (len(t.columns) - 1) * 3 // " | " separators
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Render(strings.Repeat("─", sepWidth))

	// Render rows.
	visible := t.visibleRows()
	var rowLines []string
	for i, row := range visible {
		var cells []string
		for _, col := range t.columns {
			cells = append(cells, padRight(row[col.Key], col.Width))
		}
		line := strings.Join(cells, " | ")
		if i == t.cursor {
			line = selectedStyle.Render(line)
		} else {
			line = rowStyle.Render(line)
		}
		rowLines = append(rowLines, line)
	}

	// Footer with pagination.
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginTop(1)
	footer := footerStyle.Render(
		fmt.Sprintf("Page %d/%d  |  %d rows  |  h/l: page  j/k: navigate",
			t.page+1, t.TotalPages(), len(t.rows)))

	content := header + "\n" + separator + "\n"
	content += strings.Join(rowLines, "\n")
	content += "\n" + footer

	return content
}

// padRight pads a string to the given width, truncating if necessary.
func padRight(s string, width int) string {
	if len(s) > width {
		if width > 3 {
			return s[:width-3] + "..."
		}
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}
