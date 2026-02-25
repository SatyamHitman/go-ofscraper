// =============================================================================
// FILE: internal/prompts/groups/config.go
// PURPOSE: Config editing prompts. Presents configuration options for the user
//          to modify settings interactively.
//          Ports Python prompts/prompt_groups/config.py.
// =============================================================================

package groups

import (
	"fmt"
)

// ---------------------------------------------------------------------------
// Config prompts
// ---------------------------------------------------------------------------

// PromptConfig presents config editing options.
//
// Returns:
//   - Error if cancelled.
func PromptConfig() error {
	fmt.Println("\nConfiguration editing is available via:")
	fmt.Println("  1) Edit config file directly")
	fmt.Println("  2) Use CLI flags to override settings")
	fmt.Println("\nConfig file location: ~/.config/gofscraper/config.json")
	return nil
}
