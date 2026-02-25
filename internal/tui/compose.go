// =============================================================================
// FILE: internal/tui/compose.go
// PURPOSE: Layout composition. Arranges TUI components (sidebar, table, status
//          bar) into the final screen layout. Ports Python classes/table/compose.py.
// =============================================================================

package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Layout composition
// ---------------------------------------------------------------------------

// ComposeLayout arranges sidebar and main content side by side.
//
// Parameters:
//   - sidebar: The sidebar content string.
//   - main: The main content string.
//   - width: Total available width.
//   - height: Total available height.
//
// Returns:
//   - The composed layout string.
func ComposeLayout(sidebar, main string, width, height int) string {
	sidebarWidth := SidebarWidth
	mainWidth := width - sidebarWidth - 1 // 1 for separator

	if mainWidth < MinWidth/2 {
		// Too narrow for sidebar — show only main.
		return main
	}

	sidebarBox := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(height).
		Render(sidebar)

	mainBox := lipgloss.NewStyle().
		Width(mainWidth).
		Height(height).
		Render(main)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, mainBox)
}

// ComposeVertical stacks header, body, and footer vertically.
//
// Parameters:
//   - header: Header content.
//   - body: Main body content.
//   - footer: Footer content.
//
// Returns:
//   - The composed layout string.
func ComposeVertical(header, body, footer string) string {
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}
