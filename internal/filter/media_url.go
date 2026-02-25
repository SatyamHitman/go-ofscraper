// =============================================================================
// FILE: internal/filter/media_url.go
// PURPOSE: URL presence filter. Removes media items that have no download URL.
//          Ports Python filters/media/filters.py url_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// URL presence filter
// ---------------------------------------------------------------------------

// ByURLPresence returns a filter that removes media lacking any download URL
// (neither direct nor MPD/DRM).
//
// Returns:
//   - A MediaFilter that keeps only linked media.
func ByURLPresence() MediaFilter {
	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			if m.IsLinked() {
				result = append(result, m)
			}
		}
		return result
	}
}
