// =============================================================================
// FILE: internal/commands/check.go
// PURPOSE: Check command implementation. Runs check operations for messages,
//          stories, paid content, and posts. Fetches content and displays
//          results in a table format. Ports Python runner/check.py.
// =============================================================================

package commands

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"gofscraper/internal/app"
	cmdutils "gofscraper/internal/commands/utils"
)

// ---------------------------------------------------------------------------
// CheckType enumerates the supported check operations.
// ---------------------------------------------------------------------------

// CheckType identifies which content type to check.
type CheckType string

const (
	CheckMessages CheckType = "msg_check"
	CheckStories  CheckType = "story_check"
	CheckPaid     CheckType = "paid_check"
	CheckPosts    CheckType = "post_check"
)

// AllCheckTypes returns all supported check types.
func AllCheckTypes() []CheckType {
	return []CheckType{CheckMessages, CheckStories, CheckPaid, CheckPosts}
}

// ---------------------------------------------------------------------------
// CheckCommand
// ---------------------------------------------------------------------------

// CheckCommand runs check operations that fetch and display content summaries.
type CheckCommand struct {
	cmdutils.CommandBase
	checkType CheckType
}

// NewCheckCommand creates a CheckCommand for the given check type.
//
// Parameters:
//   - logger: Structured logger for output.
//   - checkType: The type of check to run.
//
// Returns:
//   - A configured CheckCommand.
func NewCheckCommand(logger *slog.Logger, checkType CheckType) *CheckCommand {
	return &CheckCommand{
		CommandBase: cmdutils.NewCommandBase(logger),
		checkType:   checkType,
	}
}

// Name returns the command name.
func (c *CheckCommand) Name() string {
	return string(c.checkType)
}

// Run executes the check command.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance providing session and config.
//   - usernames: Usernames to check content for.
//
// Returns:
//   - Error if the check operation fails.
func (c *CheckCommand) Run(ctx context.Context, a *app.App, usernames []string) error {
	c.LogStart(c.Name(), usernames)
	defer c.LogDone(c.Name())

	if len(usernames) == 0 {
		c.Logger.Info(cmdutils.MsgNoUsers)
		return nil
	}

	for _, username := range usernames {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		results, err := c.fetchCheckData(ctx, a, username)
		if err != nil {
			c.Logger.Error("check failed", "user", username, "type", c.checkType, "error", err)
			continue
		}

		c.printTable(username, results)
	}

	return nil
}

// ---------------------------------------------------------------------------
// CheckResult holds a single row in the check output table.
// ---------------------------------------------------------------------------

// CheckResult represents one item in the check results.
type CheckResult struct {
	ID       int64
	Date     string
	Text     string
	Price    float64
	HasMedia bool
	Paid     bool
}

// fetchCheckData retrieves check data for a user based on the check type.
func (c *CheckCommand) fetchCheckData(ctx context.Context, a *app.App, username string) ([]CheckResult, error) {
	_ = a.Session() // Will be used when API is wired.

	// TODO: Wire to appropriate API endpoint based on c.checkType.
	switch c.checkType {
	case CheckMessages:
		return c.fetchMessages(ctx, username)
	case CheckStories:
		return c.fetchStories(ctx, username)
	case CheckPaid:
		return c.fetchPaid(ctx, username)
	case CheckPosts:
		return c.fetchPosts(ctx, username)
	default:
		return nil, fmt.Errorf("unknown check type: %s", c.checkType)
	}
}

func (c *CheckCommand) fetchMessages(_ context.Context, _ string) ([]CheckResult, error) {
	// TODO: Wire to messages API.
	return nil, nil
}

func (c *CheckCommand) fetchStories(_ context.Context, _ string) ([]CheckResult, error) {
	// TODO: Wire to stories API.
	return nil, nil
}

func (c *CheckCommand) fetchPaid(_ context.Context, _ string) ([]CheckResult, error) {
	// TODO: Wire to paid content API.
	return nil, nil
}

func (c *CheckCommand) fetchPosts(_ context.Context, _ string) ([]CheckResult, error) {
	// TODO: Wire to posts API.
	return nil, nil
}

// printTable displays check results in a formatted table.
func (c *CheckCommand) printTable(username string, results []CheckResult) {
	if len(results) == 0 {
		c.Logger.Info("no results found", "user", username, "type", c.checkType)
		return
	}

	c.Logger.Info(fmt.Sprintf("Check results for %s (%s): %d items", username, c.checkType, len(results)))

	header := fmt.Sprintf("%-12s %-20s %-8s %-8s %-50s", "ID", "Date", "Price", "Media", "Text")
	c.Logger.Info(header)
	c.Logger.Info(strings.Repeat("-", 100))

	for _, r := range results {
		media := "no"
		if r.HasMedia {
			media = "yes"
		}
		text := r.Text
		if len(text) > 47 {
			text = text[:47] + "..."
		}
		line := fmt.Sprintf("%-12d %-20s $%-7.2f %-8s %-50s", r.ID, r.Date, r.Price, media, text)
		c.Logger.Info(line)
	}
}
