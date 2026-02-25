// =============================================================================
// FILE: internal/cli/callbacks/username_parse.go
// PURPOSE: Username format parsing. Handles URL-style and plain usernames.
// =============================================================================

package callbacks

import (
	"net/url"
	"strings"
)

// ParseUsername extracts a clean username from a raw input that may be a full
// URL (e.g. https://onlyfans.com/username) or a plain name.
func ParseUsername(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	// Try parsing as URL.
	u, err := url.Parse(raw)
	if err == nil && u.Host != "" {
		// Take the first path segment as the username.
		path := strings.TrimPrefix(u.Path, "/")
		if idx := strings.Index(path, "/"); idx > 0 {
			path = path[:idx]
		}
		if path != "" {
			return strings.ToLower(path)
		}
	}

	// Strip leading @ if present.
	raw = strings.TrimPrefix(raw, "@")
	return strings.ToLower(raw)
}

// ParseUsernames parses a slice of raw username inputs.
func ParseUsernames(raw []string) []string {
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		name := ParseUsername(r)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}
