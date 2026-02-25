// =============================================================================
// FILE: internal/tui/constants.go
// PURPOSE: TUI constants. Defines sizing, spacing, and other constant values
//          for the terminal UI layout. Ports Python classes/table/const.py.
// =============================================================================

package tui

// ---------------------------------------------------------------------------
// Layout constants
// ---------------------------------------------------------------------------

const (
	MinWidth      = 80
	MinHeight     = 24
	SidebarWidth  = 30
	StatusHeight  = 1
	HeaderHeight  = 3
	FooterHeight  = 2
	PaddingH      = 2
	PaddingV      = 1
)

// ---------------------------------------------------------------------------
// Table constants
// ---------------------------------------------------------------------------

const (
	MaxColumnWidth  = 50
	MinColumnWidth  = 8
	DefaultPageSize = 20
)
