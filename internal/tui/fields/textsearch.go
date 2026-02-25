// =============================================================================
// FILE: internal/tui/fields/textsearch.go
// PURPOSE: Text search/regex filter field. Allows filtering by text content
//          using plain text or regular expressions.
//          Ports Python classes/table/fields/textsearch.py.
// =============================================================================

package fields

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// TextSearchField
// ---------------------------------------------------------------------------

// TextSearchMode indicates the search mode.
type TextSearchMode int

const (
	TextSearchPlain TextSearchMode = iota
	TextSearchRegex
)

// TextSearchField is a filter field for text/regex search.
type TextSearchField struct {
	name    string
	query   string
	mode    TextSearchMode
	pattern *regexp.Regexp // Compiled regex, nil if plain mode or invalid.
	err     error          // Last regex compilation error, if any.
}

// NewTextSearchField creates a new TextSearchField with the given display name.
func NewTextSearchField(name string) *TextSearchField {
	return &TextSearchField{
		name: name,
		mode: TextSearchPlain,
	}
}

// Name returns the display name of the field.
func (f *TextSearchField) Name() string {
	return f.name
}

// Value returns the current search query.
func (f *TextSearchField) Value() string {
	return f.query
}

// Reset clears the search query and resets to plain mode.
func (f *TextSearchField) Reset() {
	f.query = ""
	f.mode = TextSearchPlain
	f.pattern = nil
	f.err = nil
}

// Query returns the current search query string.
func (f *TextSearchField) Query() string {
	return f.query
}

// SetQuery sets the search query and recompiles the regex if in regex mode.
func (f *TextSearchField) SetQuery(q string) {
	f.query = q
	if f.mode == TextSearchRegex && q != "" {
		f.pattern, f.err = regexp.Compile(q)
	} else {
		f.pattern = nil
		f.err = nil
	}
}

// Mode returns the current search mode.
func (f *TextSearchField) Mode() TextSearchMode {
	return f.mode
}

// SetMode sets the search mode and recompiles if needed.
func (f *TextSearchField) SetMode(mode TextSearchMode) {
	f.mode = mode
	if mode == TextSearchRegex && f.query != "" {
		f.pattern, f.err = regexp.Compile(f.query)
	} else {
		f.pattern = nil
		f.err = nil
	}
}

// ToggleMode switches between plain and regex mode.
func (f *TextSearchField) ToggleMode() {
	if f.mode == TextSearchPlain {
		f.SetMode(TextSearchRegex)
	} else {
		f.SetMode(TextSearchPlain)
	}
}

// Pattern returns the compiled regex pattern, or nil.
func (f *TextSearchField) Pattern() *regexp.Regexp {
	return f.pattern
}

// Err returns the last regex compilation error, or nil.
func (f *TextSearchField) Err() error {
	return f.err
}

// View renders the field for display in the sidebar.
func (f *TextSearchField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	modeStr := "text"
	if f.mode == TextSearchRegex {
		modeStr = "regex"
	}

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))

	if f.query == "" {
		return fmt.Sprintf("%s %s", label, muted.Render(fmt.Sprintf("(%s)", modeStr)))
	}

	queryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	if f.err != nil {
		queryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DC2626"))
	}

	display := queryStyle.Render(fmt.Sprintf("\"%s\"", f.query))
	modeTag := muted.Render(fmt.Sprintf("(%s)", modeStr))
	return fmt.Sprintf("%s %s %s", label, display, modeTag)
}
