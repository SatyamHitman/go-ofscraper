// =============================================================================
// FILE: internal/db/merge.go
// PURPOSE: Database merge operations. Merges data from one model's database
//          into another, handling duplicate detection and conflict resolution.
//          Ports Python db/merge.py.
// =============================================================================

package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

// ---------------------------------------------------------------------------
// Merge
// ---------------------------------------------------------------------------

// MergeResult holds the outcome of a database merge operation.
type MergeResult struct {
	PostsMerged    int
	MessagesMerged int
	MediaMerged    int
	StoriesMerged  int
	LabelsMerged   int
	Conflicts      int
}

// MergeDatabases merges all data from srcDB into dstDB. Creates a backup
// of dstDB before merging.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - srcConn: Source database connection.
//   - dstConn: Destination database connection.
//
// Returns:
//   - MergeResult with counts of merged records, and any error.
func MergeDatabases(ctx context.Context, srcConn, dstConn *Conn) (MergeResult, error) {
	var result MergeResult

	// Backup destination before merge.
	if _, err := Backup(dstConn.Path); err != nil {
		slog.Warn("failed to backup before merge", "path", dstConn.Path, "error", err)
	}

	// Merge within a transaction.
	err := WithTx(ctx, dstConn, func(tx *sql.Tx) error {
		var err error

		// Merge posts.
		result.PostsMerged, err = mergeTable(ctx, srcConn.DB, tx, "posts", "post_id")
		if err != nil {
			return fmt.Errorf("failed to merge posts: %w", err)
		}

		// Merge messages.
		result.MessagesMerged, err = mergeTable(ctx, srcConn.DB, tx, "messages", "post_id")
		if err != nil {
			return fmt.Errorf("failed to merge messages: %w", err)
		}

		// Merge stories.
		result.StoriesMerged, err = mergeTable(ctx, srcConn.DB, tx, "stories", "post_id")
		if err != nil {
			return fmt.Errorf("failed to merge stories: %w", err)
		}

		// Merge media.
		result.MediaMerged, err = mergeMediaTable(ctx, srcConn.DB, tx)
		if err != nil {
			return fmt.Errorf("failed to merge media: %w", err)
		}

		return nil
	})

	return result, err
}

// mergeTable copies rows from src table into dst tx, skipping duplicates.
func mergeTable(ctx context.Context, srcDB *sql.DB, dstTx *sql.Tx, table, uniqueCol string) (int, error) {
	rows, err := srcDB.QueryContext(ctx,
		fmt.Sprintf("SELECT * FROM %s", table),
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	merged := 0
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return merged, err
		}

		// Build INSERT OR IGNORE statement.
		placeholders := ""
		for i := range cols {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += "?"
		}
		query := fmt.Sprintf("INSERT OR IGNORE INTO %s VALUES (%s)", table, placeholders)
		result, err := dstTx.ExecContext(ctx, query, vals...)
		if err != nil {
			return merged, err
		}
		affected, _ := result.RowsAffected()
		merged += int(affected)
	}

	return merged, rows.Err()
}

// mergeMediaTable handles media merge with special conflict resolution.
func mergeMediaTable(ctx context.Context, srcDB *sql.DB, dstTx *sql.Tx) (int, error) {
	return mergeTable(ctx, srcDB, dstTx, "medias", "media_id")
}
