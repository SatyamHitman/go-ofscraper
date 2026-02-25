// =============================================================================
// FILE: internal/utils/system/free.go
// PURPOSE: Disk free space checking. Provides cross-platform functions to
//          query available disk space for download destination validation.
//          Ports Python utils/system/free.py.
// =============================================================================

package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

// ---------------------------------------------------------------------------
// Disk space
// ---------------------------------------------------------------------------

// DiskUsage holds disk space information for a filesystem.
type DiskUsage struct {
	Total     uint64 // Total capacity in bytes
	Free      uint64 // Available space in bytes
	Used      uint64 // Used space in bytes
	Path      string // Mount point or drive path
}

// FreeSpaceString returns a human-readable string of available disk space.
//
// Parameters:
//   - du: The disk usage data.
//
// Returns:
//   - Formatted string like "12.5 GB free of 500 GB".
func (du DiskUsage) FreeSpaceString() string {
	return fmt.Sprintf("%s free of %s", humanize.Bytes(du.Free), humanize.Bytes(du.Total))
}

// HasMinFreeSpace checks if the disk has at least minBytes free.
//
// Parameters:
//   - path: Path on the filesystem to check.
//   - minBytes: Minimum required free bytes.
//
// Returns:
//   - true if there is sufficient space, false otherwise, and any error.
func HasMinFreeSpace(path string, minBytes uint64) (bool, error) {
	du, err := GetDiskUsage(path)
	if err != nil {
		return false, err
	}
	return du.Free >= minBytes, nil
}

// EnsureParentDir ensures the parent directory of the given path exists,
// creating it if necessary.
//
// Parameters:
//   - path: The file path whose parent should exist.
//
// Returns:
//   - Error if directory creation fails.
func EnsureParentDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
