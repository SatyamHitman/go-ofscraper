// =============================================================================
// FILE: internal/config/env/path_bytes.go
// PURPOSE: Path-related byte size environment variable defaults.
//          Ports Python of_env/values/path/bytes.py.
// =============================================================================

package env

// PathMaxBytes returns the maximum path length in bytes (0 = OS default).
func PathMaxBytes() int {
	return GetInt("OF_PATH_MAX_BYTES", 0)
}
