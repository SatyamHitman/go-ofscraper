// =============================================================================
// FILE: internal/utils/system/speed.go
// PURPOSE: Speed calculation utilities. Provides functions for measuring and
//          formatting download/upload speeds with human-readable output.
//          Ports Python utils/system/speed.py.
// =============================================================================

package system

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

// ---------------------------------------------------------------------------
// Speed measurement
// ---------------------------------------------------------------------------

// SpeedTracker measures transfer speed over a rolling window.
type SpeedTracker struct {
	startTime time.Time
	lastTime  time.Time
	lastBytes int64
	totalBytes int64
}

// NewSpeedTracker creates a new speed tracker starting now.
//
// Returns:
//   - A new SpeedTracker.
func NewSpeedTracker() *SpeedTracker {
	now := time.Now()
	return &SpeedTracker{
		startTime: now,
		lastTime:  now,
	}
}

// Update records new bytes transferred.
//
// Parameters:
//   - bytes: Number of new bytes since last update.
func (st *SpeedTracker) Update(bytes int64) {
	st.totalBytes += bytes
	st.lastBytes = bytes
	st.lastTime = time.Now()
}

// CurrentSpeed returns the instantaneous speed based on the last update.
//
// Returns:
//   - Speed in bytes per second.
func (st *SpeedTracker) CurrentSpeed() float64 {
	elapsed := time.Since(st.lastTime).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(st.lastBytes) / elapsed
}

// AverageSpeed returns the overall average speed since tracking started.
//
// Returns:
//   - Speed in bytes per second.
func (st *SpeedTracker) AverageSpeed() float64 {
	elapsed := time.Since(st.startTime).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(st.totalBytes) / elapsed
}

// TotalBytes returns the total bytes transferred.
func (st *SpeedTracker) TotalBytes() int64 {
	return st.totalBytes
}

// Elapsed returns time since tracking started.
func (st *SpeedTracker) Elapsed() time.Duration {
	return time.Since(st.startTime)
}

// ---------------------------------------------------------------------------
// Formatting
// ---------------------------------------------------------------------------

// FormatSpeed formats a bytes-per-second value as a human-readable string.
//
// Parameters:
//   - bytesPerSec: Speed in bytes per second.
//
// Returns:
//   - Formatted string like "12.5 MB/s".
func FormatSpeed(bytesPerSec float64) string {
	if bytesPerSec <= 0 {
		return "0 B/s"
	}
	return fmt.Sprintf("%s/s", humanize.Bytes(uint64(bytesPerSec)))
}

// FormatETA estimates remaining time given current speed and remaining bytes.
//
// Parameters:
//   - remaining: Bytes remaining.
//   - bytesPerSec: Current transfer speed.
//
// Returns:
//   - Human-readable ETA string, or "unknown" if speed is zero.
func FormatETA(remaining int64, bytesPerSec float64) string {
	if bytesPerSec <= 0 || remaining <= 0 {
		return "unknown"
	}
	seconds := float64(remaining) / bytesPerSec
	d := time.Duration(seconds * float64(time.Second))

	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
