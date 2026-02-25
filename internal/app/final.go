// =============================================================================
// FILE: internal/app/final.go
// PURPOSE: Final cleanup operations run after a scrape completes. Logs summary
//          statistics, executes user-defined after-action scripts, and sends
//          Discord notifications. Ports Python runner/close/final.py.
// =============================================================================

package app

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"gofscraper/internal/config"
)

// ---------------------------------------------------------------------------
// FinalCleanup
// ---------------------------------------------------------------------------

// FinalCleanup runs post-scrape cleanup operations.
//
// Parameters:
//   - logger: Structured logger for output.
//   - stats: The scrape run statistics to report.
func FinalCleanup(logger *slog.Logger, stats *Stats) {
	if stats == nil {
		return
	}

	stats.MarkDone()

	// Log summary statistics.
	logger.Info("final summary")
	stats.Log(logger)

	// Run after-action script if configured.
	runAfterActionScript(logger)

	// Send Discord notification if configured.
	sendDiscordNotification(logger, stats)
}

// ---------------------------------------------------------------------------
// After-action script
// ---------------------------------------------------------------------------

// runAfterActionScript executes the configured after-action script if set.
func runAfterActionScript(logger *slog.Logger) {
	script := config.GetAfterActionScript()
	if script == "" {
		return
	}

	logger.Info("running after-action script", "script", script)

	cmd := exec.Command("sh", "-c", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("after-action script failed",
			"script", script,
			"error", err,
			"output", string(output),
		)
		return
	}

	if len(output) > 0 {
		logger.Debug("after-action script output", "output", string(output))
	}
}

// ---------------------------------------------------------------------------
// Discord notification
// ---------------------------------------------------------------------------

// sendDiscordNotification sends a summary to the configured Discord webhook.
func sendDiscordNotification(logger *slog.Logger, stats *Stats) {
	webhook := config.GetDiscord()
	if webhook == "" {
		return
	}

	logger.Debug("sending Discord notification")

	message := formatDiscordMessage(stats)
	_ = message
	_ = webhook

	// TODO: Wire to HTTP POST to Discord webhook with the message payload.
	logger.Debug("Discord notification placeholder", "webhook_configured", true)
}

// formatDiscordMessage builds a Discord-formatted summary message from stats.
func formatDiscordMessage(stats *Stats) string {
	var sb strings.Builder
	sb.WriteString("**gofscraper Run Complete**\n")
	sb.WriteString(fmt.Sprintf("Duration: `%s`\n", stats.Duration()))
	sb.WriteString(fmt.Sprintf("Users: `%d`\n", stats.UsersProcessed.Load()))
	sb.WriteString(fmt.Sprintf("Posts: `%d`\n", stats.PostsFound.Load()))
	sb.WriteString(fmt.Sprintf("Media: `%d` downloaded, `%d` skipped, `%d` failed\n",
		stats.MediaDownloaded.Load(),
		stats.MediaSkipped.Load(),
		stats.MediaFailed.Load(),
	))
	if stats.HasErrors() {
		sb.WriteString(fmt.Sprintf("Errors: `%d`\n", stats.Errors.Load()))
	}
	return sb.String()
}
