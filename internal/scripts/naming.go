// =============================================================================
// FILE: internal/scripts/naming.go
// PURPOSE: Naming script. Runs a user-defined script to generate custom
//          filenames for downloaded media.
//          Ports Python scripts/naming_script.py.
// =============================================================================

package scripts

import (
	"context"
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Naming script
// ---------------------------------------------------------------------------

// RunNaming executes the naming script to get a custom filename.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - media: The media item to name.
//
// Returns:
//   - The custom filename string, or error.
func (m *Manager) RunNaming(ctx context.Context, media *model.Media) (string, error) {
	if m.cfg.Naming == "" {
		return "", nil
	}

	env := map[string]string{
		"OF_MEDIA_ID":   fmt.Sprintf("%d", media.ID),
		"OF_POST_ID":    fmt.Sprintf("%d", media.PostID),
		"OF_MEDIA_TYPE": media.Type,
		"OF_USERNAME":   media.Username,
		"OF_URL":        media.RawURL,
		"OF_FILENAME":   media.Filename(),
	}

	output, err := m.runScript(ctx, m.cfg.Naming, nil, env)
	if err != nil {
		return "", err
	}

	return output, nil
}
