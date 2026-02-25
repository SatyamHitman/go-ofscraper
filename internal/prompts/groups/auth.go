// =============================================================================
// FILE: internal/prompts/groups/auth.go
// PURPOSE: Auth prompts. Presents authentication setup prompts for cookie,
//          user agent, and other auth values.
//          Ports Python prompts/prompt_groups/auth.py.
// =============================================================================

package groups

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Auth prompts
// ---------------------------------------------------------------------------

// PromptAuthCookie prompts the user for their session cookie.
//
// Returns:
//   - The entered cookie string, or error.
func PromptAuthCookie() (string, error) {
	fmt.Print("Enter session cookie (sess): ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading cookie: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// PromptAuthUserAgent prompts the user for their browser user agent.
//
// Returns:
//   - The entered user agent string, or error.
func PromptAuthUserAgent() (string, error) {
	fmt.Print("Enter user agent: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading user agent: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// PromptAuthXBC prompts the user for their X-BC token.
//
// Returns:
//   - The entered X-BC string, or error.
func PromptAuthXBC() (string, error) {
	fmt.Print("Enter X-BC token: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading X-BC: %w", err)
	}
	return strings.TrimSpace(line), nil
}
