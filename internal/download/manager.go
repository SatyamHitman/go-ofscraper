// =============================================================================
// FILE: internal/download/manager.go
// PURPOSE: Download manager helpers. Provides batch download utilities and
//          convenience methods for the Orchestrator.
//          Ports Python managers/downloadmanager.py.
// =============================================================================

package download

import (
	"context"

	"gofscraper/internal/drm"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// DRM integration
// ---------------------------------------------------------------------------

// SetDRMManager configures the DRM manager for protected downloads.
//
// Parameters:
//   - dm: The DRM manager to use.
func (o *Orchestrator) SetDRMManager(dm *drm.Manager) {
	o.drm = dm
}

// ---------------------------------------------------------------------------
// Batch utilities
// ---------------------------------------------------------------------------

// DownloadBatch is a convenience wrapper that downloads media for a single model.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - username: The model's username (for logging).
//   - media: Media items to download.
//
// Returns:
//   - Result and error.
func (o *Orchestrator) DownloadBatch(ctx context.Context, username string, media []*model.Media) (*Result, error) {
	if o.logger != nil {
		o.logger.Info("starting batch download",
			"username", username,
			"count", len(media),
		)
	}

	result, err := o.Run(ctx, media)
	if err != nil {
		return result, err
	}

	if o.logger != nil {
		o.logger.Info("batch download complete",
			"username", username,
			"summary", result.Summary(),
		)
	}

	return result, nil
}

// FilterDownloadable returns media items that are eligible for download.
// Removes items without URLs and optionally skips already-downloaded items.
//
// Parameters:
//   - media: Media items to filter.
//   - downloadedIDs: Set of already-downloaded media IDs (nil = skip none).
//
// Returns:
//   - Filtered media slice.
func FilterDownloadable(media []*model.Media, downloadedIDs map[int64]bool) []*model.Media {
	var result []*model.Media
	for _, m := range media {
		if !m.IsLinked() {
			continue
		}
		if downloadedIDs != nil && downloadedIDs[m.ID] {
			continue
		}
		result = append(result, m)
	}
	return result
}
