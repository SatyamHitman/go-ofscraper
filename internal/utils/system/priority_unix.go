// =============================================================================
// FILE: internal/utils/system/priority_unix.go
// PURPOSE: Unix-specific process priority implementation using syscall.Setpriority.
// =============================================================================

//go:build !windows

package system

import "syscall"

// setLowPriority sets the process to nice value 10 on Unix systems.
func setLowPriority() error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, 0, 10)
}
