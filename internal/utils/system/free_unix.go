// =============================================================================
// FILE: internal/utils/system/free_unix.go
// PURPOSE: Unix-specific disk usage implementation using syscall.Statfs.
//          Build-tagged for linux and darwin only.
// =============================================================================

//go:build !windows

package system

import (
	"fmt"
	"syscall"
)

// GetDiskUsage returns disk usage statistics for the filesystem containing
// the given path. Unix implementation using statfs.
//
// Parameters:
//   - path: Any path on the target filesystem.
//
// Returns:
//   - DiskUsage with total/free/used bytes, and any error.
func GetDiskUsage(path string) (DiskUsage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskUsage{}, fmt.Errorf("statfs %s: %w", path, err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	return DiskUsage{
		Total: total,
		Free:  free,
		Used:  used,
		Path:  path,
	}, nil
}
