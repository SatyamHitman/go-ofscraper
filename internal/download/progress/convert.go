// =============================================================================
// FILE: internal/download/progress/convert.go
// PURPOSE: Progress conversion utilities. Converts download progress values
//          to human-readable format strings.
//          Ports Python progress/convert.py.
// =============================================================================

package progress

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

// ---------------------------------------------------------------------------
// Format helpers
// ---------------------------------------------------------------------------

// FormatBytes formats a byte count as a human-readable string.
//
// Parameters:
//   - bytes: The byte count.
//
// Returns:
//   - Formatted string (e.g. "1.5 MB").
func FormatBytes(bytes int64) string {
	return humanize.Bytes(uint64(bytes))
}

// FormatSpeed formats a bytes-per-second value as human-readable speed.
//
// Parameters:
//   - bytesPerSec: Download speed in bytes per second.
//
// Returns:
//   - Formatted string (e.g. "1.5 MB/s").
func FormatSpeed(bytesPerSec float64) string {
	return humanize.Bytes(uint64(bytesPerSec)) + "/s"
}

// FormatETA calculates and formats the estimated time remaining.
//
// Parameters:
//   - bytesRemaining: Bytes left to download.
//   - bytesPerSec: Current download speed.
//
// Returns:
//   - Formatted ETA string (e.g. "2m30s"), or "unknown" if speed is 0.
func FormatETA(bytesRemaining int64, bytesPerSec float64) string {
	if bytesPerSec <= 0 {
		return "unknown"
	}
	seconds := float64(bytesRemaining) / bytesPerSec
	return FormatDuration(time.Duration(seconds * float64(time.Second)))
}

// FormatDuration formats a duration as a compact string.
//
// Parameters:
//   - d: The duration.
//
// Returns:
//   - Formatted string (e.g. "1h2m30s").
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", m, s)
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", h, m)
}

// FormatProgress formats a progress ratio.
//
// Parameters:
//   - completed: Items completed.
//   - total: Total items.
//
// Returns:
//   - Formatted string (e.g. "42/100 (42.0%)").
func FormatProgress(completed, total int64) string {
	if total == 0 {
		return "0/0 (0.0%)"
	}
	pct := float64(completed) / float64(total) * 100
	return fmt.Sprintf("%d/%d (%.1f%%)", completed, total, pct)
}
