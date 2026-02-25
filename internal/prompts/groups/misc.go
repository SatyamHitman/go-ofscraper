// =============================================================================
// FILE: internal/prompts/groups/misc.go
// PURPOSE: Miscellaneous prompts. Includes confirmation dialogs, yes/no
//          prompts, and utility prompts. Ports Python prompts/prompt_groups/misc.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Misc prompts
// ---------------------------------------------------------------------------

// PromptConfirm presents a yes/no confirmation prompt.
//
// Parameters:
//   - message: The message to display.
//   - defaultYes: Whether the default is Yes (true) or No (false).
//
// Returns:
//   - true if confirmed, false otherwise.
func PromptConfirm(message string, defaultYes bool) (bool, error) {
	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}

	fmt.Printf("%s %s: ", message, hint)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes, nil
	}

	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" {
		return defaultYes, nil
	}

	return line == "y" || line == "yes", nil
}

// PromptInput presents a text input prompt.
//
// Parameters:
//   - message: The prompt message.
//   - defaultVal: Default value if empty input.
//
// Returns:
//   - The entered string, or default.
func PromptInput(message, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", message, defaultVal)
	} else {
		fmt.Printf("%s: ", message)
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal, nil
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal, nil
	}
	return line, nil
}
