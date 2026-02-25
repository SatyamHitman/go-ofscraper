// =============================================================================
// FILE: internal/prompts/validators.go
// PURPOSE: Input validators for interactive prompts. Provides validation
//          functions for usernames, dates, numbers, and paths.
//          Ports Python prompts/prompt_validators.py.
// =============================================================================

package prompts

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Validators
// ---------------------------------------------------------------------------

// ValidateUsername checks if a username string is valid.
// Usernames must be alphanumeric with underscores, 1-50 chars.
//
// Parameters:
//   - username: The username to validate.
//
// Returns:
//   - nil if valid, error describing the issue otherwise.
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > 50 {
		return fmt.Errorf("username too long (max 50 characters)")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

// ValidateDate checks if a date string can be parsed.
//
// Parameters:
//   - date: The date string to validate.
//
// Returns:
//   - nil if valid, error otherwise.
func ValidateDate(date string) error {
	date = strings.TrimSpace(date)
	if date == "" {
		return nil // Empty is allowed (means no constraint).
	}
	_, err := utils.ParseFlexibleDate(date)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", date)
	}
	return nil
}

// ValidatePositiveInt checks if a string represents a positive integer.
//
// Parameters:
//   - s: The string to validate.
//
// Returns:
//   - nil if valid, error otherwise.
func ValidatePositiveInt(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number: %s", s)
	}
	if n < 0 {
		return fmt.Errorf("must be a positive number")
	}
	return nil
}

// ValidateNonEmpty checks that the string is not empty after trimming.
//
// Parameters:
//   - s: The string to validate.
//
// Returns:
//   - nil if non-empty, error otherwise.
func ValidateNonEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

// ValidateChoice checks that the value is one of the allowed choices.
//
// Parameters:
//   - value: The value to check.
//   - choices: Allowed values.
//
// Returns:
//   - nil if value is in choices, error otherwise.
func ValidateChoice(value string, choices []string) error {
	for _, c := range choices {
		if strings.EqualFold(value, c) {
			return nil
		}
	}
	return fmt.Errorf("invalid choice %q, must be one of: %s", value, strings.Join(choices, ", "))
}
