// =============================================================================
// FILE: internal/filter/model_subtype.go
// PURPOSE: Model subscription type filter. Filters users by their subscription
//          type (e.g., active, expired, all). Ports Python
//          filters/models/subtype.py.
// =============================================================================

package filter

import (
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Subscription type filter
// ---------------------------------------------------------------------------

// BySubType returns a filter that keeps users matching the subscription type.
//
// Parameters:
//   - subType: Subscription type to filter by ("active", "expired", "all").
//              Empty or "all" means no filter.
//
// Returns:
//   - A ModelFilter, or nil if subType is "all" or empty.
func BySubType(subType string) ModelFilter {
	subType = strings.ToLower(subType)
	if subType == "" || subType == "all" {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			switch subType {
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
