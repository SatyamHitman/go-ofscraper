// =============================================================================
// FILE: internal/prompts/groups/actions.go
// PURPOSE: Action selection prompts. Presents the user with available actions
//          (download, like, unlike, metadata) and returns the selection.
//          Ports Python prompts/prompt_groups/actions.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Action prompt
// ---------------------------------------------------------------------------

// PromptAction presents action choices to the user and returns the selection.
// In non-interactive mode, returns "download" as default.
//
// Returns:
//   - The selected action string, or error if cancelled.
func PromptAction() (string, error) {
	actions := []string{"download", "like", "unlike", "metadata"}

	fmt.Println("\nSelect Action:")
	for i, a := range actions {
		fmt.Printf("  %d) %s\n", i+1, strings.Title(a))
	}
	fmt.Print("\nChoice [1]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "download", nil
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return "download", nil
	}

	idx := 0
	_, err = fmt.Sscanf(line, "%d", &idx)
	if err != nil || idx < 1 || idx > len(actions) {
		return "download", nil
	}

	return actions[idx-1], nil
}
