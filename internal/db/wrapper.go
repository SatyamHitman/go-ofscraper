// =============================================================================
// FILE: internal/db/wrapper.go
// PURPOSE: Database operation wrappers. Provides retry logic, error wrapping,
//          and logging for database operations. Handles transient SQLite
//          errors (busy, locked). Ports Python db/operations_/wrapper.py.
// =============================================================================

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Retry wrapper
// ---------------------------------------------------------------------------

// maxRetries is the number of times to retry on transient DB errors.
const maxRetries = 3

// retryDelay is the initial delay between retries (doubles each attempt).
const retryDelay = 100 * time.Millisecond

// WithRetry wraps a database operation with retry logic for transient errors
// (SQLITE_BUSY, SQLITE_LOCKED).
//
// Parameters:
//   - ctx: Context for cancellation.
//   - operation: Name of the operation (for logging).
//   - fn: The database function to execute.
//
// Returns:
//   - Error from the final attempt, or nil.
func WithRetry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error
	delay := retryDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Only retry on transient SQLite errors.
		if !isTransientError(lastErr) {
			return fmt.Errorf("%s: %w", operation, lastErr)
		}

		slog.Debug("retrying DB operation",
			"operation", operation,
			"attempt", attempt+1,
			"error", lastErr,
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			delay *= 2
		}
	}

	return fmt.Errorf("%s: max retries exceeded: %w", operation, lastErr)
}

// isTransientError checks if a database error is transient and can be retried.
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "database is locked") ||
		strings.Contains(msg, "SQLITE_BUSY") ||
		strings.Contains(msg, "SQLITE_LOCKED")
}

// ---------------------------------------------------------------------------
// Error wrapper
// ---------------------------------------------------------------------------

// DBError wraps a database error with context about the operation.
type DBError struct {
	Op  string // Operation name (e.g. "UpsertPost", "GetMedia")
	Err error  // Underlying error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("db %s: %v", e.Op, e.Err)
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// IsNotFound reports whether the error indicates no rows were found.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error is sql.ErrNoRows.
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
