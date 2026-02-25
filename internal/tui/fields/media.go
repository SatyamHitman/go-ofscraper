// =============================================================================
// FILE: internal/tui/fields/media.go
// PURPOSE: Media type filter field. Filters by media type: images, videos,
//          audio, or all. Ports Python classes/table/fields/media.py.
// =============================================================================

package fields

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// MediaField
// ---------------------------------------------------------------------------

// MediaFilter represents the media type filter option.
type MediaFilter int

const (
	MediaAll MediaFilter = iota
	MediaImages
	MediaVideos
	MediaAudio
)

// mediaFilterNames maps filter values to display names.
var mediaFilterNames = []string{"All", "Images", "Videos", "Audio"}

// MediaField is a filter field for media type selection.
type MediaField struct {
	name   string
	filter MediaFilter
}

// NewMediaField creates a new MediaField with the given display name.
func NewMediaField(name string) *MediaField {
	return &MediaField{
		name:   name,
		filter: MediaAll,
	}
}

// Name returns the display name of the field.
func (f *MediaField) Name() string {
	return f.name
}

// Value returns the current filter as a string.
func (f *MediaField) Value() string {
	if f.filter == MediaAll {
		return ""
	}
	return mediaFilterNames[f.filter]
}

// Reset clears the filter back to All.
func (f *MediaField) Reset() {
	f.filter = MediaAll
}

// Filter returns the current MediaFilter.
func (f *MediaField) Filter() MediaFilter {
	return f.filter
}

// Cycle advances to the next filter option.
func (f *MediaField) Cycle() {
	f.filter = (f.filter + 1) % MediaFilter(len(mediaFilterNames))
}

// View renders the field for display in the sidebar.
func (f *MediaField) View() string {
	label := lipgloss.NewStyle().Bold(true).Render(f.name)

	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	if f.filter != MediaAll {
		valStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	}

	indicator := valStyle.Render(fmt.Sprintf("[%s]", mediaFilterNames[f.filter]))
	return fmt.Sprintf("%s %s", label, indicator)
}
