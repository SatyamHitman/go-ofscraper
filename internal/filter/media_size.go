// =============================================================================
// FILE: internal/filter/media_size.go
// PURPOSE: File size filter. Filters media by their file size (when known),
//          keeping only those within configured min/max byte ranges.
//          Ports download skip logic from Python.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Size filter
// ---------------------------------------------------------------------------

// BySize returns a filter that keeps only media within the given file size range.
// Media with unknown size (0) is included by default.
//
// Parameters:
//   - minBytes: Minimum size in bytes (0 = no lower bound).
//   - maxBytes: Maximum size in bytes (0 = no upper bound).
//
// Returns:
//   - A MediaFilter, or nil if both bounds are zero.
func BySize(minBytes, maxBytes float64) MediaFilter {
	if minBytes == 0 && maxBytes == 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			if m.Size == 0 {
				// Unknown size — include by default.
				result = append(result, m)
				continue
			}

			minOK := minBytes == 0 || m.Size >= minBytes
			maxOK := maxBytes == 0 || m.Size <= maxBytes

			if minOK && maxOK {
				result = append(result, m)
			}
		}
		return result
	}
}
