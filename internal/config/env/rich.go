// =============================================================================
// FILE: internal/config/env/rich.go
// PURPOSE: Rich display environment variable defaults for formatted output.
//          Ports Python of_env/values/rich.py.
// =============================================================================

package env

// RichOutputEnabled returns whether rich/formatted terminal output is enabled.
func RichOutputEnabled() bool {
	return GetBool("OF_RICH_OUTPUT", true)
}
