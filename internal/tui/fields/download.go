// =============================================================================
// FILE: internal/tui/fields/download.go
// PURPOSE: Download status filter field. Filters media by download state:
//          downloaded, not downloaded, or all.
//          Ports Python classes/table/fields/download.py.
// =============================================================================

package fields

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// DownloadField
// ---------------------------------------------------------------------------

// DownloadFilter represents the download status filter option.
type DownloadFilter int

const (
	DownloadAll DownloadFilter = iota
	DownloadCompleted
	DownloadNotCompleted
)

// downloadFilterNames maps filter values to display names.
var downloadFilterNames = []string{"All", "Downloaded", "Not Downloaded"}

// DownloadField is a filter field for download status.
type DownloadField struct {
	name   string
	filter DownloadFilter
}

// NewDownloadField creates a new DownloadField with the given display name.
func NewDownloadField(name string) *DownloadField {
	return &DownloadField{
		name:   name,
		filter: DownloadAll,
	}
}

// Name returns the display name of the field.
func (f *DownloadField) Name() string {
	return f.name
}

// Value returns the current filter as a string.
func (f *DownloadField) Value() string {
	if f.filter == DownloadAll {
		return ""
	}
	return downloadFilterNames[f.filter]
}

// Reset clears the filter back to All.
func (f *DownloadField) Reset() {
	f.filter = DownloadAll
}

// Filter returns the current DownloadFilter.
func (f *DownloadField) Filter() DownloadFilter {
	return f.filter
}

// Cycle advances to the next filter option.
func (f *DownloadField) Cycle() {
	f.filter = (f.filter + 1) % DownloadFilter(len(downloadFilterNames))
}

// View renders the field for display in the sidebar.
func (f *DownloadField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	if f.filter != DownloadAll {
		valStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	}

	indicator := valStyle.Render(fmt.Sprintf("[%s]", downloadFilterNames[f.filter]))
	return fmt.Sprintf("%s %s", label, indicator)
}
