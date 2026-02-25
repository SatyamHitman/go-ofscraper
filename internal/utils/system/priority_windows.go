// =============================================================================
// FILE: internal/utils/system/priority_windows.go
// PURPOSE: Windows-specific process priority implementation using
//          SetPriorityClass to BELOW_NORMAL_PRIORITY_CLASS.
// =============================================================================

//go:build windows

package system

import (
	"syscall"
)

// BELOW_NORMAL_PRIORITY_CLASS = 0x00004000
const belowNormalPriorityClass uintptr = 0x00004000

// setLowPriority sets the process to BELOW_NORMAL_PRIORITY_CLASS on Windows.
func setLowPriority() error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setPriorityClass := kernel32.NewProc("SetPriorityClass")

	handle, err := syscall.GetCurrentProcess()
	if err != nil {
		return err
	}

	ret, _, callErr := setPriorityClass.Call(
		uintptr(handle),
		belowNormalPriorityClass,
	)
	if ret == 0 {
		return callErr
	}
	return nil
}
