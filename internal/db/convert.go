// =============================================================================
// FILE: internal/db/convert.go
// PURPOSE: Data conversion utilities for database operations. Converts between
//          domain model types and database row types. Ports Python
//          db/utils/convert.py.
// =============================================================================

package db

import (
	"database/sql"
)

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

// NullString creates a sql.NullString from a Go string.
//
// Parameters:
//   - s: The string value.
//
// Returns:
//   - A NullString that is valid if s is non-empty.
func NullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// NullInt64 creates a sql.NullInt64 from a Go int64.
//
// Parameters:
//   - n: The int64 value.
//
// Returns:
//   - A NullInt64 that is valid if n is non-zero.
func NullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: n,
		Valid: n != 0,
	}
}

// NullFloat64 creates a sql.NullFloat64 from a Go float64.
//
// Parameters:
//   - f: The float64 value.
//
// Returns:
//   - A NullFloat64 that is valid if f is non-zero.
func NullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   f != 0,
	}
}

// StringFromNull extracts a string from a NullString, returning empty if null.
//
// Parameters:
//   - ns: The NullString to read.
//
// Returns:
//   - The string value, or "" if null.
func StringFromNull(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Int64FromNull extracts an int64 from a NullInt64, returning 0 if null.
//
// Parameters:
//   - ni: The NullInt64 to read.
//
// Returns:
//   - The int64 value, or 0 if null.
func Int64FromNull(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}
