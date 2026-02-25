// =============================================================================
// FILE: internal/config/env/like.go
// PURPOSE: Like action environment variable defaults.
//          Ports Python of_env/values/action/like.py.
// =============================================================================

package env

// LikeEnabled returns whether the like action is enabled.
func LikeEnabled() bool {
	return GetBool("OF_LIKE_ENABLED", true)
}
