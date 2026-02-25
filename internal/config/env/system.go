// =============================================================================
// FILE: internal/config/env/system.go
// PURPOSE: System-level environment variable defaults (priority, free space).
//          Ports Python of_env/values/system.py.
// =============================================================================

package env

// SystemFreeMin returns the minimum free disk space in bytes (0 = no check).
func SystemFreeMin() int64 {
	return GetInt64("OF_SYSTEM_FREE_MIN", 0)
}
