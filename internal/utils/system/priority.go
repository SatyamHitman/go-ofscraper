// =============================================================================
// FILE: internal/utils/system/priority.go
// PURPOSE: Process priority management. Provides functions to lower the
//          process priority so gofscraper runs as a background-friendly task
//          without starving other applications. Ports Python
//          utils/system/priority.py.
// =============================================================================

package system

import (
	"log/slog"
	"os"
	"runtime"
)

// ---------------------------------------------------------------------------
// Process priority
// ---------------------------------------------------------------------------

// SetLowPriority attempts to lower the current process priority. On Unix
// this uses nice(10), on Windows it sets BELOW_NORMAL_PRIORITY_CLASS.
// Errors are logged but not returned since this is a best-effort operation.
func SetLowPriority() {
	if err := setLowPriority(); err != nil {
		slog.Debug("failed to set low priority", "error", err, "os", runtime.GOOS)
	} else {
		slog.Debug("process priority set to low", "pid", os.Getpid())
	}
}

// SetMaxProcs configures GOMAXPROCS to use at most n CPUs. If n <= 0 or
// greater than available CPUs, uses all available CPUs.
//
// Parameters:
//   - n: Maximum number of CPUs to use.
//
// Returns:
//   - The previous GOMAXPROCS setting.
func SetMaxProcs(n int) int {
	if n <= 0 || n > runtime.NumCPU() {
		n = runtime.NumCPU()
	}
	return runtime.GOMAXPROCS(n)
}
