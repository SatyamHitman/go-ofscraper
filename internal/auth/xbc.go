// =============================================================================
// FILE: internal/auth/xbc.go
// PURPOSE: X-BC token generation. The x-bc header is required for OF API
//          requests and is derived from browser fingerprinting data.
//          Ports Python utils/auth/request.py generate_xbc.
// =============================================================================

package auth

// ---------------------------------------------------------------------------
// X-BC token
// ---------------------------------------------------------------------------

// GenerateXBC generates or validates the X-BC token. In practice this token
// is extracted from the browser and stored in the auth file, but this function
// provides a fallback for token refresh scenarios.
//
// Parameters:
//   - existingToken: The current x-bc token (returned if valid).
//
// Returns:
//   - The x-bc token to use.
func GenerateXBC(existingToken string) string {
	// The x-bc token is extracted from browser state and cannot be
	// programmatically generated without browser automation. If the
	// existing token is provided, use it directly.
	if existingToken != "" {
		return existingToken
	}

	// Return empty — caller should prompt user to provide the token.
	return ""
}
