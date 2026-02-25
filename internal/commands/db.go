// =============================================================================
// FILE: internal/commands/db.go
// PURPOSE: Database management command. Provides backup and merge operations
//          for model databases. Ports Python runner/db.py.
// =============================================================================

package commands

import (
	"context"
	"fmt"
	"log/slog"

	"gofscraper/internal/app"
	cmdutils "gofscraper/internal/commands/utils"
	"gofscraper/internal/db"
)

// ---------------------------------------------------------------------------
// DBOperation enumerates the supported database operations.
// ---------------------------------------------------------------------------

// DBOperation identifies which database operation to run.
type DBOperation string

const (
	DBOpBackup DBOperation = "backup"
	DBOpMerge  DBOperation = "merge"
)

// ---------------------------------------------------------------------------
// DBCommand
// ---------------------------------------------------------------------------

// DBCommand handles database management operations.
type DBCommand struct {
	cmdutils.CommandBase
	operation DBOperation
}

// NewDBCommand creates a DBCommand for the given operation.
//
// Parameters:
//   - logger: Structured logger for output.
//   - operation: The database operation to perform.
//
// Returns:
//   - A configured DBCommand.
func NewDBCommand(logger *slog.Logger, operation DBOperation) *DBCommand {
	return &DBCommand{
		CommandBase: cmdutils.NewCommandBase(logger),
		operation:   operation,
	}
}

// Name returns the command name.
func (d *DBCommand) Name() string {
	return fmt.Sprintf("db_%s", d.operation)
}

// Run executes the database command.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance providing config.
//   - args: Command arguments (usernames for backup, source+dest for merge).
//
// Returns:
//   - Error if the operation fails.
func (d *DBCommand) Run(ctx context.Context, a *app.App, args []string) error {
	d.LogStart(d.Name(), args)
	defer d.LogDone(d.Name())

	switch d.operation {
	case DBOpBackup:
		return d.runBackup(ctx, a, args)
	case DBOpMerge:
		return d.runMerge(ctx, a, args)
	default:
		return fmt.Errorf("unknown db operation: %s", d.operation)
	}
}

// runBackup creates backups for the specified model databases.
func (d *DBCommand) runBackup(_ context.Context, _ *app.App, usernames []string) error {
	if len(usernames) == 0 {
		return fmt.Errorf("no usernames specified for backup")
	}

	var succeeded, failed int
	for _, username := range usernames {
		conn := db.GetConn(username)
		if conn == nil {
			d.Logger.Warn("no open database for user", "user", username)
			failed++
			continue
		}

		backupPath, err := db.Backup(conn.Path)
		if err != nil {
			d.Logger.Error("backup failed", "user", username, "error", err)
			failed++
			continue
		}

		d.Logger.Info("backup created", "user", username, "path", backupPath)
		succeeded++
	}

	d.Logger.Info("backup operation complete",
		"succeeded", succeeded,
		"failed", failed,
	)
	return nil
}

// runMerge merges one model's database into another.
func (d *DBCommand) runMerge(ctx context.Context, _ *app.App, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("merge requires source and destination usernames")
	}

	srcUsername := args[0]
	dstUsername := args[1]

	srcConn := db.GetConn(srcUsername)
	if srcConn == nil {
		return fmt.Errorf("no open database for source user %q", srcUsername)
	}

	dstConn := db.GetConn(dstUsername)
	if dstConn == nil {
		return fmt.Errorf("no open database for destination user %q", dstUsername)
	}

	d.Logger.Info("starting database merge",
		"source", srcUsername,
		"destination", dstUsername,
	)

	result, err := db.MergeDatabases(ctx, srcConn, dstConn)
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	d.Logger.Info("merge complete",
		"posts_merged", result.PostsMerged,
		"messages_merged", result.MessagesMerged,
		"media_merged", result.MediaMerged,
		"stories_merged", result.StoriesMerged,
		"conflicts", result.Conflicts,
	)

	return nil
}
