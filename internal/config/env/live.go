// =============================================================================
// FILE: internal/config/env/live.go
// PURPOSE: Live display environment variable defaults for TUI rendering.
//          Ports Python of_env/values/live.py.
// =============================================================================

package env

// LiveDisplayEnabled returns whether the live TUI display is enabled.
func LiveDisplayEnabled() bool {
	return GetBool("OF_LIVE_DISPLAY_ENABLED", true)
}
