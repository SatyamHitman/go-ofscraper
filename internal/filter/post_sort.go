// =============================================================================
// FILE: internal/filter/post_sort.go
// PURPOSE: Post sort filter. Sorts posts by various criteria (date, ID, text,
//          price). Ports Python filters final_post_sort.
// =============================================================================

package filter

import (
	"sort"
	"strings"

	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Post sort
// ---------------------------------------------------------------------------

// SortPosts returns a filter that sorts posts by the given key and direction.
//
// Parameters:
//   - sortBy: Sort key ("date", "id", "text", "price", "random"). Empty = no sort.
//   - descending: If true, sort in descending order.
//
// Returns:
//   - A PostFilter, or nil if sortBy is empty.
func SortPosts(sortBy string, descending bool) PostFilter {
	if sortBy == "" {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		sorted := make([]model.Post, len(posts))
		copy(sorted, posts)

		sort.SliceStable(sorted, func(i, j int) bool {
			less := comparePosts(sorted[i], sorted[j], sortBy)
			if descending {
				return !less
			}
			return less
		})

		return sorted
	}
}

// comparePosts returns true if a should sort before b for the given key.
func comparePosts(a, b model.Post, key string) bool {
	switch strings.ToLower(key) {
	case "date":
		ta, errA := utils.ParseFlexibleDate(a.Date())
		tb, errB := utils.ParseFlexibleDate(b.Date())
		if errA != nil || errB != nil {
			return a.ID < b.ID
		}
		return ta.Before(tb)
	case "id":
		return a.ID < b.ID
	case "text":
		return strings.ToLower(a.RawText) < strings.ToLower(b.RawText)
	case "price":
		return a.Price < b.Price
	default:
		return a.ID < b.ID
	}
}
