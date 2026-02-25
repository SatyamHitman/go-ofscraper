// =============================================================================
// FILE: internal/filter/post_date.go
// PURPOSE: Post date range filter. Filters posts by their creation date,
//          keeping only those within the configured before/after date range.
//          Ports Python filters post_date_filter.
// =============================================================================

package filter

import (
	"time"

	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Post date filter
// ---------------------------------------------------------------------------

// ByPostDate returns a filter that keeps posts within the given date range.
//
// Parameters:
//   - after: Only keep posts after this time (zero = no lower bound).
//   - before: Only keep posts before this time (zero = no upper bound).
//
// Returns:
//   - A PostFilter.
func ByPostDate(after, before time.Time) PostFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		var result []model.Post
		for _, p := range posts {
			t, err := utils.ParseFlexibleDate(p.CreatedAt)
			if err != nil {
				// Can't parse date — include by default.
				result = append(result, p)
				continue
			}

			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)

			if afterOK && beforeOK {
				result = append(result, p)
			}
		}
		return result
	}
}
