// =============================================================================
// FILE: internal/tui/tui.go
// PURPOSE: Bubbletea TUI application model. Main entry point for the terminal
//          user interface. Ports Python classes/table/app.py.
// =============================================================================

package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------------------------------------------------------
// App model
// ---------------------------------------------------------------------------

// App is the top-level Bubbletea model for the TUI.
type App struct {
	width   int
	height  int
	ready   bool
	quitting bool
	err     error
	view    View
}

// View represents which screen is currently active.
type View int

const (
	ViewMenu View = iota
	ViewTable
	ViewProgress
)

// New creates a new TUI App model.
func New() App {
	return App{
		view: ViewMenu,
	}
}

// Init implements tea.Model.
func (a App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			a.quitting = true
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
	}

	return a, nil
}

// View implements tea.Model.
func (a App) View() string {
	if a.quitting {
		return "Goodbye!\n"
	}

	if !a.ready {
		return "Initializing...\n"
	}

	if a.err != nil {
		return fmt.Sprintf("Error: %v\n", a.err)
	}

	switch a.view {
	case ViewMenu:
		return a.menuView()
	case ViewTable:
		return a.tableView()
	case ViewProgress:
		return a.progressView()
	default:
		return "Unknown view\n"
	}
}

// menuView renders the main menu.
func (a App) menuView() string {
	return "GoFScraper - Main Menu\n\nPress 'q' to quit\n"
}

// tableView renders the data table.
func (a App) tableView() string {
	return "Table View\n"
}

// progressView renders the download progress.
func (a App) progressView() string {
	return "Progress View\n"
}

// Run starts the TUI application.
func Run() error {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
