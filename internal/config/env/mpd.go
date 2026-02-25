// =============================================================================
// FILE: internal/config/env/mpd.go
// PURPOSE: MPD (MPEG-DASH) related environment variable defaults.
//          Ports Python of_env/values/req/mpd.py.
// =============================================================================

package env

// MPDEnabled returns whether MPD/DASH downloading is enabled.
func MPDEnabled() bool {
	return GetBool("OF_MPD_ENABLED", true)
}
