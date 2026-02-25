// =============================================================================
// FILE: internal/filter/media_date.go
// PURPOSE: Media date range filter. Filters media items by their creation date,
//          keeping only those within the configured before/after date range.
//          Ports Python filters/media/filters.py posts_date_filter_media.
// =============================================================================

package filter

import (
	"time"

	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Media date filter
// ---------------------------------------------------------------------------

// ByMediaDate returns a filter that keeps only media within the given date range.
//
// Parameters:
//   - after: Only keep media after this time (zero = no lower bound).
//   - before: Only keep media before this time (zero = no upper bound).
//
// Returns:
//   - A MediaFilter that removes out-of-range media.
func ByMediaDate(after, before time.Time) MediaFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			dateStr := m.Date()
			if dateStr == "" {
				// No date — include by default.
				result = append(result, m)
				continue
			}

			t, err := utils.ParseFlexibleDate(dateStr)
			if err != nil {
				result = append(result, m)
				continue
			}

			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)

			if afterOK && beforeOK {
				result = append(result, m)
			}
		}
		return result
	}
}
