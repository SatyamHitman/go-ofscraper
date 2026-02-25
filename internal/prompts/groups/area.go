// =============================================================================
// FILE: internal/prompts/groups/area.go
// PURPOSE: Area selection prompts. Presents content area choices (timeline,
//          messages, archived, etc.) for multi-select.
//          Ports Python prompts/prompt_groups/area.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// Area prompt
// ---------------------------------------------------------------------------

// PromptAreas presents content area choices and returns selected areas.
// The user can enter comma-separated numbers or "all".
//
// Returns:
//   - Slice of selected area names, or error if cancelled.
func PromptAreas() ([]string, error) {
	areas := []string{
		"timeline", "messages", "archived", "highlights",
		"stories", "pinned", "streams", "labels", "purchased",
	}

	fmt.Println("\nSelect Content Areas (comma-separated, or 'all'):")
	for i, a := range areas {
		fmt.Printf("  %d) %s\n", i+1, strings.Title(a))
	}
	fmt.Print("\nChoice [all]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return areas, nil
	}

	line = strings.TrimSpace(line)
	if line == "" || strings.EqualFold(line, "all") {
		return areas, nil
	}

	var selected []string
	for _, part := range strings.Split(line, ",") {
		part = strings.TrimSpace(part)
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 || idx > len(areas) {
			continue
		}
		selected = append(selected, areas[idx-1])
	}

	if len(selected) == 0 {
		return areas, nil
	}
	return selected, nil
}
