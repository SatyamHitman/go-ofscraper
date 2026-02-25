// =============================================================================
// FILE: internal/commands/metadata/consumer.go
// PURPOSE: MetadataConsumer processes individual metadata update items.
//          Compares current metadata with stored values and updates the
//          database when changes are detected. Ports Python
//          metadata/consumer.py.
// =============================================================================

package metadata

import (
	"context"
	"log/slog"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// MetadataConsumer
// ---------------------------------------------------------------------------

// MetadataConsumer processes individual media items for metadata updates.
type MetadataConsumer struct {
	logger *slog.Logger
}

// NewConsumer creates a MetadataConsumer.
//
// Parameters:
//   - logger: Structured logger.
//
// Returns:
//   - A configured MetadataConsumer.
func NewConsumer(logger *slog.Logger) *MetadataConsumer {
	if logger == nil {
		logger = slog.Default()
	}
	return &MetadataConsumer{
		logger: logger,
	}
}

// ProcessItem checks and updates metadata for a single media item.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - media: The media item to process.
//
// Returns:
//   - The MetadataStatus indicating the outcome.
func (mc *MetadataConsumer) ProcessItem(ctx context.Context, media *model.Media) model.MetadataStatus {
	if ctx.Err() != nil {
		return model.MetadataStatusSkipped
	}

	if media == nil {
		return model.MetadataStatusSkipped
	}

	// Skip media without links.
	if !media.IsLinked() {
		mc.logger.Debug("skipping unlinked media",
			"media_id", media.ID,
			"post_id", media.PostID,
		)
		media.MarkMetadataUnchanged()
		return model.MetadataStatusSkipped
	}

	// Extract description metadata.
	desc := ExtractDesc(media)

	// TODO: Compare extracted metadata with stored values in the database.
	// If changed, update the database and mark as changed.
	// For now, treat all as unchanged.
	_ = desc

	mc.logger.Debug("metadata checked",
		"media_id", media.ID,
		"type", media.Type,
		"linked", media.IsLinked(),
	)

	media.MarkMetadataUnchanged()
	return model.MetadataStatusUnchanged
}

// ProcessBatch processes a slice of media items and returns aggregate results.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - items: The media items to process.
//
// Returns:
//   - An UpdateResult with counts.
func (mc *MetadataConsumer) ProcessBatch(ctx context.Context, items []*model.Media) UpdateResult {
	var result UpdateResult

	for _, item := range items {
		if ctx.Err() != nil {
			break
		}

		status := mc.ProcessItem(ctx, item)
		switch status {
		case model.MetadataStatusChanged:
			result.Changed++
		case model.MetadataStatusUnchanged:
			result.Unchanged++
		case model.MetadataStatusFailed:
			result.Failed++
		}
	}

	return result
}
