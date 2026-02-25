// =============================================================================
// FILE: internal/filter/post_temp.go
// PURPOSE: Timed/temporary post filter. Filters posts by their expiry status,
//          either keeping or excluding posts that have an expiration date.
//          Ports Python filters temp_post_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Temporary post filter
// ---------------------------------------------------------------------------

// ByTempPost returns a filter based on temporary/expiring post status.
//
// Parameters:
//   - mode: "only" keeps only expiring posts, "exclude" removes them, "" = no filter.
//
// Returns:
//   - A PostFilter, or nil if mode is empty.
func ByTempPost(mode string) PostFilter {
	if mode == "" {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		var result []model.Post
		for _, p := range posts {
			switch mode {
			case "only":
				if p.HasExpiry {
					result = append(result, p)
				}
			case "exclude":
				if !p.HasExpiry {
					result = append(result, p)
				}
			default:
				result = append(result, p)
			}
		}
		return result
	}
}
