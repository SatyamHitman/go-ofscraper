// =============================================================================
// FILE: internal/filter/media_sort.go
// PURPOSE: Media sort filter. Sorts media by various criteria (date, ID, type,
//          size). Ports Python filters/media/filters.py final_media_sort.
// =============================================================================

package filter

import (
	"sort"
	"strings"

	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Media sort
// ---------------------------------------------------------------------------

// SortMedia returns a filter that sorts media by the given key and direction.
//
// Parameters:
//   - sortBy: Sort key ("date", "id", "type", "size", "random"). Empty = no sort.
//   - descending: If true, sort in descending order.
//
// Returns:
//   - A MediaFilter, or nil if sortBy is empty.
func SortMedia(sortBy string, descending bool) MediaFilter {
	if sortBy == "" {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		sorted := make([]*model.Media, len(media))
		copy(sorted, media)

		sort.SliceStable(sorted, func(i, j int) bool {
			less := compareMedia(sorted[i], sorted[j], sortBy)
			if descending {
				return !less
			}
			return less
		})

		return sorted
	}
}

// compareMedia returns true if a should sort before b for the given key.
func compareMedia(a, b *model.Media, key string) bool {
	switch strings.ToLower(key) {
	case "date":
		ta, errA := utils.ParseFlexibleDate(a.Date())
		tb, errB := utils.ParseFlexibleDate(b.Date())
		if errA != nil || errB != nil {
			return a.ID < b.ID
		}
		return ta.Before(tb)
	case "id":
		return a.ID < b.ID
	case "type":
		return a.Type < b.Type
	case "size":
		return a.Size < b.Size
	default:
		return a.ID < b.ID
	}
}
