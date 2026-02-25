// =============================================================================
// FILE: internal/filter/post_ad.go
// PURPOSE: Ad/promotional post filter. Removes posts that appear to be
//          advertisements based on text pattern matching.
//          Ports Python filters + utils/ads.py.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Ad post filter
// ---------------------------------------------------------------------------

// ByAdPost returns a filter that removes posts identified as advertisements.
// Uses the ad detection patterns from utils.IsAdPost.
//
// Parameters:
//   - filterAds: Whether to filter ads. If false, returns nil.
//
// Returns:
//   - A PostFilter, or nil if filterAds is false.
func ByAdPost(filterAds bool) PostFilter {
	if !filterAds {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		var result []model.Post
		for _, p := range posts {
			if !utils.IsAdPost(p.RawText) {
				result = append(result, p)
			}
		}
		return result
	}
}
