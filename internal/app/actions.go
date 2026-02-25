// =============================================================================
// FILE: internal/app/actions.go
// PURPOSE: Action routing. Dispatches user-selected actions (download, like,
//          unlike, metadata) to the appropriate handlers.
//          Ports Python utils/actions.py.
// =============================================================================

package app

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// Action router
// ---------------------------------------------------------------------------

// RunAction dispatches the given action.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - action: The action name ("download", "like", "unlike", "metadata").
//   - areas: Content areas to process.
//   - usernames: Usernames to process.
//
// Returns:
//   - Error if the action fails.
func (a *App) RunAction(ctx context.Context, action string, areas, usernames []string) error {
	switch action {
	case "download":
		return a.runDownload(ctx, areas, usernames)
	case "like":
		return a.runLike(ctx, usernames, true)
	case "unlike":
		return a.runLike(ctx, usernames, false)
	case "metadata":
		return a.runMetadata(ctx, areas, usernames)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

// runDownload handles the download action.
func (a *App) runDownload(ctx context.Context, areas, usernames []string) error {
	a.logger.Info("starting download",
		"areas", areas,
		"users", usernames,
	)
	// TODO: Wire to download orchestrator.
	return nil
}

// runLike handles the like/unlike action.
func (a *App) runLike(ctx context.Context, usernames []string, like bool) error {
	action := "like"
	if !like {
		action = "unlike"
	}
	a.logger.Info("starting "+action,
		"users", usernames,
	)
	// TODO: Wire to like handler.
	return nil
}

// runMetadata handles the metadata update action.
func (a *App) runMetadata(ctx context.Context, areas, usernames []string) error {
	a.logger.Info("starting metadata update",
		"areas", areas,
		"users", usernames,
	)
	// TODO: Wire to metadata handler.
	return nil
}
