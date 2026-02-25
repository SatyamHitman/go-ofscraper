// =============================================================================
// FILE: internal/filter/model_flags.go
// PURPOSE: Model flag-based filters. Filters users by boolean status flags
//          like active/expired, restricted, performer status.
//          Ports Python filters/models/flags.py.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Active subscription filter
// ---------------------------------------------------------------------------

// ByActiveStatus returns a filter based on subscription active status.
//
// Parameters:
//   - mode: "active" keeps only active, "expired" keeps only expired, "" = no filter.
//
// Returns:
//   - A ModelFilter, or nil if mode is empty.
func ByActiveStatus(mode string) ModelFilter {
	if mode == "" {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			switch mode {
			case "active":
				if u.IsActive() {
					result = append(result, u)
				}
			case "expired":
				if !u.IsActive() {
					result = append(result, u)
				}
			default:
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Restricted filter
// ---------------------------------------------------------------------------

// ByRestricted returns a filter based on restricted status.
//
// Parameters:
//   - mode: "restricted" keeps only restricted, "unrestricted" keeps only unrestricted, "" = no filter.
//
// Returns:
//   - A ModelFilter, or nil if mode is empty.
func ByRestricted(mode string) ModelFilter {
	if mode == "" {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			switch mode {
			case "restricted":
				if u.IsRestricted {
					result = append(result, u)
				}
			case "unrestricted":
				if !u.IsRestricted {
					result = append(result, u)
				}
			default:
				result = append(result, u)
			}
		}
		return result
	}
}
