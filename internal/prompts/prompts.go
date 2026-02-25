// =============================================================================
// FILE: internal/prompts/prompts.go
// PURPOSE: Prompt aggregator. Orchestrates the interactive TUI prompt flow,
//          managing the sequence of user selections from main menu through
//          model selection. Ports Python prompts/prompts.py.
// =============================================================================

package prompts

import (
	"fmt"

	"gofscraper/internal/prompts/groups"
)

// ---------------------------------------------------------------------------
// Prompt flow orchestrator
// ---------------------------------------------------------------------------

// Result holds the collected prompt selections.
type Result struct {
	Action   string   // Selected action (e.g. "download", "like", "unlike")
	Areas    []string // Selected content areas (e.g. "timeline", "messages")
	Users    []string // Selected usernames
	Profile  string   // Selected profile name
	Continue bool     // Whether the user chose to continue
}

// RunMainMenu presents the main interactive menu and collects all selections.
// Returns a Result with the user's choices, or an error if cancelled.
//
// Returns:
//   - The collected Result, or error if the user exits.
func RunMainMenu() (*Result, error) {
	// Step 1: Action selection
	action, err := groups.PromptAction()
	if err != nil {
		return nil, fmt.Errorf("action prompt: %w", err)
	}

	// Step 2: Area selection
	areas, err := groups.PromptAreas()
	if err != nil {
		return nil, fmt.Errorf("area prompt: %w", err)
	}

	return &Result{
		Action:   action,
		Areas:    areas,
		Continue: true,
	}, nil
}

// RunProfileMenu presents the profile selection menu.
//
// Parameters:
//   - profiles: Available profile names.
//
// Returns:
//   - The selected profile name, or error.
func RunProfileMenu(profiles []string) (string, error) {
	return groups.PromptProfile(profiles)
}

// RunConfigMenu presents the config editing menu.
//
// Returns:
//   - Error if cancelled.
func RunConfigMenu() error {
	return groups.PromptConfig()
}
