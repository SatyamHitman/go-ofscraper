// =============================================================================
// FILE: internal/filter/post_mass.go
// PURPOSE: Mass message filter. Filters posts by their mass message status.
//          Ports Python filters mass_msg_filter.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Mass message filter
// ---------------------------------------------------------------------------

// ByMassMessage returns a filter based on mass message status.
//
// Parameters:
//   - mode: "only" keeps only mass messages, "exclude" removes them, "" = no filter.
//
// Returns:
//   - A PostFilter, or nil if mode is empty.
func ByMassMessage(mode string) PostFilter {
	if mode == "" {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		var result []model.Post
		for _, p := range posts {
			switch mode {
			case "only":
				if p.Mass {
					result = append(result, p)
				}
			case "exclude":
				if !p.Mass {
					result = append(result, p)
				}
			default:
				result = append(result, p)
			}
		}
		return result
	}
}
