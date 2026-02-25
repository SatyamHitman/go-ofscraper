// =============================================================================
// FILE: internal/filter/media_dupe.go
// PURPOSE: Duplicate media filter. Removes media items with duplicate IDs,
//          keeping only the first occurrence. Ports Python
//          filters/media/filters.py dupefiltermedia.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Duplicate media filter
// ---------------------------------------------------------------------------

// ByMediaDupe returns a filter that removes duplicate media by ID.
// The first occurrence of each ID is kept.
//
// Returns:
//   - A MediaFilter.
func ByMediaDupe() MediaFilter {
	return func(media []*model.Media) []*model.Media {
		seen := make(map[int64]bool)
		var result []*model.Media
		for _, m := range media {
			if !seen[m.ID] {
				seen[m.ID] = true
				result = append(result, m)
			}
		}
		return result
	}
}
