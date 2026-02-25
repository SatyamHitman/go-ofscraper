// =============================================================================
// FILE: internal/prompts/groups/menu.go
// PURPOSE: Main menu prompts. Presents the top-level menu options and routes
//          to sub-menus. Ports Python prompts/prompt_groups/menu.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Menu options
// ---------------------------------------------------------------------------

const (
	MenuScraper  = "scraper"
	MenuConfig   = "config"
	MenuProfile  = "profile"
	MenuAuth     = "auth"
	MenuExit     = "exit"
)

// ---------------------------------------------------------------------------
// Main menu prompt
// ---------------------------------------------------------------------------

// PromptMainMenu presents the top-level menu and returns the selection.
//
// Returns:
//   - The selected menu option string, or error.
func PromptMainMenu() (string, error) {
	options := []struct {
		key   string
		label string
	}{
		{MenuScraper, "Run Scraper"},
		{MenuConfig, "Edit Configuration"},
		{MenuProfile, "Manage Profiles"},
		{MenuAuth, "Setup Authentication"},
		{MenuExit, "Exit"},
	}

	fmt.Println("\n=== GoFScraper ===")
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt.label)
	}
	fmt.Print("\nChoice [1]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return MenuExit, nil
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return MenuScraper, nil
	}

	idx := 0
	_, err = fmt.Sscanf(line, "%d", &idx)
	if err != nil || idx < 1 || idx > len(options) {
		return MenuScraper, nil
	}

	return options[idx-1].key, nil
}
