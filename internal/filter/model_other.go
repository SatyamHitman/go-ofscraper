// =============================================================================
// FILE: internal/filter/model_other.go
// PURPOSE: Miscellaneous model filters. Includes username include/exclude
//          lists, user list filtering, and promo availability.
//          Ports Python filters/models/other.py.
// =============================================================================

package filter

import (
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Username include filter
// ---------------------------------------------------------------------------

// ByUsernameInclude returns a filter that keeps only users whose name matches
// one of the given usernames (case-insensitive).
//
// Parameters:
//   - usernames: List of usernames to include. Empty = no filter.
//
// Returns:
//   - A ModelFilter, or nil if list is empty.
func ByUsernameInclude(usernames []string) ModelFilter {
	if len(usernames) == 0 {
		return nil
	}

	allowed := make(map[string]bool)
	for _, u := range usernames {
		allowed[strings.ToLower(u)] = true
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if allowed[strings.ToLower(u.Name)] {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Username exclude filter
// ---------------------------------------------------------------------------

// ByUsernameExclude returns a filter that removes users whose name matches
// one of the given usernames (case-insensitive).
//
// Parameters:
//   - usernames: List of usernames to exclude. Empty = no filter.
//
// Returns:
//   - A ModelFilter, or nil if list is empty.
func ByUsernameExclude(usernames []string) ModelFilter {
	if len(usernames) == 0 {
		return nil
	}

	blocked := make(map[string]bool)
	for _, u := range usernames {
		blocked[strings.ToLower(u)] = true
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if !blocked[strings.ToLower(u.Name)] {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// User list filter
// ---------------------------------------------------------------------------

// ByUserList returns a filter that keeps only users whose ID is in the given
// list of allowed IDs.
//
// Parameters:
//   - allowedIDs: Set of user IDs to keep. Empty = no filter.
//
// Returns:
//   - A ModelFilter, or nil if allowedIDs is empty.
func ByUserList(allowedIDs map[int64]bool) ModelFilter {
	if len(allowedIDs) == 0 {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if allowedIDs[u.ID] {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Has promo filter
// ---------------------------------------------------------------------------

// ByHasPromo returns a filter based on promo availability.
//
// Parameters:
//   - mode: "with_promo" keeps users with promos, "without_promo" keeps those without, "" = no filter.
//
// Returns:
//   - A ModelFilter, or nil if mode is empty.
func ByHasPromo(mode string) ModelFilter {
	if mode == "" {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			hasPromo := len(u.Promos) > 0
			switch mode {
			case "with_promo":
				if hasPromo {
					result = append(result, u)
				}
			case "without_promo":
				if !hasPromo {
					result = append(result, u)
				}
			default:
				result = append(result, u)
			}
		}
		return result
	}
}
