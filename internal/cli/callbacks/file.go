// =============================================================================
// FILE: internal/cli/callbacks/file.go
// PURPOSE: File path validation callbacks.
// =============================================================================

package callbacks

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateFileExists checks that the given path exists on disk.
func ValidateFileExists(path string) error {
	if path == "" {
		return nil // empty means not specified, which is valid
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

// ValidateDirExists checks that the given path is an existing directory.
func ValidateDirExists(path string) error {
	if path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", path)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}

// NormalizePath cleans and converts a path to an absolute path.
func NormalizePath(path string) string {
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return path
	}
	return abs
}
