// =============================================================================
// FILE: internal/filter/post_neg_text.go
// PURPOSE: Negative text regex filter. Removes posts whose text matches any of
//          the given regex patterns. Ports Python filters post_neg_text_filter.
// =============================================================================

package filter

import (
	"regexp"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Negative text regex filter
// ---------------------------------------------------------------------------

// ByPostNegText returns a filter that removes posts matching any of the given
// regex patterns. The inverse of ByPostText.
//
// Parameters:
//   - patterns: Regex pattern strings. Empty slice = no filter.
//
// Returns:
//   - A PostFilter, or nil if no patterns given.
func ByPostNegText(patterns []string) PostFilter {
	if len(patterns) == 0 {
		return nil
	}

	var compiled []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		compiled = append(compiled, re)
	}

	if len(compiled) == 0 {
		return nil
	}

	return func(posts []model.Post) []model.Post {
		var result []model.Post
		for _, post := range posts {
			matched := false
			for _, re := range compiled {
				if re.MatchString(post.RawText) {
					matched = true
					break
				}
			}
			if !matched {
				result = append(result, post)
			}
		}
		return result
	}
}
