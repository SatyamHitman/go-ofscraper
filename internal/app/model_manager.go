// =============================================================================
// FILE: internal/app/model_manager.go
// PURPOSE: ModelManager processes a single user/model by fetching their content
//          areas, applying filters, and dispatching to download/like/metadata
//          handlers. Ports Python runner/scraper/model_manager.py.
// =============================================================================

package app

import (
	"context"
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// ModelManager
// ---------------------------------------------------------------------------

// ModelManager handles the complete processing pipeline for a single user.
type ModelManager struct {
	app     *App
	user    *model.User
	areas   []string
	actions []string
}

// NewModelManager creates a ModelManager for the given user.
//
// Parameters:
//   - a: The application instance.
//   - user: The user/model to process.
//   - areas: Content areas to fetch.
//   - actions: Actions to perform.
//
// Returns:
//   - A configured ModelManager.
func NewModelManager(a *App, user *model.User, areas, actions []string) *ModelManager {
	return &ModelManager{
		app:     a,
		user:    user,
		areas:   areas,
		actions: actions,
	}
}

// Process runs the complete pipeline for this user: fetch content areas,
// apply filters, and dispatch actions.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - A ModelResult with processing outcomes, and any error.
func (mm *ModelManager) Process(ctx context.Context) (*ModelResult, error) {
	mm.app.logger.Info("processing model",
		"user", mm.user.Name,
		"user_id", mm.user.ID,
		"areas", mm.areas,
		"actions", mm.actions,
	)

	result := &ModelResult{
		Username: mm.user.Name,
		UserID:   mm.user.ID,
	}

	// Mark user as active in the global manager.
	mgr := GetManager()
	if !mgr.MarkUserActive(mm.user.Name) {
		return result, fmt.Errorf("user %s is already being processed", mm.user.Name)
	}
	defer mgr.MarkUserDone(mm.user.Name)

	// Fetch posts from each content area.
	for _, area := range mm.areas {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		posts, err := mm.fetchArea(ctx, area)
		if err != nil {
			mm.app.logger.Error("failed to fetch area",
				"user", mm.user.Name,
				"area", area,
				"error", err,
			)
			result.Errors++
			continue
		}

		result.PostsFound += len(posts)

		// Collect media from posts.
		for _, post := range posts {
			media := post.ViewableMedia()
			result.MediaFound += len(media)
		}
	}

	// Dispatch actions.
	for _, action := range mm.actions {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		if err := mm.app.RunAction(ctx, action, mm.areas, []string{mm.user.Name}); err != nil {
			mm.app.logger.Error("action failed",
				"user", mm.user.Name,
				"action", action,
				"error", err,
			)
			result.Errors++
		}
	}

	mm.app.logger.Info("model processing complete",
		"user", mm.user.Name,
		"posts_found", result.PostsFound,
		"media_found", result.MediaFound,
		"errors", result.Errors,
	)

	return result, nil
}

// fetchArea retrieves posts for a single content area.
func (mm *ModelManager) fetchArea(_ context.Context, area string) ([]*model.Post, error) {
	// TODO: Wire to API post fetcher for the given area.
	_ = area
	return nil, nil
}

// ---------------------------------------------------------------------------
// ModelResult
// ---------------------------------------------------------------------------

// ModelResult holds the outcomes from processing a single model.
type ModelResult struct {
	Username   string
	UserID     int64
	PostsFound int
	MediaFound int
	Downloaded int
	Skipped    int
	Failed     int
	Errors     int
}
