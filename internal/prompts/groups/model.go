// =============================================================================
// FILE: internal/prompts/groups/model.go
// PURPOSE: Model/user selection prompts. Presents creator usernames for the
//          user to select from. Ports Python prompts/prompt_groups/model.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Model selection prompt
// ---------------------------------------------------------------------------

// PromptModelSelect presents a list of models/creators for selection.
// Supports comma-separated numbers and "all".
//
// Parameters:
//   - users: Available users to choose from.
//
// Returns:
//   - Slice of selected usernames, or error.
func PromptModelSelect(users []model.User) ([]string, error) {
	if len(users) == 0 {
		return nil, fmt.Errorf("no users available")
	}

	fmt.Println("\nSelect Users (comma-separated, or 'all'):")
	for i, u := range users {
		status := "expired"
		if u.IsActive() {
			status = "active"
		}
		fmt.Printf("  %d) %s [%s]\n", i+1, u.Name, status)
	}
	fmt.Print("\nChoice [all]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return allUsernames(users), nil
	}

	line = strings.TrimSpace(line)
	if line == "" || strings.EqualFold(line, "all") {
		return allUsernames(users), nil
	}

	var selected []string
	for _, part := range strings.Split(line, ",") {
		part = strings.TrimSpace(part)
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 || idx > len(users) {
			continue
		}
		selected = append(selected, users[idx-1].Name)
	}

	if len(selected) == 0 {
		return allUsernames(users), nil
	}
	return selected, nil
}

// allUsernames extracts all usernames from a user slice.
func allUsernames(users []model.User) []string {
	names := make([]string, len(users))
	for i, u := range users {
		names[i] = u.Name
	}
	return names
}
