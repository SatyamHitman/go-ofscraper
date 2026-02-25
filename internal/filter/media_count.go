// =============================================================================
// FILE: internal/filter/media_count.go
// PURPOSE: Count limit filter. Limits media to the first N items after other
//          filters have been applied. Ports Python
//          filters/media/filters.py ele_count_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Count limit filter
// ---------------------------------------------------------------------------

// ByCount returns a filter that keeps at most maxCount media items.
// A zero or negative maxCount means no limit.
//
// Parameters:
//   - maxCount: Maximum number of media items to keep.
//
// Returns:
//   - A MediaFilter, or nil if maxCount <= 0.
func ByCount(maxCount int) MediaFilter {
	if maxCount <= 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		if len(media) <= maxCount {
			return media
		}
		return media[:maxCount]
	}
}
