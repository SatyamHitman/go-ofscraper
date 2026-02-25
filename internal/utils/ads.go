// =============================================================================
// FILE: internal/utils/ads.go
// PURPOSE: Ad detection utilities. Identifies promotional/advertisement posts
//          based on text patterns and structural markers. Used by post filters
//          to exclude or flag sponsored content. Ports Python utils/ads.py.
// =============================================================================

package utils

import (
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------
// Ad detection
// ---------------------------------------------------------------------------

// adPatterns are compiled regexes that match common ad/promo indicators.
var adPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bpromotion\b`),
	regexp.MustCompile(`(?i)\bsponsored\b`),
	regexp.MustCompile(`(?i)\bad\b`),
	regexp.MustCompile(`(?i)\b#ad\b`),
	regexp.MustCompile(`(?i)\bcollaboration\b`),
	regexp.MustCompile(`(?i)\bpromo\b`),
	regexp.MustCompile(`(?i)\bpartner(ship)?\b`),
	regexp.MustCompile(`(?i)\baffiliate\b`),
	regexp.MustCompile(`(?i)\bdiscount code\b`),
	regexp.MustCompile(`(?i)\buse (my )?code\b`),
	regexp.MustCompile(`(?i)\blink in bio\b`),
}

// IsAdPost checks whether a post's text contains common advertisement or
// promotional patterns.
//
// Parameters:
//   - text: The post text to analyse.
//
// Returns:
//   - true if the text matches any known ad pattern.
func IsAdPost(text string) bool {
	if text == "" {
		return false
	}
	for _, pat := range adPatterns {
		if pat.MatchString(text) {
			return true
		}
	}
	return false
}

// CountAdIndicators returns the number of distinct ad patterns found in the
// text. A higher count gives stronger confidence the post is promotional.
//
// Parameters:
//   - text: The post text to analyse.
//
// Returns:
//   - Count of matching ad patterns.
func CountAdIndicators(text string) int {
	if text == "" {
		return 0
	}
	count := 0
	for _, pat := range adPatterns {
		if pat.MatchString(text) {
			count++
		}
	}
	return count
}

// ContainsExternalLink checks if text contains URLs pointing outside the
// platform, which is a common indicator of promotional content.
//
// Parameters:
//   - text: The post text to check.
//
// Returns:
//   - true if external links are found.
func ContainsExternalLink(text string) bool {
	// Simple heuristic: check for http(s) URLs that aren't onlyfans.com.
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "http") {
		return false
	}
	// If it has a URL but not an OF URL, it's external.
	hasURL := strings.Contains(lower, "http://") || strings.Contains(lower, "https://")
	isOFOnly := strings.Contains(lower, "onlyfans.com")
	return hasURL && !isOFOnly
}
