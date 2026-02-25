// =============================================================================
// FILE: internal/filter/post_text.go
// PURPOSE: Post text regex filter. Keeps only posts whose text matches one or
//          more regular expression patterns. Ports Python filters post_text_filter.
// =============================================================================

package filter

import (
	"regexp"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Text regex filter
// ---------------------------------------------------------------------------

// ByPostText returns a filter that keeps only posts matching at least one of
// the given regex patterns. Patterns are matched against the post's raw text.
//
// Parameters:
//   - patterns: Regex pattern strings. Empty slice = no filter.
//
// Returns:
//   - A PostFilter, or nil if no patterns given.
func ByPostText(patterns []string) PostFilter {
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
			for _, re := range compiled {
				if re.MatchString(post.RawText) {
					result = append(result, post)
					break
				}
			}
		}
		return result
	}
}
