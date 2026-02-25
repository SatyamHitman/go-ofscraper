// =============================================================================
// FILE: internal/auth/headers.go
// PURPOSE: HTTP header construction for OF API requests. Builds the complete
//          set of headers including authentication, signing, and standard
//          browser headers. Ports Python utils/auth/request.py make_headers.
// =============================================================================

package auth

import (
	"gofscraper/internal/config"
)

// ---------------------------------------------------------------------------
// Header construction
// ---------------------------------------------------------------------------

// MakeHeaders builds the complete header map for an OF API request.
// Includes auth cookies, signing headers, and browser identification.
//
// Parameters:
//   - data: The auth credentials.
//
// Returns:
//   - A map of header name to value.
func MakeHeaders(data *Data) map[string]string {
	headers := map[string]string{
		"accept":          "application/json, text/plain, */*",
		"app-token":       config.AppToken,
		"user-id":         data.UserID,
		"user-agent":      data.UserAgent,
		"x-bc":            data.XBC,
		"referer":         "https://onlyfans.com",
		"origin":          "https://onlyfans.com",
		"sec-fetch-dest":  "empty",
		"sec-fetch-mode":  "cors",
		"sec-fetch-site":  "same-site",
		"sec-ch-ua-mobile": "?0",
	}

	// Override app token if auth data provides one.
	if data.AppToken != "" {
		headers["app-token"] = data.AppToken
	}

	return headers
}

// MakeCookies builds the cookie header string from auth data.
//
// Parameters:
//   - data: The auth credentials.
//
// Returns:
//   - The cookie header value.
func MakeCookies(data *Data) string {
	return "sess=" + data.Cookie
}

// AddAuthHeaders adds authentication headers to an existing header map.
// This is used to update headers on an already-constructed request.
//
// Parameters:
//   - headers: The existing header map (modified in place).
//   - data: The auth credentials.
func AddAuthHeaders(headers map[string]string, data *Data) {
	headers["cookie"] = MakeCookies(data)
	headers["user-id"] = data.UserID
	headers["user-agent"] = data.UserAgent
	headers["x-bc"] = data.XBC
	if data.AppToken != "" {
		headers["app-token"] = data.AppToken
	} else {
		headers["app-token"] = config.AppToken
	}
}
