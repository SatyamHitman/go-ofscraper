// =============================================================================
// FILE: internal/scripts/after_download_action.go
// PURPOSE: After download action script. Runs a user-defined script after
//          all downloads for a model are complete.
//          Ports Python scripts/after_download_action_script.py.
// =============================================================================

package scripts

import (
	"context"
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// After download action
// ---------------------------------------------------------------------------

// RunAfterDownloadAction executes the after-download-action script.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - username: The model's username.
//   - media: The media items that were downloaded.
//   - action: The action that was performed (e.g. "download").
//
// Returns:
//   - Error if the script fails.
func (m *Manager) RunAfterDownloadAction(ctx context.Context, username string, media []*model.Media, action string) error {
	if m.cfg.AfterDownloadAction == "" {
		return nil
	}

	env := map[string]string{
		"OF_USERNAME":    username,
		"OF_ACTION":      action,
		"OF_MEDIA_COUNT": fmt.Sprintf("%d", len(media)),
	}

	_, err := m.runScript(ctx, m.cfg.AfterDownloadAction, nil, env)
	return err
}
