// =============================================================================
// FILE: internal/db/empty.go
// PURPOSE: Empty/null handling for database operations. Provides functions to
//          check for empty or null values in query results and supply defaults.
//          Ports Python db/operations_/empty.py.
// =============================================================================

package db

import (
	"context"
	"database/sql"
)

// ---------------------------------------------------------------------------
// Empty checks
// ---------------------------------------------------------------------------

// IsTableEmpty reports whether a table has zero rows.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - conn: Database connection.
//   - table: Table name to check.
//
// Returns:
//   - true if the table is empty, and any error.
func IsTableEmpty(ctx context.Context, conn *Conn, table string) (bool, error) {
	var count int
	err := conn.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM "+table,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// TableRowCount returns the number of rows in a table.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//   - table: Table name.
//
// Returns:
//   - Row count, and any error.
func TableRowCount(ctx context.Context, conn *Conn, table string) (int64, error) {
	var count int64
	err := conn.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM "+table,
	).Scan(&count)
	return count, err
}

// DefaultIfNull returns the value if the NullString is valid, otherwise
// returns the default.
//
// Parameters:
//   - ns: The nullable string.
//   - def: Default value if null.
//
// Returns:
//   - The string value or default.
func DefaultIfNull(ns sql.NullString, def string) string {
	if ns.Valid && ns.String != "" {
		return ns.String
	}
	return def
}

// DefaultIntIfNull returns the value if the NullInt64 is valid, otherwise
// returns the default.
//
// Parameters:
//   - ni: The nullable int.
//   - def: Default value if null.
//
// Returns:
//   - The int value or default.
func DefaultIntIfNull(ni sql.NullInt64, def int64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return def
}
