// =============================================================================
// FILE: internal/filter/post_dupe.go
// PURPOSE: Duplicate post filter. Removes posts with duplicate IDs, keeping
//          only the first occurrence. Ports Python filters dupefilterPost.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Duplicate post filter
// ---------------------------------------------------------------------------

// ByPostDupe returns a filter that removes duplicate posts by ID.
// The first occurrence of each ID is kept.
//
// Returns:
//   - A PostFilter.
func ByPostDupe() PostFilter {
	return func(posts []model.Post) []model.Post {
		seen := make(map[int64]bool)
		var result []model.Post
		for _, p := range posts {
			if !seen[p.ID] {
				seen[p.ID] = true
				result = append(result, p)
			}
		}
		return result
	}
}
