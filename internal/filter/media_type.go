// =============================================================================
// FILE: internal/filter/media_type.go
// PURPOSE: Media type filter. Filters media items by their content type
//          (photo, video, audio, gif) based on user-configured preferences.
//          Ports Python filters/media/filters.py mediatype_type_filter.
// =============================================================================

package filter

import (
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Media type filter
// ---------------------------------------------------------------------------

// ByMediaType returns a filter that keeps only media matching the allowed types.
//
// Parameters:
//   - allowedTypes: List of allowed type strings (e.g. "images", "videos", "audios").
//
// Returns:
//   - A MediaFilter that removes non-matching media.
func ByMediaType(allowedTypes []string) MediaFilter {
	if len(allowedTypes) == 0 {
		return nil // No filter = pass all.
	}

	// Build a lookup set.
	allowed := make(map[string]bool)
	for _, t := range allowedTypes {
		// Normalise: "images" -> "photo", "videos" -> "video", etc.
		switch strings.ToLower(t) {
		case "images", "image", "photo", "photos":
			allowed["photo"] = true
		case "videos", "video":
			allowed["video"] = true
		case "audios", "audio":
			allowed["audio"] = true
		case "gifs", "gif":
			allowed["gif"] = true
		default:
			allowed[strings.ToLower(t)] = true
		}
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			mediaType := strings.ToLower(m.Type)
			if allowed[mediaType] {
				result = append(result, m)
			}
		}
		return result
	}
}
