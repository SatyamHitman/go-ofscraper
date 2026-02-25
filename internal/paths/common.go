// =============================================================================
// FILE: internal/paths/common.go
// PURPOSE: Common path operations. Provides shared helpers for path
//          construction, validation, and directory management used across
//          multiple subsystems. Ports Python utils/paths/common.py.
// =============================================================================

package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// Directory operations
// ---------------------------------------------------------------------------

// EnsureDir creates a directory and all parents if they don't exist.
//
// Parameters:
//   - dir: The directory path to create.
//
// Returns:
//   - Error if creation fails.
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// EnsureParentDir creates the parent directory of the given file path.
//
// Parameters:
//   - filePath: The file whose parent directory should exist.
//
// Returns:
//   - Error if creation fails.
func EnsureParentDir(filePath string) error {
	return os.MkdirAll(filepath.Dir(filePath), 0755)
}

// ---------------------------------------------------------------------------
// Path construction
// ---------------------------------------------------------------------------

// JoinSafe joins path components, ensuring the result stays within the base
// directory (prevents path traversal attacks).
//
// Parameters:
//   - base: The base directory.
//   - components: Path components to join.
//
// Returns:
//   - The joined path, and error if the result escapes base.
func JoinSafe(base string, components ...string) (string, error) {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", fmt.Errorf("invalid base path: %w", err)
	}

	joined := filepath.Join(append([]string{absBase}, components...)...)
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("invalid joined path: %w", err)
	}

	// Ensure the resolved path is still within base.
	if !isSubPath(absBase, absJoined) {
		return "", fmt.Errorf("path traversal detected: %q escapes %q", absJoined, absBase)
	}

	return absJoined, nil
}

// isSubPath checks if child is a subdirectory of parent.
func isSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// If the relative path starts with "..", it escapes the parent.
	return rel != ".." && len(rel) >= 1 && rel[0] != '.'
}

// ---------------------------------------------------------------------------
// Path info
// ---------------------------------------------------------------------------

// Exists reports whether a path exists (file or directory).
//
// Parameters:
//   - path: The path to check.
//
// Returns:
//   - true if the path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir reports whether a path is an existing directory.
//
// Parameters:
//   - path: The path to check.
//
// Returns:
//   - true if the path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile reports whether a path is an existing regular file.
//
// Parameters:
//   - path: The path to check.
//
// Returns:
//   - true if the path is a regular file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
