// =============================================================================
// FILE: internal/prompts/utils/model_helpers.go
// PURPOSE: Model prompt helper functions. Provides formatting and display
//          helpers for model/user selection prompts.
//          Ports Python prompts/utils/model_helpers.py.
// =============================================================================

package promptutils

import (
	"fmt"
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Model display helpers
// ---------------------------------------------------------------------------

// FormatUserLine formats a user for display in selection lists.
//
// Parameters:
//   - u: The user to format.
//   - showPrice: Whether to show price info.
//
// Returns:
//   - Formatted display string.
func FormatUserLine(u model.User, showPrice bool) string {
	var parts []string
	parts = append(parts, u.Name)

	if u.IsActive() {
		parts = append(parts, "[active]")
	} else {
		parts = append(parts, "[expired]")
	}

	if showPrice {
		price := u.FinalCurrentPrice()
		if price == 0 {
			parts = append(parts, "(free)")
		} else {
			parts = append(parts, fmt.Sprintf("($%.2f)", price))
		}
	}

	return strings.Join(parts, " ")
}

// FormatUserList formats a slice of users for display.
//
// Parameters:
//   - users: The users to format.
//   - showPrice: Whether to show price info.
//
// Returns:
//   - Slice of formatted display strings.
func FormatUserList(users []model.User, showPrice bool) []string {
	lines := make([]string, len(users))
	for i, u := range users {
		lines[i] = FormatUserLine(u, showPrice)
	}
	return lines
}

// FilterUsersByName returns users whose names contain the search string
// (case-insensitive).
//
// Parameters:
//   - users: Users to filter.
//   - search: Search string.
//
// Returns:
//   - Matching users.
func FilterUsersByName(users []model.User, search string) []model.User {
	if search == "" {
		return users
	}
	search = strings.ToLower(search)
	var result []model.User
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.Name), search) {
			result = append(result, u)
		}
	}
	return result
}
