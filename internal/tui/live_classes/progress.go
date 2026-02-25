// =============================================================================
// FILE: internal/tui/live_classes/progress.go
// PURPOSE: Download progress data class. Stores bytes transferred, speed, and
//          ETA information for a single download.
//          Ports Python utils/live/classes/progress.py.
// =============================================================================

package liveclasses

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// ProgressInfo
// ---------------------------------------------------------------------------

// ProgressInfo holds download progress data for a single transfer.
type ProgressInfo struct {
	// Identification.
	ID       string
	Filename string

	// Byte counts.
	TotalBytes      int64
	DownloadedBytes int64

	// Timing.
	StartTime time.Time
	UpdateTime time.Time

	// Speed in bytes per second (calculated).
	Speed float64

	// Status.
	Done   bool
	Failed bool
	Error  error
}

// NewProgressInfo creates a new ProgressInfo for the given file.
func NewProgressInfo(id, filename string, totalBytes int64) *ProgressInfo {
	now := time.Now()
	return &ProgressInfo{
		ID:         id,
		Filename:   filename,
		TotalBytes: totalBytes,
		StartTime:  now,
		UpdateTime: now,
	}
}

// Update records new progress data.
func (p *ProgressInfo) Update(downloadedBytes int64) {
	p.DownloadedBytes = downloadedBytes
	p.UpdateTime = time.Now()

	elapsed := p.UpdateTime.Sub(p.StartTime).Seconds()
	if elapsed > 0 {
		p.Speed = float64(p.DownloadedBytes) / elapsed
	}
}

// MarkDone marks the download as completed.
func (p *ProgressInfo) MarkDone() {
	p.Done = true
	p.DownloadedBytes = p.TotalBytes
	p.UpdateTime = time.Now()
}

// MarkFailed marks the download as failed.
func (p *ProgressInfo) MarkFailed(err error) {
	p.Failed = true
	p.Error = err
	p.UpdateTime = time.Now()
}

// Percent returns the download progress as a percentage (0-100).
func (p *ProgressInfo) Percent() float64 {
	if p.TotalBytes <= 0 {
		return 0
	}
	pct := float64(p.DownloadedBytes) / float64(p.TotalBytes) * 100
	if pct > 100 {
		return 100
	}
	return pct
}

// ETA returns the estimated time remaining.
func (p *ProgressInfo) ETA() time.Duration {
	if p.Speed <= 0 || p.Done {
		return 0
	}
	remaining := float64(p.TotalBytes - p.DownloadedBytes)
	if remaining <= 0 {
		return 0
	}
	return time.Duration(remaining/p.Speed) * time.Second
}

// Elapsed returns the time since the download started.
func (p *ProgressInfo) Elapsed() time.Duration {
	return p.UpdateTime.Sub(p.StartTime)
}

// StatusString returns a human-readable status string.
func (p *ProgressInfo) StatusString() string {
	if p.Failed {
		return "Failed"
	}
	if p.Done {
		return "Done"
	}
	return fmt.Sprintf("%.1f%%", p.Percent())
}

// FormatSpeed returns the speed formatted for display.
func (p *ProgressInfo) FormatSpeed() string {
	return FormatBytes(p.Speed) + "/s"
}

// FormatBytes formats a byte count for human-readable display.
func FormatBytes(bytes float64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", bytes/gb)
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", bytes/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", bytes/kb)
	default:
		return fmt.Sprintf("%.0f B", bytes)
	}
}
