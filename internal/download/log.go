// =============================================================================
// FILE: internal/download/log.go
// PURPOSE: Download logging helpers. Provides structured logging for download
//          events (start, progress, completion, failure).
//          Ports Python actions/utils/log.py.
// =============================================================================

package download

import (
	"log/slog"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Download logging
// ---------------------------------------------------------------------------

// LogDownloadStart logs the beginning of a media download.
func LogDownloadStart(logger *slog.Logger, m *model.Media) {
	if logger == nil {
		return
	}
	logger.Info("download starting",
		"media_id", m.ID,
		"type", m.Type,
		"url", m.Link(),
		"protected", m.IsProtected(),
	)
}

// LogDownloadComplete logs successful completion of a download.
func LogDownloadComplete(logger *slog.Logger, m *model.Media, path string) {
	if logger == nil {
		return
	}
	logger.Info("download complete",
		"media_id", m.ID,
		"type", m.Type,
		"path", path,
	)
}

// LogDownloadFailed logs a download failure.
func LogDownloadFailed(logger *slog.Logger, m *model.Media, err error) {
	if logger == nil {
		return
	}
	logger.Error("download failed",
		"media_id", m.ID,
		"type", m.Type,
		"url", m.Link(),
		"error", err,
	)
}

// LogDownloadSkipped logs a skipped download.
func LogDownloadSkipped(logger *slog.Logger, m *model.Media, reason string) {
	if logger == nil {
		return
	}
	logger.Debug("download skipped",
		"media_id", m.ID,
		"reason", reason,
	)
}

// LogBatchSummary logs the summary of a download batch.
func LogBatchSummary(logger *slog.Logger, result *Result, username string) {
	if logger == nil {
		return
	}
	logger.Info("batch summary",
		"username", username,
		"total", result.Total,
		"succeeded", result.Succeeded,
		"failed", result.Failed,
		"skipped", result.Skipped,
	)
}
