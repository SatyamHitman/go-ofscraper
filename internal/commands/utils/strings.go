// =============================================================================
// FILE: internal/commands/utils/strings.go
// PURPOSE: Command output string constants and formatters. Provides consistent
//          messaging templates used across all commands for user-facing output.
//          Ports Python utils/args/strings.py output templates.
// =============================================================================

package cmdutils

import "fmt"

// ---------------------------------------------------------------------------
// Output string constants
// ---------------------------------------------------------------------------

const (
	// MsgStarting is displayed when a command begins execution.
	MsgStarting = "Starting %s..."

	// MsgCompleted is displayed when a command finishes successfully.
	MsgCompleted = "%s completed successfully"

	// MsgNoUsers is displayed when no users match the configured filters.
	MsgNoUsers = "No users matched the current filters"

	// MsgNoPosts is displayed when no posts are found for a user.
	MsgNoPosts = "No posts found for %s"

	// MsgNoMedia is displayed when no downloadable media is found.
	MsgNoMedia = "No downloadable media found for %s"

	// MsgProcessingUser is displayed when beginning to process a user.
	MsgProcessingUser = "Processing user: %s (%d/%d)"

	// MsgFetchingPosts is displayed when fetching posts for a content area.
	MsgFetchingPosts = "Fetching %s posts for %s..."

	// MsgDownloadProgress shows download progress for a user.
	MsgDownloadProgress = "Downloaded %d/%d media for %s"

	// MsgSkippedUser is displayed when a user is skipped due to filters.
	MsgSkippedUser = "Skipped user: %s (reason: %s)"

	// MsgDaemonSleep is displayed when daemon mode pauses between runs.
	MsgDaemonSleep = "Daemon sleeping for %s until next run"

	// MsgSummaryHeader is the header for the final summary section.
	MsgSummaryHeader = "===== Scrape Summary ====="

	// MsgSummaryFooter is the footer for the final summary section.
	MsgSummaryFooter = "========================="
)

// ---------------------------------------------------------------------------
// Formatters
// ---------------------------------------------------------------------------

// FormatUserProgress returns a formatted progress string for user processing.
//
// Parameters:
//   - username: The user being processed.
//   - current: The current user number (1-based).
//   - total: The total number of users.
//
// Returns:
//   - A formatted progress string.
func FormatUserProgress(username string, current, total int) string {
	return fmt.Sprintf(MsgProcessingUser, username, current, total)
}

// FormatSummaryLine returns a formatted key-value line for summary output.
//
// Parameters:
//   - label: The label for the metric.
//   - value: The metric value.
//
// Returns:
//   - A formatted summary line.
func FormatSummaryLine(label string, value interface{}) string {
	return fmt.Sprintf("  %-20s %v", label+":", value)
}
