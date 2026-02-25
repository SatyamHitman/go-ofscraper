// =============================================================================
// FILE: internal/config/env/list.go
// PURPOSE: List-related environment variable defaults.
//          Ports Python of_env/values/list.py.
// =============================================================================

package env

// DefaultUserListName returns the default user list name.
func DefaultUserListName() string {
	return GetString("OF_DEFAULT_USER_LIST", "main")
}

// DefaultBlackListName returns the default blacklist name.
func DefaultBlackListName() string {
	return GetString("OF_DEFAULT_BLACK_LIST", "")
}
