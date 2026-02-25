// =============================================================================
// FILE: internal/filter/media_length.go
// PURPOSE: Media duration filter. Filters video/audio media by their duration,
//          keeping only those within configured min/max length bounds.
//          Ports Python filters/media/filters.py media_length_filter.
// =============================================================================

package filter

import (
	"strconv"
	"strings"
	"time"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Media length filter
// ---------------------------------------------------------------------------

// ByMediaLength returns a filter that keeps only media within the given
// duration range. Media without a parseable duration is included by default.
//
// Parameters:
//   - minLen: Minimum duration (zero = no lower bound).
//   - maxLen: Maximum duration (zero = no upper bound).
//
// Returns:
//   - A MediaFilter.
func ByMediaLength(minLen, maxLen time.Duration) MediaFilter {
	if minLen == 0 && maxLen == 0 {
		return nil
	}

	return func(media []*model.Media) []*model.Media {
		var result []*model.Media
		for _, m := range media {
			dur := parseDuration(m.Duration)
			if dur < 0 {
				// Unparseable — include by default.
				result = append(result, m)
				continue
			}

			minOK := minLen == 0 || dur >= minLen
			maxOK := maxLen == 0 || dur <= maxLen

			if minOK && maxOK {
				result = append(result, m)
			}
		}
		return result
	}
}

// parseDuration attempts to parse a duration string in formats like
// "123" (seconds), "1:23" (mm:ss), "1:23:45" (hh:mm:ss).
// Returns -1 if the string can't be parsed.
func parseDuration(s string) time.Duration {
	if s == "" {
		return -1
	}

	// Try plain seconds first.
	if secs, err := strconv.ParseFloat(s, 64); err == nil {
		return time.Duration(secs * float64(time.Second))
	}

	// Try HH:MM:SS or MM:SS.
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 2:
		mins, err1 := strconv.Atoi(parts[0])
		secs, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return -1
		}
		return time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second
	case 3:
		hrs, err1 := strconv.Atoi(parts[0])
		mins, err2 := strconv.Atoi(parts[1])
		secs, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return -1
		}
		return time.Duration(hrs)*time.Hour + time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second
	}

	return -1
}
