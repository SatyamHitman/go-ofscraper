// =============================================================================
// FILE: internal/auth/auth.go
// PURPOSE: Auth data loading and management. Loads authentication credentials
//          (cookies, user-agent, x-bc token) from the auth JSON file and
//          provides them to the HTTP session for API requests. Ports Python
//          utils/auth/data.py.
// =============================================================================

package auth

import (
	"fmt"
	"sync"
)

// ---------------------------------------------------------------------------
// Auth data
// ---------------------------------------------------------------------------

// Data holds the authentication credentials needed for OF API access.
type Data struct {
	Cookie    string `json:"auth_cookie"`     // sess cookie value
	UserAgent string `json:"auth_user_agent"` // Browser user-agent
	XBC       string `json:"auth_x_bc"`       // x-bc header token
	UserID    string `json:"auth_user_id"`    // Numeric user ID string
	AppToken  string `json:"auth_app_token"`  // App token (usually static)
}

var (
	// current holds the loaded auth data.
	current *Data
	authMu  sync.RWMutex
)

// Load reads auth credentials from the auth file and stores them globally.
//
// Parameters:
//   - authPath: Path to the auth JSON file.
//
// Returns:
//   - The loaded Data, and any error.
func Load(authPath string) (*Data, error) {
	data, err := ReadAuthFile(authPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load auth: %w", err)
	}

	if err := Validate(data); err != nil {
		return nil, fmt.Errorf("auth validation failed: %w", err)
	}

	authMu.Lock()
	current = data
	authMu.Unlock()

	return data, nil
}

// Get returns the currently loaded auth data.
//
// Returns:
//   - The current *Data, or nil if not loaded.
func Get() *Data {
	authMu.RLock()
	defer authMu.RUnlock()
	return current
}

// Set replaces the current auth data.
//
// Parameters:
//   - data: The new auth credentials.
func Set(data *Data) {
	authMu.Lock()
	defer authMu.Unlock()
	current = data
}

// IsLoaded reports whether auth data has been loaded.
//
// Returns:
//   - true if credentials are available.
func IsLoaded() bool {
	authMu.RLock()
	defer authMu.RUnlock()
	return current != nil
}

// Validate checks that all required auth fields are present.
//
// Parameters:
//   - data: The auth data to validate.
//
// Returns:
//   - Error describing the first missing field, or nil.
func Validate(data *Data) error {
	if data == nil {
		return fmt.Errorf("auth data is nil")
	}
	if data.Cookie == "" {
		return fmt.Errorf("auth_cookie is required")
	}
	if data.UserAgent == "" {
		return fmt.Errorf("auth_user_agent is required")
	}
	if data.XBC == "" {
		return fmt.Errorf("auth_x_bc is required")
	}
	if data.UserID == "" {
		return fmt.Errorf("auth_user_id is required")
	}
	return nil
}
