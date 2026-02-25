// =============================================================================
// FILE: internal/paths/paths.go
// PURPOSE: General path resolution. Provides functions for resolving absolute
//          paths, expanding home directories, and building standard application
//          paths (config, data, cache, temp). Ports Python utils/paths/paths.py.
// =============================================================================

package paths

import (
	"os"
	"path/filepath"
	"strings"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Standard application paths
// ---------------------------------------------------------------------------

// ConfigDir returns the gofscraper configuration directory.
//
// Returns:
//   - Absolute path to the config directory.
func ConfigDir() string {
	return env.ConfigDir()
}

// TempDir returns the gofscraper temporary directory for in-progress downloads.
//
// Returns:
//   - Absolute path to the temp directory.
func TempDir() string {
	d := env.TempDir()
	if d != "" {
		return d
	}
	return filepath.Join(os.TempDir(), "gofscraper")
}

// LogDir returns the log file directory inside the config directory.
//
// Returns:
//   - Absolute path to the log directory.
func LogDir() string {
	return filepath.Join(ConfigDir(), "logs")
}

// CacheDir returns the cache directory inside the config directory.
//
// Returns:
//   - Absolute path to the cache directory.
func CacheDir() string {
	return filepath.Join(ConfigDir(), "cache")
}

// ---------------------------------------------------------------------------
// Path resolution
// ---------------------------------------------------------------------------

// ExpandHome expands a leading "~" in a path to the user's home directory.
//
// Parameters:
//   - path: The path to expand.
//
// Returns:
//   - The expanded path, or the original if no ~ prefix.
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}

// Resolve converts a potentially relative path to an absolute path, expanding
// home directory references.
//
// Parameters:
//   - path: The path to resolve.
//
// Returns:
//   - The absolute path, and any error.
func Resolve(path string) (string, error) {
	path = ExpandHome(path)
	return filepath.Abs(path)
}

// MustResolve is like Resolve but panics on error. Use for paths known to be valid.
//
// Parameters:
//   - path: The path to resolve.
//
// Returns:
//   - The absolute path.
func MustResolve(path string) string {
	abs, err := Resolve(path)
	if err != nil {
		panic("failed to resolve path: " + err.Error())
	}
	return abs
}
