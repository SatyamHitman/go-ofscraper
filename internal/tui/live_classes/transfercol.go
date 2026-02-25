// =============================================================================
// FILE: internal/tui/live_classes/transfercol.go
// PURPOSE: Transfer rate column formatter. Formats bytes-per-second transfer
//          rates into human-readable columns for the live display.
//          Ports Python utils/live/classes/transfercol.py.
// =============================================================================

package liveclasses

import (
	"fmt"
)

// ---------------------------------------------------------------------------
// TransferColumn
// ---------------------------------------------------------------------------

// TransferColumn formats a transfer rate for display as a table column.
type TransferColumn struct {
	BytesPerSecond float64
	Downloaded     int64
	Total          int64
}

// NewTransferColumn creates a new TransferColumn.
func NewTransferColumn(bytesPerSec float64, downloaded, total int64) TransferColumn {
	return TransferColumn{
		BytesPerSecond: bytesPerSec,
		Downloaded:     downloaded,
		Total:          total,
	}
}

// SpeedString returns the formatted transfer speed.
func (tc TransferColumn) SpeedString() string {
	return formatRate(tc.BytesPerSecond)
}

// DownloadedString returns the formatted downloaded amount.
func (tc TransferColumn) DownloadedString() string {
	return FormatBytes(float64(tc.Downloaded))
}

// TotalString returns the formatted total size.
func (tc TransferColumn) TotalString() string {
	if tc.Total <= 0 {
		return "?"
	}
	return FormatBytes(float64(tc.Total))
}

// TransferString returns "downloaded / total" formatted string.
func (tc TransferColumn) TransferString() string {
	return fmt.Sprintf("%s / %s", tc.DownloadedString(), tc.TotalString())
}

// View returns the complete transfer column display string.
func (tc TransferColumn) View() string {
	return fmt.Sprintf("%s  %s", tc.TransferString(), tc.SpeedString())
}

// formatRate formats bytes per second into a human-readable rate.
func formatRate(bytesPerSec float64) string {
	const (
		kb = 1024.0
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytesPerSec >= gb:
		return fmt.Sprintf("%.2f GB/s", bytesPerSec/gb)
	case bytesPerSec >= mb:
		return fmt.Sprintf("%.2f MB/s", bytesPerSec/mb)
	case bytesPerSec >= kb:
		return fmt.Sprintf("%.2f KB/s", bytesPerSec/kb)
	case bytesPerSec > 0:
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	default:
		return "0 B/s"
	}
}
