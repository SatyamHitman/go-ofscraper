// =============================================================================
// FILE: internal/prompts/groups/merge.go
// PURPOSE: Merge prompts. Presents database merge options and confirmation.
//          Ports Python prompts/prompt_groups/merge.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Merge prompt
// ---------------------------------------------------------------------------

// PromptMergeConfirm asks the user to confirm a database merge operation.
//
// Parameters:
//   - sourcePath: Path of the source database.
//   - targetPath: Path of the target database.
//
// Returns:
//   - true if confirmed, false otherwise.
func PromptMergeConfirm(sourcePath, targetPath string) (bool, error) {
	fmt.Printf("\nMerge database:\n  Source: %s\n  Target: %s\n", sourcePath, targetPath)
	fmt.Print("\nConfirm merge? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, nil
	}

	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
}
