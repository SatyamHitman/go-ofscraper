// =============================================================================
// FILE: internal/scripts/after_like_action.go
// PURPOSE: After like action script. Runs a user-defined script after
//          like/unlike operations are complete.
//          Ports Python scripts/after_like_action_script.py.
// =============================================================================

package scripts

import (
	"context"
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// After like action
// ---------------------------------------------------------------------------

// RunAfterLikeAction executes the after-like-action script.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - username: The model's username.
//   - posts: The posts that were liked/unliked.
//
// Returns:
//   - Error if the script fails.
func (m *Manager) RunAfterLikeAction(ctx context.Context, username string, posts []model.Post) error {
	if m.cfg.AfterLikeAction == "" {
		return nil
	}

	env := map[string]string{
		"OF_USERNAME":   username,
		"OF_POST_COUNT": fmt.Sprintf("%d", len(posts)),
	}

	_, err := m.runScript(ctx, m.cfg.AfterLikeAction, nil, env)
	return err
}
