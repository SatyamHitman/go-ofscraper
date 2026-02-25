// =============================================================================
// FILE: internal/commands/scraper/print.go
// PURPOSE: Summary printing for scrape runs. Formats and outputs scrape
//          results including users processed, posts found, media downloaded,
//          and errors encountered. Ports Python utils/logs/stdout.py summary.
// =============================================================================

package scraper

import (
	"log/slog"

	cmdutils "gofscraper/internal/commands/utils"
)

// ---------------------------------------------------------------------------
// PrintSummary
// ---------------------------------------------------------------------------

// PrintSummary formats and logs the scrape run results.
//
// Parameters:
//   - sc: The ScrapeContext containing accumulated results.
//   - logger: The structured logger for output.
func PrintSummary(sc *cmdutils.ScrapeContext, logger *slog.Logger) {
	logger.Info(cmdutils.MsgSummaryHeader)

	logger.Info(cmdutils.FormatSummaryLine("Users Processed", sc.UsersProcessed.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Posts Found", sc.PostsFound.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Media Found", sc.MediaFound.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Media Downloaded", sc.MediaDownloaded.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Media Skipped", sc.MediaSkipped.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Media Failed", sc.MediaFailed.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Likes Attempted", sc.LikesAttempted.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Likes Succeeded", sc.LikesSucceeded.Load()))
	logger.Info(cmdutils.FormatSummaryLine("Errors", sc.Errors.Load()))

	logger.Info(cmdutils.MsgSummaryFooter)
}
