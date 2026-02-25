// =============================================================================
// FILE: internal/config/env/path_files.go
// PURPOSE: Path-related file environment variable defaults.
//          Ports Python of_env/values/path/files.py.
// =============================================================================

package env

// AuthFilePath returns the auth configuration file path override.
func AuthFilePath() string {
	return GetString("OF_AUTH_FILE", "")
}

// ConfigFilePath returns the config file path override.
func ConfigFilePath() string {
	return GetString("OF_CONFIG_FILE", "")
}
