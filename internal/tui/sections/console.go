// =============================================================================
// FILE: internal/tui/sections/console.go
// PURPOSE: Console section. Fallback plain-text table renderer for
//          non-interactive mode (piped output, CI, or no TTY).
//          Ports Python classes/table/sections/console.py.
// =============================================================================

package sections

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// ConsoleSection
// ---------------------------------------------------------------------------

// ConsoleSection renders a plain-text table to stdout for non-interactive use.
type ConsoleSection struct {
	columns []Column
	rows    []Row
	writer  io.Writer
}

// NewConsoleSection creates a new ConsoleSection with the given columns.
func NewConsoleSection(columns []Column) *ConsoleSection {
	return &ConsoleSection{
		columns: columns,
		writer:  os.Stdout,
	}
}

// SetWriter sets the output writer (defaults to os.Stdout).
func (c *ConsoleSection) SetWriter(w io.Writer) {
	c.writer = w
}

// SetRows replaces the row data.
func (c *ConsoleSection) SetRows(rows []Row) {
	c.rows = rows
}

// Render outputs the table as plain text.
func (c *ConsoleSection) Render() string {
	if len(c.columns) == 0 {
		return ""
	}

	// Calculate column widths.
	widths := make([]int, len(c.columns))
	for i, col := range c.columns {
		widths[i] = len(col.Title)
		if col.Width > widths[i] {
			widths[i] = col.Width
		}
	}
	for _, row := range c.rows {
		for i, col := range c.columns {
			val := row[col.Key]
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
	}

	// Cap column widths.
	for i := range widths {
		if widths[i] > 50 {
			widths[i] = 50
		}
	}

	var sb strings.Builder

	// Header.
	var headerParts []string
	for i, col := range c.columns {
		headerParts = append(headerParts, padRight(col.Title, widths[i]))
	}
	sb.WriteString(strings.Join(headerParts, " | "))
	sb.WriteString("\n")

	// Separator.
	var sepParts []string
	for _, w := range widths {
		sepParts = append(sepParts, strings.Repeat("-", w))
	}
	sb.WriteString(strings.Join(sepParts, "-+-"))
	sb.WriteString("\n")

	// Rows.
	for _, row := range c.rows {
		var cells []string
		for i, col := range c.columns {
			cells = append(cells, padRight(row[col.Key], widths[i]))
		}
		sb.WriteString(strings.Join(cells, " | "))
		sb.WriteString("\n")
	}

	// Summary.
	sb.WriteString(fmt.Sprintf("\nTotal: %d rows\n", len(c.rows)))

	return sb.String()
}

// Print renders the table and writes it to the output writer.
func (c *ConsoleSection) Print() {
	content := c.Render()
	fmt.Fprint(c.writer, content)
}
