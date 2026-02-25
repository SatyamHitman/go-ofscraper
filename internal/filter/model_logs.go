// =============================================================================
// FILE: internal/filter/model_logs.go
// PURPOSE: Model filter logging. Provides structured logging for filter
//          operations to help users understand what was filtered and why.
//          Ports Python filters/models/utils/logs.py.
// =============================================================================

package filter

import (
	"log/slog"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Filter logging
// ---------------------------------------------------------------------------

// LogFilterResult logs the before/after counts of a filter operation.
//
// Parameters:
//   - logger: The slog logger to write to.
//   - filterName: Human-readable name of the filter.
//   - before: Count before filtering.
//   - after: Count after filtering.
func LogFilterResult(logger *slog.Logger, filterName string, before, after int) {
	if logger == nil {
		return
	}
	removed := before - after
	if removed == 0 {
		logger.Debug("filter passed all",
			"filter", filterName,
			"count", after,
		)
		return
	}
	logger.Info("filter applied",
		"filter", filterName,
		"before", before,
		"after", after,
		"removed", removed,
	)
}

// LogMediaFilterChain logs media filter chain results.
//
// Parameters:
//   - logger: The slog logger.
//   - media: The media slice after filtering.
//   - original: The original count before any filtering.
func LogMediaFilterChain(logger *slog.Logger, media []*model.Media, original int) {
	if logger == nil {
		return
	}
	logger.Info("media filter chain complete",
		"original", original,
		"remaining", len(media),
		"removed", original-len(media),
	)
}

// LogPostFilterChain logs post filter chain results.
//
// Parameters:
//   - logger: The slog logger.
//   - posts: The posts slice after filtering.
//   - original: The original count before any filtering.
func LogPostFilterChain(logger *slog.Logger, posts []model.Post, original int) {
	if logger == nil {
		return
	}
	logger.Info("post filter chain complete",
		"original", original,
		"remaining", len(posts),
		"removed", original-len(posts),
	)
}

// LogModelFilterChain logs model filter chain results.
//
// Parameters:
//   - logger: The slog logger.
//   - users: The users slice after filtering.
//   - original: The original count before any filtering.
func LogModelFilterChain(logger *slog.Logger, users []model.User, original int) {
	if logger == nil {
		return
	}
	logger.Info("model filter chain complete",
		"original", original,
		"remaining", len(users),
		"removed", original-len(users),
	)
}
