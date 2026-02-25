// =============================================================================
// FILE: internal/utils/me.go
// PURPOSE: Current user ("me") utility functions. Stores and retrieves the
//          authenticated user's profile information for use across subsystems.
//          Ports Python utils/me.py.
// =============================================================================

package utils

import (
	"sync"
)

// ---------------------------------------------------------------------------
// Me data (authenticated user)
// ---------------------------------------------------------------------------

// MeData holds the authenticated user's basic profile information retrieved
// from the /me API endpoint.
type MeData struct {
	ID       int64  // Numeric user ID
	Name     string // Display name
	Username string // URL-safe username
	IsPerformer bool // Whether the user is a creator
}

var (
	// meData is the global authenticated user info.
	meData *MeData

	// meMu protects meData.
	meMu sync.RWMutex
)

// SetMe stores the authenticated user's profile data.
//
// Parameters:
//   - data: The user's profile data from /me.
func SetMe(data *MeData) {
	meMu.Lock()
	defer meMu.Unlock()
	meData = data
}

// GetMe retrieves the authenticated user's profile data.
//
// Returns:
//   - The MeData, or nil if not yet set.
func GetMe() *MeData {
	meMu.RLock()
	defer meMu.RUnlock()
	return meData
}

// GetMyID returns the authenticated user's numeric ID, or 0 if not set.
//
// Returns:
//   - The user's ID.
func GetMyID() int64 {
	meMu.RLock()
	defer meMu.RUnlock()
	if meData == nil {
		return 0
	}
	return meData.ID
}

// GetMyUsername returns the authenticated user's username, or empty if not set.
//
// Returns:
//   - The username string.
func GetMyUsername() string {
	meMu.RLock()
	defer meMu.RUnlock()
	if meData == nil {
		return ""
	}
	return meData.Username
}
