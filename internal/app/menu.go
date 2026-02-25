// =============================================================================
// FILE: internal/app/menu.go
// PURPOSE: Interactive menu system. Presents the main menu with choices for
//          scraper, config editing, profile management, auth setup, and exit.
//          Dispatches to the appropriate handler. Ports Python prompts/prompt.py.
// =============================================================================

package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Menu choices
// ---------------------------------------------------------------------------

// MenuChoice represents a selectable menu option.
type MenuChoice struct {
	Key         string
	Label       string
	Description string
}

// MainMenuChoices returns the available main menu options.
func MainMenuChoices() []MenuChoice {
	return []MenuChoice{
		{Key: "1", Label: "Scraper", Description: "Run the scraper to download content"},
		{Key: "2", Label: "Config", Description: "Edit configuration settings"},
		{Key: "3", Label: "Profile", Description: "Manage configuration profiles"},
		{Key: "4", Label: "Auth", Description: "Setup or update authentication"},
		{Key: "5", Label: "Exit", Description: "Exit the application"},
	}
}

// ---------------------------------------------------------------------------
// RunMenu
// ---------------------------------------------------------------------------

// RunMenu presents the interactive main menu and dispatches user selections.
// Loops until the user selects exit or the context is cancelled.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - Error if a dispatched operation fails fatally.
func (a *App) RunMenu(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		if ctx.Err() != nil {
			return nil
		}

		a.printMenu()

		fmt.Print("\nSelect an option: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice := strings.TrimSpace(input)

		if err := a.dispatchMenu(ctx, choice); err != nil {
			a.logger.Error("menu action failed", "choice", choice, "error", err)
		}

		// Check for exit.
		if choice == "5" {
			return nil
		}
	}
}

// printMenu displays the main menu options.
func (a *App) printMenu() {
	fmt.Println()
	fmt.Println("========== gofscraper ==========")
	for _, choice := range MainMenuChoices() {
		fmt.Printf("  [%s] %-12s %s\n", choice.Key, choice.Label, choice.Description)
	}
	fmt.Println("================================")
}

// dispatchMenu routes the user's menu selection to the appropriate handler.
func (a *App) dispatchMenu(ctx context.Context, choice string) error {
	switch choice {
	case "1":
		return a.menuScraper(ctx)
	case "2":
		return a.menuConfig()
	case "3":
		return a.menuProfile()
	case "4":
		return a.menuAuth()
	case "5":
		return nil // Exit handled by caller.
	default:
		fmt.Printf("Unknown option: %s\n", choice)
		return nil
	}
}

// menuScraper runs the scraper from the interactive menu.
func (a *App) menuScraper(ctx context.Context) error {
	a.logger.Info("starting scraper from menu")
	// TODO: Wire to scraper command with interactive area/action selection.
	return a.RunAction(ctx, "download", []string{"timeline"}, nil)
}

// menuConfig opens the configuration editor.
func (a *App) menuConfig() error {
	a.logger.Info("opening configuration editor")
	// TODO: Wire to interactive config editor TUI.
	fmt.Println("Configuration editor not yet implemented.")
	return nil
}

// menuProfile opens the profile manager.
func (a *App) menuProfile() error {
	a.logger.Info("opening profile manager")
	// TODO: Wire to profile manager TUI.
	fmt.Println("Profile manager not yet implemented.")
	return nil
}

// menuAuth starts the authentication setup flow.
func (a *App) menuAuth() error {
	a.logger.Info("starting auth setup")
	// TODO: Wire to auth setup flow.
	fmt.Println("Auth setup not yet implemented.")
	return nil
}
