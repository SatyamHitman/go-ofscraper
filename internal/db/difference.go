// =============================================================================
// FILE: internal/db/difference.go
// PURPOSE: Database difference/comparison operations. Compares two model
//          databases and reports records that exist in one but not the other.
//          Used for migration validation and debugging. Ports Python
//          db/difference.py.
// =============================================================================

package db

import (
	"context"
	"database/sql"
	"fmt"
)

// ---------------------------------------------------------------------------
// Difference
// ---------------------------------------------------------------------------

// DiffResult holds the difference between two databases.
type DiffResult struct {
	OnlyInSource []int64 // IDs present in source but not destination.
	OnlyInDest   []int64 // IDs present in destination but not source.
	InBoth       int     // Count of IDs present in both.
}

// DiffTable compares a table between two databases by a unique ID column.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - srcDB: Source database.
//   - dstDB: Destination database.
//   - table: Table name to compare.
//   - idCol: The unique ID column name.
//
// Returns:
//   - DiffResult with the comparison, and any error.
func DiffTable(ctx context.Context, srcDB, dstDB *sql.DB, table, idCol string) (DiffResult, error) {
	var result DiffResult

	srcIDs, err := collectIDs(ctx, srcDB, table, idCol)
	if err != nil {
		return result, fmt.Errorf("failed to read source %s: %w", table, err)
	}

	dstIDs, err := collectIDs(ctx, dstDB, table, idCol)
	if err != nil {
		return result, fmt.Errorf("failed to read dest %s: %w", table, err)
	}

	// Build sets for comparison.
	srcSet := make(map[int64]struct{}, len(srcIDs))
	for _, id := range srcIDs {
		srcSet[id] = struct{}{}
	}

	dstSet := make(map[int64]struct{}, len(dstIDs))
	for _, id := range dstIDs {
		dstSet[id] = struct{}{}
	}

	// Find differences.
	for _, id := range srcIDs {
		if _, ok := dstSet[id]; ok {
			result.InBoth++
		} else {
			result.OnlyInSource = append(result.OnlyInSource, id)
		}
	}
	for _, id := range dstIDs {
		if _, ok := srcSet[id]; !ok {
			result.OnlyInDest = append(result.OnlyInDest, id)
		}
	}

	return result, nil
}

// collectIDs reads all ID values from a table column.
func collectIDs(ctx context.Context, db *sql.DB, table, idCol string) ([]int64, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", idCol, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
