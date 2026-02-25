// =============================================================================
// FILE: internal/utils/system/system.go
// PURPOSE: System information utilities. Retrieves OS, architecture, and
//          runtime details for diagnostics and logging. Ports Python
//          utils/system/system.py.
// =============================================================================

package system

import (
	"fmt"
	"runtime"
)

// ---------------------------------------------------------------------------
// System info
// ---------------------------------------------------------------------------

// Info holds basic system information.
type Info struct {
	OS       string // Operating system (e.g. "linux", "darwin", "windows")
	Arch     string // Architecture (e.g. "amd64", "arm64")
	NumCPU   int    // Number of logical CPUs
	GoVer    string // Go runtime version
	Compiler string // Go compiler used
}

// GetInfo returns current system information.
//
// Returns:
//   - An Info struct with the current platform details.
func GetInfo() Info {
	return Info{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCPU:   runtime.NumCPU(),
		GoVer:    runtime.Version(),
		Compiler: runtime.Compiler,
	}
}

// Summary returns a one-line summary of system information suitable for
// log output.
//
// Returns:
//   - A formatted summary string.
func Summary() string {
	info := GetInfo()
	return fmt.Sprintf("%s/%s cpus=%d go=%s", info.OS, info.Arch, info.NumCPU, info.GoVer)
}

// IsWindows reports whether the current OS is Windows.
func IsWindows() bool { return runtime.GOOS == "windows" }

// IsDarwin reports whether the current OS is macOS.
func IsDarwin() bool { return runtime.GOOS == "darwin" }

// IsLinux reports whether the current OS is Linux.
func IsLinux() bool { return runtime.GOOS == "linux" }
