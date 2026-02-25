// =============================================================================
// FILE: internal/cli/callbacks/username.go
// PURPOSE: Username parsing and validation callbacks.
// =============================================================================

package callbacks

import (
	"fmt"
	"strings"
)

// ValidateUsername checks that a single username string is non-empty and does
// not contain disallowed characters.
func ValidateUsername(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("username must not be empty")
	}
	if strings.ContainsAny(name, " \t\n") {
		return fmt.Errorf("username %q must not contain whitespace", name)
	}
	return nil
}

// ValidateUsernames validates a slice of usernames.
func ValidateUsernames(names []string) error {
	for _, n := range names {
		if err := ValidateUsername(n); err != nil {
			return err
		}
	}
	return nil
}
