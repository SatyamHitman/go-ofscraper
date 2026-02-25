// =============================================================================
// FILE: internal/scripts/skip_download.go
// PURPOSE: Skip download script. Runs a user-defined script to determine
//          whether a media item should be skipped during download.
//          Ports Python scripts/skip_download_script.py.
// =============================================================================

package scripts

import (
	"context"
	"fmt"
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Skip download check
// ---------------------------------------------------------------------------

// RunSkipCheck executes the skip-download script to check if a download
// should be skipped.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - total: Total number of already-downloaded items.
//   - media: The media item to check.
//
// Returns:
//   - true if the download should be skipped, or error.
func (m *Manager) RunSkipCheck(ctx context.Context, total int64, media *model.Media) (bool, error) {
	if m.cfg.SkipDownload == "" {
		return false, nil
	}

	env := map[string]string{
		"OF_MEDIA_ID":      fmt.Sprintf("%d", media.ID),
		"OF_POST_ID":       fmt.Sprintf("%d", media.PostID),
		"OF_MEDIA_TYPE":    media.Type,
		"OF_USERNAME":      media.Username,
		"OF_TOTAL":         fmt.Sprintf("%d", total),
		"OF_URL":           media.RawURL,
		"OF_DOWNLOAD_TYPE": string(media.DownloadKind()),
	}

	output, err := m.runScript(ctx, m.cfg.SkipDownload, nil, env)
	if err != nil {
		return false, err
	}

	// Script returns "true"/"skip" to skip, anything else to proceed.
	lower := strings.ToLower(strings.TrimSpace(output))
	return lower == "true" || lower == "skip" || lower == "1", nil
}
