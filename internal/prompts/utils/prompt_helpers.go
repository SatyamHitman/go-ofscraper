// =============================================================================
// FILE: internal/prompts/utils/prompt_helpers.go
// PURPOSE: General prompt helper functions. Provides common display formatting
//          and input parsing helpers used across prompt groups.
//          Ports Python prompts/utils/prompt_helpers.py.
// =============================================================================

package promptutils

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Display helpers
// ---------------------------------------------------------------------------

// PrintHeader prints a formatted section header.
//
// Parameters:
//   - title: The header title.
func PrintHeader(title string) {
	line := strings.Repeat("=", len(title)+4)
	fmt.Printf("\n%s\n  %s\n%s\n", line, title, line)
}

// PrintDivider prints a horizontal divider line.
func PrintDivider() {
	fmt.Println(strings.Repeat("-", 40))
}

// PrintNumberedList prints a numbered list of items.
//
// Parameters:
//   - items: Items to display.
//   - startIdx: Starting index number.
func PrintNumberedList(items []string, startIdx int) {
	for i, item := range items {
		fmt.Printf("  %d) %s\n", startIdx+i, item)
	}
}

// ---------------------------------------------------------------------------
// Input helpers
// ---------------------------------------------------------------------------

// ParseIntList parses a comma-separated string of integers.
//
// Parameters:
//   - input: The comma-separated input string.
//
// Returns:
//   - Slice of parsed integers.
func ParseIntList(input string) []int {
	var result []int
	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(part, "%d", &n); err == nil {
			result = append(result, n)
		}
	}
	return result
}

// TruncateDisplay truncates a string for display, adding "..." if needed.
//
// Parameters:
//   - s: The string to truncate.
//   - maxLen: Maximum display length.
//
// Returns:
//   - Truncated string.
func TruncateDisplay(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
