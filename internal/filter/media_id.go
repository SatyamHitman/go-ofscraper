// =============================================================================
// FILE: internal/filter/media_id.go
// PURPOSE: Media ID filter. Keeps only media whose IDs fall within a specified
//          numeric range. Ports Python filters/media/filters.py media_id_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Media ID filter
// ---------------------------------------------------------------------------

// ByMediaID returns a filter that keeps only media with IDs in [minID, maxID].
// A zero value for either bound means no bound on that side.
//
// Parameters:
//   - minID: Minimum media ID (0 = no lower bound).
//   - maxID: Maximum media ID (0 = no upper bound).
//
// Returns:
//   - A MediaFilter, or nil if both bounds are zero.
func ByMediaID(minID, maxID int64) MediaFilter {
	if minID == 0 && maxID == 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			minOK := minID == 0 || m.ID >= minID
			maxOK := maxID == 0 || m.ID <= maxID

			if minOK && maxOK {
				result = append(result, m)
			}
		}
		return result
	}
}
