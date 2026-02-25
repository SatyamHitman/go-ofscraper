// =============================================================================
// FILE: internal/commands/metadata/manager.go
// PURPOSE: MetadataManager manages metadata update operations for a single
//          user. Fetches posts per area and dispatches to consumers.
//          Ports Python metadata/manager.py.
// =============================================================================

package metadata

import (
	"context"
	"log/slog"

	"gofscraper/internal/app"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// MetadataManager
// ---------------------------------------------------------------------------

// MetadataManager handles metadata updates for a single user across
// configured content areas.
type MetadataManager struct {
	logger *slog.Logger
	user   *model.User
	areas  []string
}

// NewManager creates a MetadataManager for the given user.
//
// Parameters:
//   - logger: Structured logger.
//   - user: The user whose metadata to update.
//   - areas: Content areas to process.
//
// Returns:
//   - A configured MetadataManager.
func NewManager(logger *slog.Logger, user *model.User, areas []string) *MetadataManager {
	return &MetadataManager{
		logger: logger,
		user:   user,
		areas:  areas,
	}
}

// Process runs the metadata update for this user across all areas.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance.
//
// Returns:
//   - An UpdateResult with counts, and any error.
func (mm *MetadataManager) Process(ctx context.Context, a *app.App) (UpdateResult, error) {
	var result UpdateResult

	for _, area := range mm.areas {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		mm.logger.Debug("updating metadata",
			"user", mm.user.Name,
			"area", area,
		)

		areaResult, err := mm.processArea(ctx, a, area)
		if err != nil {
			mm.logger.Error("metadata area failed",
				"user", mm.user.Name,
				"area", area,
				"error", err,
			)
			result.Failed++
			continue
		}

		result.Changed += areaResult.Changed
		result.Unchanged += areaResult.Unchanged
		result.Failed += areaResult.Failed
	}

	return result, nil
}

// processArea fetches posts for a content area and updates metadata for each.
func (mm *MetadataManager) processArea(ctx context.Context, a *app.App, area string) (UpdateResult, error) {
	var result UpdateResult
	_ = a.Session() // Will be used when API is wired.

	// TODO: Fetch posts for this area from the API.
	var posts []*model.Post
	_ = area

	// Process each post's media through the consumer.
	consumer := NewConsumer(mm.logger)
	for _, post := range posts {
		for _, media := range post.AllMedia {
			if ctx.Err() != nil {
				return result, ctx.Err()
			}

			status := consumer.ProcessItem(ctx, media)
			switch status {
			case model.MetadataStatusChanged:
				result.Changed++
			case model.MetadataStatusUnchanged:
				result.Unchanged++
			case model.MetadataStatusFailed:
				result.Failed++
			}
		}
	}

	return result, nil
}
