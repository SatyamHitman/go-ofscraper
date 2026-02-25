// =============================================================================
// FILE: internal/filter/media_post_id.go
// PURPOSE: Post ID filter for media. Keeps only media whose parent post ID
//          falls within a specified range. Ports Python
//          filters/media/filters.py post_id_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Post ID filter
// ---------------------------------------------------------------------------

// ByPostID returns a filter that keeps only media whose PostID is in [minID, maxID].
// A zero value for either bound means no bound on that side.
//
// Parameters:
//   - minID: Minimum post ID (0 = no lower bound).
//   - maxID: Maximum post ID (0 = no upper bound).
//
// Returns:
//   - A MediaFilter, or nil if both bounds are zero.
func ByPostID(minID, maxID int64) MediaFilter {
	if minID == 0 && maxID == 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			minOK := minID == 0 || m.PostID >= minID
			maxOK := maxID == 0 || m.PostID <= maxID

			if minOK && maxOK {
				result = append(result, m)
			}
		}
		return result
	}
}
