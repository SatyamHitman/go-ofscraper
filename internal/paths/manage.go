// =============================================================================
// FILE: internal/paths/manage.go
// PURPOSE: Path management utilities. Provides functions for managing
//          temporary files, cleanup of incomplete downloads, and directory
//          listing operations. Ports Python utils/paths/manage.py.
// =============================================================================

package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Temp file management
// ---------------------------------------------------------------------------

// TempFilePath returns a temporary file path for an in-progress download.
// The temp file is placed in the gofscraper temp directory with a .part suffix.
//
// Parameters:
//   - finalPath: The intended final file path.
//
// Returns:
//   - The temporary file path.
func TempFilePath(finalPath string) string {
	name := filepath.Base(finalPath) + ".part"
	return filepath.Join(TempDir(), name)
}

// CleanupTempFiles removes .part files from the temp directory that are older
// than the current session. Called during startup.
//
// Returns:
//   - Number of files removed, and any error.
func CleanupTempFiles() (int, error) {
	tempDir := TempDir()
	if !Exists(tempDir) {
		return 0, nil
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read temp dir: %w", err)
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".part") {
			path := filepath.Join(tempDir, entry.Name())
			if err := os.Remove(path); err == nil {
				removed++
			}
		}
	}

	return removed, nil
}

// PromoteTemp moves a temporary file to its final path. Creates parent
// directories if needed.
//
// Parameters:
//   - tempPath: Path to the temp file.
//   - finalPath: Desired final location.
//
// Returns:
//   - Error if the move fails.
func PromoteTemp(tempPath, finalPath string) error {
	if err := EnsureParentDir(finalPath); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", finalPath, err)
	}

	if err := os.Rename(tempPath, finalPath); err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", tempPath, finalPath, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Directory listing
// ---------------------------------------------------------------------------

// ListFiles returns all regular files in a directory (non-recursive).
//
// Parameters:
//   - dir: The directory to list.
//
// Returns:
//   - Slice of absolute file paths, and any error.
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory %s: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}

// ListDirs returns all subdirectories in a directory (non-recursive).
//
// Parameters:
//   - dir: The directory to list.
//
// Returns:
//   - Slice of absolute directory paths, and any error.
func ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory %s: %w", dir, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(dir, entry.Name()))
		}
	}
	return dirs, nil
}
