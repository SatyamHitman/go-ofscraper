// =============================================================================
// FILE: internal/filter/model_sort.go
// PURPOSE: Model sort filter. Sorts users/models by various criteria (name,
//          price, date, activity). Ports Python filters/models/sort.py.
// =============================================================================

package filter

import (
	"sort"
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Model sort
// ---------------------------------------------------------------------------

// SortModels returns a filter that sorts users by the given key and direction.
//
// Parameters:
//   - sortBy: Sort key ("name", "subscribed", "expired", "current-price",
//             "regular-price", "promo-price", "renewal-price", "last-seen").
//             Empty = no sort.
//   - descending: If true, sort in descending order.
//
// Returns:
//   - A ModelFilter, or nil if sortBy is empty.
func SortModels(sortBy string, descending bool) ModelFilter {
	if sortBy == "" {
		return nil
	}

	return func(users []model.User) []model.User {
		sorted := make([]model.User, len(users))
		copy(sorted, users)

		sort.SliceStable(sorted, func(i, j int) bool {
			less := compareModels(sorted[i], sorted[j], sortBy)
			if descending {
				return !less
			}
			return less
		})

		return sorted
	}
}

// compareModels returns true if a should sort before b for the given key.
func compareModels(a, b model.User, key string) bool {
	switch strings.ToLower(key) {
	case "name":
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	case "subscribed":
		return a.FinalSubscribed() < b.FinalSubscribed()
	case "expired":
		return a.FinalExpired() < b.FinalExpired()
	case "current-price":
		return a.FinalCurrentPrice() < b.FinalCurrentPrice()
	case "regular-price":
		return a.RegularPrice() < b.RegularPrice()
	case "promo-price":
		return a.FinalPromoPrice() < b.FinalPromoPrice()
	case "renewal-price":
		return a.FinalRenewalPrice() < b.FinalRenewalPrice()
	case "last-seen":
		return a.FinalLastSeen() < b.FinalLastSeen()
	default:
		return a.ID < b.ID
	}
}
