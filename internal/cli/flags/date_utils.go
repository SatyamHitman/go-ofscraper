// =============================================================================
// FILE: internal/cli/flags/date_utils.go
// PURPOSE: Date parsing utilities for flag values. Provides helpers that
//          convert CLI date strings into time.Time.
// =============================================================================

package flags

import (
	"fmt"
	"time"
)

// Supported date layouts tried in order when parsing flag values.
var dateLayouts = []string{
	"2006-01-02",
	"2006-01-02T15:04:05",
	"2006/01/02",
	"01-02-2006",
	"Jan 2, 2006",
}

// ParseDate attempts to parse a date string using the supported layouts.
// Returns the parsed time or an error if none of the layouts match.
func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date %q: expected format like 2006-01-02", value)
}

// ParseDateOrDefault parses a date string, returning the fallback value if the
// string is empty or unparseable.
func ParseDateOrDefault(value string, fallback time.Time) time.Time {
	t, err := ParseDate(value)
	if err != nil {
		return fallback
	}
	return t
}

// ParseEpoch converts a Unix epoch (seconds) to a time.Time.
// Returns the zero time if epoch is 0.
func ParseEpoch(epoch int64) time.Time {
	if epoch == 0 {
		return time.Time{}
	}
	return time.Unix(epoch, 0)
}
