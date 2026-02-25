// =============================================================================
// FILE: internal/tui/live/panel.go
// PURPOSE: Bordered panel display component. Renders content inside a labeled
//          border box for visual grouping in the live display.
//          Ports Python utils/live/panel.py.
// =============================================================================

package live

import (
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Panel
// ---------------------------------------------------------------------------

// Panel renders content inside a bordered box with an optional title.
type Panel struct {
	Title       string
	Content     string
	Width       int
	BorderColor string
}

// NewPanel creates a new Panel with the given title.
func NewPanel(title string) *Panel {
	return &Panel{
		Title:       title,
		BorderColor: "#374151",
	}
}

// SetContent sets the panel's inner content.
func (p *Panel) SetContent(content string) {
	p.Content = content
}

// SetWidth sets the panel's width. A value of 0 uses automatic sizing.
func (p *Panel) SetWidth(w int) {
	p.Width = w
}

// SetBorderColor sets the border color as a hex string.
func (p *Panel) SetBorderColor(color string) {
	p.BorderColor = color
}

// Render returns the panel as a styled string.
func (p *Panel) Render() string {
	borderColor := lipgloss.Color(p.BorderColor)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	if p.Width > 0 {
		style = style.Width(p.Width)
	}

	content := p.Content
	if p.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED"))
		content = titleStyle.Render(p.Title) + "\n" + content
	}

	return style.Render(content)
}
