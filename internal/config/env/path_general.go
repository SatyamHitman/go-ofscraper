// =============================================================================
// FILE: internal/config/env/path_general.go
// PURPOSE: General path environment variable defaults.
//          Ports Python of_env/values/path/general.py.
// =============================================================================

package env

import (
	"os"
	"path/filepath"
)

// ConfigDir returns the base configuration directory.
// Default: ~/.config/gofscraper (or OF_CONFIG_DIR env override).
func ConfigDir() string {
	override := GetString("OF_CONFIG_DIR", "")
	if override != "" {
		return override
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".config/gofscraper"
	}
	return filepath.Join(homeDir, ".config", "gofscraper")
}

// TempDir returns the temporary files directory override.
// Empty means use system default.
func TempDir() string {
	return GetString("OF_TEMP_DIR", "")
}
