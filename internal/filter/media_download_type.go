// =============================================================================
// FILE: internal/filter/media_download_type.go
// PURPOSE: Download type filter. Separates media into normal (direct HTTP) and
//          protected (DRM/DASH) categories. Ports Python
//          filters/media/filters.py download_type_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Download type filter
// ---------------------------------------------------------------------------

// ByDownloadType returns a filter that keeps only media matching the specified
// download type (normal or protected).
//
// Parameters:
//   - dtype: The download type to keep ("normal" or "protected").
//
// Returns:
//   - A MediaFilter, or nil if dtype is empty (pass all).
func ByDownloadType(dtype string) MediaFilter {
	if dtype == "" {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			switch dtype {
			case "protected":
				if m.IsProtected() {
					result = append(result, m)
				}
			case "normal":
				if !m.IsProtected() {
					result = append(result, m)
				}
			default:
				result = append(result, m)
			}
		}
		return result
	}
}
