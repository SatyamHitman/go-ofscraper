// =============================================================================
// FILE: internal/model/base.go
// PURPOSE: Provides shared text processing utilities used across all domain
//          models. Handles text truncation for filenames, HTML cleanup for
//          database storage, and filename sanitization. Ports the Python
//          classes/of/base.py base class methods.
// =============================================================================

package model

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// htmlTagPattern matches HTML tags for removal during text cleanup.
var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)

// multiSpacePattern matches consecutive whitespace for normalization.
var multiSpacePattern = regexp.MustCompile(`\s+`)

// unsafeFilenamePattern matches characters unsafe for filenames across OS platforms.
var unsafeFilenamePattern = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// TextTruncate shortens text to the specified limit using the given truncation mode.
//
// Parameters:
//   - text: The input text to truncate.
//   - limit: Maximum length (characters for "letter" mode, words for "word" mode).
//     A limit of 0 means no truncation.
//   - mode: Truncation mode — "letter" truncates by character count,
//     "word" truncates by word count.
//
// Returns:
//   - The truncated text string.
func TextTruncate(text string, limit int, mode string) string {
	if limit <= 0 || text == "" {
		return text
	}

	switch mode {
	case "word":
		words := strings.Fields(text)
		if len(words) <= limit {
			return text
		}
		return strings.Join(words[:limit], " ")

	default: // "letter" mode
		if utf8.RuneCountInString(text) <= limit {
			return text
		}
		runes := []rune(text)
		return string(runes[:limit])
	}
}

// FileCleanup sanitizes text for use in file and directory names.
// Removes HTML tags, strips unsafe filename characters, and normalizes whitespace.
//
// Parameters:
//   - text: The raw text to sanitize for filesystem use.
//
// Returns:
//   - A sanitized string safe for use in file paths.
func FileCleanup(text string) string {
	if text == "" {
		return text
	}

	// Strip HTML tags
	result := htmlTagPattern.ReplaceAllString(text, "")

	// Remove characters unsafe for filenames
	result = unsafeFilenamePattern.ReplaceAllString(result, "")

	// Normalize whitespace (collapse multiple spaces, trim edges)
	result = multiSpacePattern.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}

// DBCleanup sanitizes text for database storage. Removes HTML tags and
// normalizes whitespace without stripping filename-unsafe characters.
//
// Parameters:
//   - text: The raw text to clean for database insertion.
//
// Returns:
//   - A sanitized string suitable for DB storage.
func DBCleanup(text string) string {
	if text == "" {
		return text
	}

	// Strip HTML tags
	result := htmlTagPattern.ReplaceAllString(text, "")

	// Normalize whitespace
	result = multiSpacePattern.ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
