// =============================================================================
// FILE: internal/utils/text.go
// PURPOSE: Text processing utilities. Handles text extraction, cleaning, and
//          formatting for post text content that may contain HTML entities,
//          emoji, and special characters. Ports Python utils/text.py.
// =============================================================================

package utils

import (
	"html"
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------
// Text processing
// ---------------------------------------------------------------------------

// Pre-compiled regexes for text cleaning.
var (
	htmlTagRe    = regexp.MustCompile(`<[^>]+>`)
	multiSpaceRe = regexp.MustCompile(`\s{2,}`)
	newlineRe    = regexp.MustCompile(`\n{3,}`)
)

// CleanText removes HTML tags, decodes HTML entities, normalises whitespace,
// and trims the result.
//
// Parameters:
//   - text: The raw text to clean.
//
// Returns:
//   - The cleaned text.
func CleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove HTML tags.
	text = htmlTagRe.ReplaceAllString(text, " ")

	// Decode HTML entities (&amp; → &, etc.).
	text = html.UnescapeString(text)

	// Collapse multiple spaces.
	text = multiSpaceRe.ReplaceAllString(text, " ")

	// Collapse excessive newlines (3+ → 2).
	text = newlineRe.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// ExtractURLs finds all URLs in the given text.
//
// Parameters:
//   - text: The text to scan.
//
// Returns:
//   - A slice of URL strings found in the text.
func ExtractURLs(text string) []string {
	urlRe := regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	return urlRe.FindAllString(text, -1)
}

// StripEmoji removes common emoji characters from text. This is a best-effort
// removal using Unicode ranges for common emoji blocks.
//
// Parameters:
//   - text: The text to strip.
//
// Returns:
//   - Text with emoji removed.
func StripEmoji(text string) string {
	emojiRe := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2702}-\x{27B0}]|[\x{FE00}-\x{FE0F}]|[\x{1F900}-\x{1F9FF}]`)
	return emojiRe.ReplaceAllString(text, "")
}

// SanitizeForDB cleans text for safe database storage. Removes null bytes
// and normalises unicode.
//
// Parameters:
//   - text: The raw text.
//
// Returns:
//   - Database-safe text.
func SanitizeForDB(text string) string {
	// Remove null bytes which can cause SQLite issues.
	text = strings.ReplaceAll(text, "\x00", "")

	// Clean HTML and normalise.
	return CleanText(text)
}

// SanitizeForFile cleans text for safe use in filenames. More aggressive than
// DB sanitisation — removes special characters, limits length.
//
// Parameters:
//   - text: The raw text.
//   - maxLen: Maximum length of the result.
//
// Returns:
//   - Filename-safe text.
func SanitizeForFile(text string, maxLen int) string {
	// Start with basic cleaning.
	text = CleanText(text)

	// Remove characters unsafe for filenames.
	unsafeRe := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	text = unsafeRe.ReplaceAllString(text, "")

	// Collapse spaces.
	text = multiSpaceRe.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Truncate.
	if maxLen > 0 {
		text = Truncate(text, maxLen)
	}

	return text
}
