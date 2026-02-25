// =============================================================================
// FILE: internal/config/env/date.go
// PURPOSE: Date formatting environment variable defaults.
//          Ports Python of_env/values/date.py.
// =============================================================================

package env

// DefaultDateFormat returns the default date display format string.
func DefaultDateFormat() string {
	return GetString("OF_DATE_FORMAT", "MM-DD-YYYY")
}
