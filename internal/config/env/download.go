// =============================================================================
// FILE: internal/config/env/download.go
// PURPOSE: Download action environment variable defaults.
//          Ports Python of_env/values/action/download.py.
// =============================================================================

package env

// DownloadEnabled returns whether downloading is enabled.
func DownloadEnabled() bool {
	return GetBool("OF_DOWNLOAD_ENABLED", true)
}
