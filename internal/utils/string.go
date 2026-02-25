// =============================================================================
// FILE: internal/utils/string.go
// PURPOSE: String utility functions. Provides helpers for string manipulation,
//          truncation, formatting, and common transformations used throughout
//          the application. Ports Python utils/string.py.
// =============================================================================

package utils

import (
	"strings"
	"unicode/utf8"
)

// ---------------------------------------------------------------------------
// String manipulation
// ---------------------------------------------------------------------------

// Truncate shortens a string to at most maxLen runes, appending "..." if
// truncated. Operates on rune boundaries to avoid splitting multi-byte chars.
//
// Parameters:
//   - s: The string to truncate.
//   - maxLen: Maximum rune count (including the "..." suffix).
//
// Returns:
//   - The possibly truncated string.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// TruncateBytes shortens a string to at most maxBytes bytes, appending "..."
// if truncated. Cuts at valid UTF-8 boundaries.
//
// Parameters:
//   - s: The string to truncate.
//   - maxBytes: Maximum byte length.
//
// Returns:
//   - The possibly truncated string.
func TruncateBytes(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	if maxBytes <= 3 {
		return s[:maxBytes]
	}
	// Walk backwards to find a valid rune boundary.
	end := maxBytes - 3
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	return s[:end] + "..."
}

// NormalizeWhitespace collapses all runs of whitespace (spaces, tabs, newlines)
// into single spaces and trims leading/trailing whitespace.
//
// Parameters:
//   - s: The input string.
//
// Returns:
//   - The normalised string.
func NormalizeWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// RemovePrefix removes the given prefix from s if present.
//
// Parameters:
//   - s: The string.
//   - prefix: The prefix to remove.
//
// Returns:
//   - s without the prefix (or s unchanged if prefix not found).
func RemovePrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// RemoveSuffix removes the given suffix from s if present.
//
// Parameters:
//   - s: The string.
//   - suffix: The suffix to remove.
//
// Returns:
//   - s without the suffix.
func RemoveSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// ContainsAny reports whether s contains any of the given substrings.
//
// Parameters:
//   - s: The string to search.
//   - subs: The substrings to look for.
//
// Returns:
//   - true if any substring is found in s.
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAnyFold is like ContainsAny but case-insensitive.
//
// Parameters:
//   - s: The string to search (compared case-insensitively).
//   - subs: Substrings to look for.
//
// Returns:
//   - true if any substring matches case-insensitively.
func ContainsAnyFold(s string, subs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range subs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// PadRight pads a string on the right to the desired width with spaces.
//
// Parameters:
//   - s: The string to pad.
//   - width: Desired minimum width.
//
// Returns:
//   - The padded string.
func PadRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
