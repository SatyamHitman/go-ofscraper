// =============================================================================
// FILE: internal/tui/styles.go
// PURPOSE: TUI styling definitions using lipgloss. Centralizes all visual
//          styling for the terminal UI. Ports Python classes/table/css.py.
// =============================================================================

package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Color palette
// ---------------------------------------------------------------------------

var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary = lipgloss.Color("#059669") // Green
	ColorAccent    = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#DC2626") // Red
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorBorder    = lipgloss.Color("#374151") // Dark gray
)

// ---------------------------------------------------------------------------
// Styles
// ---------------------------------------------------------------------------

var (
	// Title style for headers.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	// Subtle text for descriptions.
	SubtleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Active item in a list.
	ActiveStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// Inactive item in a list.
	InactiveStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Success status.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	// Error status.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	// Warning status.
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	// Border box.
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	// Status bar at the bottom.
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)
)
