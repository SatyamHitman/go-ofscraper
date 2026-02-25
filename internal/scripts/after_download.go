// =============================================================================
// FILE: internal/scripts/after_download.go
// PURPOSE: After download script. Runs a user-defined script after each
//          individual media file is downloaded.
//          Ports Python scripts/after_download_script.py.
// =============================================================================

package scripts

import (
	"context"
)

// ---------------------------------------------------------------------------
// After download
// ---------------------------------------------------------------------------

// RunAfterDownload executes the after-download script for a single file.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - filePath: Path to the downloaded file.
//
// Returns:
//   - Error if the script fails.
func (m *Manager) RunAfterDownload(ctx context.Context, filePath string) error {
	if m.cfg.AfterDownload == "" {
		return nil
	}

	env := map[string]string{
		"OF_FILE_PATH": filePath,
	}

	_, err := m.runScript(ctx, m.cfg.AfterDownload, nil, env)
	return err
}
