// =============================================================================
// FILE: internal/commands/metadata/metadata.go
// PURPOSE: Metadata update orchestrator. Coordinates the metadata update flow
//          across users and content areas. Ports Python runner/metadata.py.
// =============================================================================

package metadata

import (
	"context"
	"fmt"
	"log/slog"

	"gofscraper/internal/app"
	cmdutils "gofscraper/internal/commands/utils"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// MetadataOrchestrator
// ---------------------------------------------------------------------------

// MetadataOrchestrator coordinates metadata updates across all users.
type MetadataOrchestrator struct {
	logger *slog.Logger
	areas  []string
}

// NewOrchestrator creates a new MetadataOrchestrator.
//
// Parameters:
//   - logger: Structured logger for output.
//   - areas: Content areas to update metadata for.
//
// Returns:
//   - A configured MetadataOrchestrator.
func NewOrchestrator(logger *slog.Logger, areas []string) *MetadataOrchestrator {
	if logger == nil {
		logger = slog.Default()
	}
	return &MetadataOrchestrator{
		logger: logger,
		areas:  areas,
	}
}

// Run executes the metadata update flow for the given users.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance.
//   - users: The users to update metadata for.
//
// Returns:
//   - Error if the orchestration fails.
func (mo *MetadataOrchestrator) Run(ctx context.Context, a *app.App, users []*model.User) error {
	mo.logger.Info("metadata update starting",
		"user_count", len(users),
		"areas", mo.areas,
	)

	if len(users) == 0 {
		mo.logger.Info(cmdutils.MsgNoUsers)
		return nil
	}

	var totalChanged, totalUnchanged, totalFailed int

	for i, user := range users {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		mo.logger.Info(cmdutils.FormatUserProgress(user.Name, i+1, len(users)))

		manager := NewManager(mo.logger, user, mo.areas)
		result, err := manager.Process(ctx, a)
		if err != nil {
			mo.logger.Error("metadata update failed for user",
				"user", user.Name,
				"error", err,
			)
			totalFailed++
			continue
		}

		totalChanged += result.Changed
		totalUnchanged += result.Unchanged
		totalFailed += result.Failed
	}

	mo.logger.Info("metadata update complete",
		"changed", totalChanged,
		"unchanged", totalUnchanged,
		"failed", totalFailed,
	)

	return nil
}

// ---------------------------------------------------------------------------
// UpdateResult aggregates the outcomes of metadata updates.
// ---------------------------------------------------------------------------

// UpdateResult holds the counts from a metadata update operation.
type UpdateResult struct {
	Changed   int
	Unchanged int
	Failed    int
}

// Total returns the total number of items processed.
func (ur UpdateResult) Total() int {
	return ur.Changed + ur.Unchanged + ur.Failed
}

// String returns a formatted summary of the update result.
func (ur UpdateResult) String() string {
	return fmt.Sprintf("changed=%d unchanged=%d failed=%d", ur.Changed, ur.Unchanged, ur.Failed)
}
