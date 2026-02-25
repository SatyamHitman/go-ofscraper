// =============================================================================
// FILE: internal/prompts/groups/profile.go
// PURPOSE: Profile selection prompts. Presents available profiles for the user
//          to choose from. Ports Python prompts/prompt_groups/profile.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Profile prompt
// ---------------------------------------------------------------------------

// PromptProfile presents profile choices and returns the selected profile name.
//
// Parameters:
//   - profiles: Available profile names.
//
// Returns:
//   - The selected profile name, or error.
func PromptProfile(profiles []string) (string, error) {
	if len(profiles) == 0 {
		return "default", nil
	}

	fmt.Println("\nSelect Profile:")
	for i, p := range profiles {
		fmt.Printf("  %d) %s\n", i+1, p)
	}
	fmt.Print("\nChoice [1]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return profiles[0], nil
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return profiles[0], nil
	}

	idx := 0
	_, err = fmt.Sscanf(line, "%d", &idx)
	if err != nil || idx < 1 || idx > len(profiles) {
		return profiles[0], nil
	}

	return profiles[idx-1], nil
}
