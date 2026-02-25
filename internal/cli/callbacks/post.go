// =============================================================================
// FILE: internal/cli/callbacks/post.go
// PURPOSE: Post argument validation callbacks.
// =============================================================================

package callbacks

import (
	"fmt"
	"strings"
)

// validPostAreas is the set of allowed content area values.
var validPostAreas = map[string]bool{
	"timeline":   true,
	"messages":   true,
	"archived":   true,
	"stories":    true,
	"highlights": true,
	"purchased":  true,
	"labels":     true,
	"pinned":     true,
}

// ValidatePostAreas checks that every provided area string is a recognised
// content area. Returns an error listing any invalid values.
func ValidatePostAreas(areas []string) error {
	var bad []string
	for _, a := range areas {
		if !validPostAreas[strings.ToLower(strings.TrimSpace(a))] {
			bad = append(bad, a)
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("invalid post area(s): %s", strings.Join(bad, ", "))
	}
	return nil
}

// NormalizePostAreas lowercases and trims each area string.
func NormalizePostAreas(areas []string) []string {
	out := make([]string, len(areas))
	for i, a := range areas {
		out[i] = strings.ToLower(strings.TrimSpace(a))
	}
	return out
}
