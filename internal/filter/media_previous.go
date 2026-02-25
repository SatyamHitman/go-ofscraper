// =============================================================================
// FILE: internal/filter/media_previous.go
// PURPOSE: Previous download filter. Removes media items that have already
//          been downloaded (based on a set of known media IDs).
//          Ports Python filters/media/filters.py previous_download_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Previous download filter
// ---------------------------------------------------------------------------

// ByPreviousDownload returns a filter that removes media already downloaded.
// The downloadedIDs set contains media IDs that should be skipped.
//
// Parameters:
//   - downloadedIDs: Set of already-downloaded media IDs.
//   - skipPrevious: Whether to actually filter. If false, returns nil.
//
// Returns:
//   - A MediaFilter, or nil if skipPrevious is false or set is empty.
func ByPreviousDownload(downloadedIDs map[int64]bool, skipPrevious bool) MediaFilter {
	if !skipPrevious || len(downloadedIDs) == 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			if !downloadedIDs[m.ID] {
				result = append(result, m)
			}
		}
		return result
	}
}
