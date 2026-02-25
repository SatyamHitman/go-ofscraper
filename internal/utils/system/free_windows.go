// =============================================================================
// FILE: internal/utils/system/free_windows.go
// PURPOSE: Windows-specific disk usage implementation using
//          GetDiskFreeSpaceExW. Build-tagged for windows only.
// =============================================================================

//go:build windows

package system

import (
	"fmt"
	"syscall"
	"unsafe"
)

// GetDiskUsage returns disk usage statistics for the filesystem containing
// the given path. Windows implementation using GetDiskFreeSpaceExW.
//
// Parameters:
//   - path: Any path on the target filesystem.
//
// Returns:
//   - DiskUsage with total/free/used bytes, and any error.
func GetDiskUsage(path string) (DiskUsage, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable, totalBytes, totalFreeBytes uint64

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return DiskUsage{}, fmt.Errorf("invalid path %s: %w", path, err)
	}

	ret, _, callErr := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		return DiskUsage{}, fmt.Errorf("GetDiskFreeSpaceExW failed for %s: %w", path, callErr)
	}

	return DiskUsage{
		Total: totalBytes,
		Free:  freeBytesAvailable,
		Used:  totalBytes - freeBytesAvailable,
		Path:  path,
	}, nil
}
