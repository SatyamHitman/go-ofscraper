// =============================================================================
// FILE: internal/cli/callbacks/choice.go
// PURPOSE: Choice/enum validation callbacks for flag values.
// =============================================================================

package callbacks

import (
	"fmt"
	"strings"
)

// ValidateChoice checks that value is one of the allowed choices (case-insensitive).
func ValidateChoice(value string, choices []string) error {
	lower := strings.ToLower(strings.TrimSpace(value))
	for _, c := range choices {
		if strings.ToLower(c) == lower {
			return nil
		}
	}
	return fmt.Errorf("invalid value %q: must be one of [%s]", value, strings.Join(choices, ", "))
}

// ValidateAction validates an action flag value.
func ValidateAction(value string) error {
	return ValidateChoice(value, []string{"download", "like", "unlike"})
}

// ValidateQuality validates a quality flag value.
func ValidateQuality(value string) error {
	return ValidateChoice(value, []string{"source", "high", "medium", "low"})
}

// ValidateKeyMode validates a key-mode flag value.
func ValidateKeyMode(value string) error {
	return ValidateChoice(value, []string{"auto", "manual", "keydb"})
}

// ValidateSortType validates a sort-type flag value.
func ValidateSortType(value string) error {
	return ValidateChoice(value, []string{"name", "subscribed", "expired"})
}
