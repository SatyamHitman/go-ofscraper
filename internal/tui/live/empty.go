// =============================================================================
// FILE: internal/tui/live/empty.go
// PURPOSE: Empty display placeholder. Shown when no data is available for the
//          current view (e.g., no active downloads, no tasks).
//          Ports Python utils/live/empty.py.
// =============================================================================

package live

import (
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// EmptyDisplay
// ---------------------------------------------------------------------------

// EmptyDisplay is a placeholder component shown when no data is available.
type EmptyDisplay struct {
	message string
	width   int
}

// NewEmptyDisplay creates a new EmptyDisplay with the given message.
func NewEmptyDisplay(message string) *EmptyDisplay {
	return &EmptyDisplay{
		message: message,
	}
}

// SetMessage sets the placeholder message.
func (e *EmptyDisplay) SetMessage(msg string) {
	e.message = msg
}

// SetWidth sets the display width.
func (e *EmptyDisplay) SetWidth(w int) {
	e.width = w
}

// Render returns the empty display as a styled string.
func (e *EmptyDisplay) Render() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Align(lipgloss.Center).
		Italic(true)

	if e.width > 0 {
		style = style.Width(e.width)
	}

	msg := e.message
	if msg == "" {
		msg = "No data available"
	}

	return style.Render(msg)
}
