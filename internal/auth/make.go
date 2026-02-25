// =============================================================================
// FILE: internal/auth/make.go
// PURPOSE: Auth creation helpers. Provides functions for creating new auth
//          files and interactively collecting auth credentials from the user.
//          Ports Python utils/auth/make.py.
// =============================================================================

package auth

// ---------------------------------------------------------------------------
// Auth creation
// ---------------------------------------------------------------------------

// NewEmptyData creates a Data struct with empty fields, suitable for creating
// a new auth file template.
//
// Returns:
//   - A Data struct with all fields empty.
func NewEmptyData() *Data {
	return &Data{}
}

// NewData creates a Data struct from individual credential strings.
//
// Parameters:
//   - cookie: The sess cookie value.
//   - userAgent: The browser user-agent string.
//   - xbc: The x-bc header token.
//   - userID: The numeric user ID.
//
// Returns:
//   - A populated Data struct.
func NewData(cookie, userAgent, xbc, userID string) *Data {
	return &Data{
		Cookie:    cookie,
		UserAgent: userAgent,
		XBC:       xbc,
		UserID:    userID,
	}
}

// MissingFields returns the names of required fields that are empty.
//
// Parameters:
//   - data: The auth data to check.
//
// Returns:
//   - Slice of missing field labels.
func MissingFields(data *Data) []string {
	if data == nil {
		return []string{"all fields"}
	}

	var missing []string
	if data.Cookie == "" {
		missing = append(missing, "Session Cookie")
	}
	if data.UserAgent == "" {
		missing = append(missing, "User Agent")
	}
	if data.XBC == "" {
		missing = append(missing, "X-BC Token")
	}
	if data.UserID == "" {
		missing = append(missing, "User ID")
	}
	return missing
}
