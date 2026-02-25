// =============================================================================
// FILE: internal/utils/dates.go
// PURPOSE: Date/time utility functions. Provides parsing, formatting, and
//          comparison helpers for the various date formats used across the OF
//          API and config system. Ports Python utils/dates.py.
// =============================================================================

package utils

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Date format layouts
// ---------------------------------------------------------------------------

// Common date/time layouts used throughout the application.
const (
	LayoutISO8601     = "2006-01-02T15:04:05Z07:00"
	LayoutISO8601NoTZ = "2006-01-02T15:04:05"
	LayoutDateOnly    = "2006-01-02"
	LayoutDisplay     = "01-02-2006"
	LayoutDisplayTime = "01-02-2006 15:04:05"
	LayoutCompact     = "20060102"
	LayoutLog         = "2006-01-02 15:04:05"
)

// parseLayouts is the list of layouts tried in order by ParseFlexibleDate.
var parseLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	LayoutISO8601,
	LayoutISO8601NoTZ,
	LayoutDateOnly,
	LayoutDisplay,
	LayoutCompact,
}

// ---------------------------------------------------------------------------
// Parsing
// ---------------------------------------------------------------------------

// ParseFlexibleDate tries multiple date layouts and returns the first match.
// Returns a zero time and error if none match.
//
// Parameters:
//   - s: The date string to parse.
//
// Returns:
//   - The parsed time.Time, and any error.
func ParseFlexibleDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}
	for _, layout := range parseLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %q", s)
}

// ParseUnixTimestamp converts a Unix timestamp (seconds since epoch) to time.
//
// Parameters:
//   - ts: Unix timestamp as float64 (supports fractional seconds).
//
// Returns:
//   - The corresponding time.Time in UTC.
func ParseUnixTimestamp(ts float64) time.Time {
	sec := int64(ts)
	nsec := int64((ts - float64(sec)) * 1e9)
	return time.Unix(sec, nsec).UTC()
}

// ---------------------------------------------------------------------------
// Formatting
// ---------------------------------------------------------------------------

// FormatDate formats a time using the given layout string. If the layout is
// empty, uses the default display format.
//
// Parameters:
//   - t: The time to format.
//   - layout: A Go time layout string, or "" for default.
//
// Returns:
//   - The formatted date string.
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = LayoutDisplay
	}
	return t.Format(layout)
}

// FormatDateISO formats a time as ISO 8601 (RFC 3339).
//
// Parameters:
//   - t: The time to format.
//
// Returns:
//   - The ISO 8601 formatted string.
func FormatDateISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ---------------------------------------------------------------------------
// Comparison helpers
// ---------------------------------------------------------------------------

// IsAfter reports whether t is strictly after the reference time.
//
// Parameters:
//   - t: The time to check.
//   - ref: The reference time.
//
// Returns:
//   - true if t is after ref.
func IsAfter(t, ref time.Time) bool {
	return t.After(ref)
}

// IsBefore reports whether t is strictly before the reference time.
//
// Parameters:
//   - t: The time to check.
//   - ref: The reference time.
//
// Returns:
//   - true if t is before ref.
func IsBefore(t, ref time.Time) bool {
	return t.Before(ref)
}

// IsBetween reports whether t falls within the inclusive range [start, end].
//
// Parameters:
//   - t: The time to check.
//   - start: Range start (inclusive).
//   - end: Range end (inclusive).
//
// Returns:
//   - true if start <= t <= end.
func IsBetween(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

// NowUTC returns the current time in UTC.
//
// Returns:
//   - Current time.Time in UTC.
func NowUTC() time.Time {
	return time.Now().UTC()
}
