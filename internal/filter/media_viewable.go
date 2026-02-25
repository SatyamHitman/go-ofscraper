// =============================================================================
// FILE: internal/filter/media_viewable.go
// PURPOSE: Viewable media filter. Removes media items that the user does not
//          have permission to view. Ports Python
//          filters/media/filters.py unviewable_media_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Viewable filter
// ---------------------------------------------------------------------------

// ByViewable returns a filter that keeps only media the user can view.
// If skipUnviewable is false, no filtering is performed.
//
// Parameters:
//   - skipUnviewable: Whether to remove unviewable media.
//
// Returns:
//   - A MediaFilter, or nil if skipUnviewable is false.
func ByViewable(skipUnviewable bool) MediaFilter {
	if !skipUnviewable {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			if m.CanView {
				result = append(result, m)
			}
		}
		return result
	}
}
