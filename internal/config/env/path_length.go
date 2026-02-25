// =============================================================================
// FILE: internal/config/env/path_length.go
// PURPOSE: Path length limit environment variable defaults.
//          Ports Python of_env/values/path/length.py.
// =============================================================================

package env

// MaxPathLength returns the maximum total file path length.
// 0 = use OS default (260 on Windows, 4096 on Linux/macOS).
func MaxPathLength() int {
	return GetInt("OF_MAX_PATH_LENGTH", 0)
}

// MaxFilenameLength returns the maximum filename length.
// 0 = use OS default (255 on most systems).
func MaxFilenameLength() int {
	return GetInt("OF_MAX_FILENAME_LENGTH", 0)
}
