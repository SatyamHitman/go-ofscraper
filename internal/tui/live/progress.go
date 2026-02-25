// =============================================================================
// FILE: internal/tui/live/progress.go
// PURPOSE: Progress bar display. Renders download and processing progress
//          bars in the terminal. Ports Python utils/live/progress.py.
// =============================================================================

package live

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Progress bar
// ---------------------------------------------------------------------------

// ProgressBar renders a text-based progress bar.
type ProgressBar struct {
	Width    int
	Total    int64
	Current  int64
	Label    string
}

// Render returns the progress bar as a string.
func (pb *ProgressBar) Render() string {
	if pb.Width <= 0 {
		pb.Width = 40
	}

	pct := float64(0)
	if pb.Total > 0 {
		pct = float64(pb.Current) / float64(pb.Total) * 100
	}

	filled := int(float64(pb.Width) * pct / 100)
	if filled > pb.Width {
		filled = pb.Width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.Width-filled)

	if pb.Label != "" {
		return fmt.Sprintf("%s [%s] %.1f%% (%d/%d)", pb.Label, bar, pct, pb.Current, pb.Total)
	}
	return fmt.Sprintf("[%s] %.1f%% (%d/%d)", bar, pct, pb.Current, pb.Total)
}
