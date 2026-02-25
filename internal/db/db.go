// =============================================================================
// FILE: internal/db/db.go
// PURPOSE: Database connection pool and initialisation. Manages per-model
//          SQLite database connections with WAL mode, busy timeouts, and
//          schema migration. Ports Python db/__init__.py.
// =============================================================================

package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite" // Pure-Go SQLite driver.
)

// ---------------------------------------------------------------------------
// Connection pool
// ---------------------------------------------------------------------------

var (
	// connPool caches open database connections keyed by model username.
	connPool sync.Map
)

// Conn represents a database connection for a specific model.
type Conn struct {
	DB       *sql.DB
	Username string
	Path     string
}

// Open returns a database connection for the given model. Creates the DB file
// and runs migrations if it doesn't exist. Connections are cached so repeated
// calls for the same model return the same connection.
//
// Parameters:
//   - username: The model username.
//   - dbPath: Absolute path to the SQLite database file.
//
// Returns:
//   - A *Conn wrapping the database, and any error.
func Open(username, dbPath string) (*Conn, error) {
	// Return cached connection if available.
	if cached, ok := connPool.Load(username); ok {
		return cached.(*Conn), nil
	}

	// Ensure directory exists.
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %w", err)
	}

	// Open SQLite with WAL mode and busy timeout.
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", dbPath)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %s: %w", dbPath, err)
	}

	// Set connection pool limits — SQLite doesn't handle concurrent writes.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	// Verify connection.
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database %s: %w", dbPath, err)
	}

	// Run schema migrations.
	if err := migrate(sqlDB); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to migrate database %s: %w", dbPath, err)
	}

	conn := &Conn{
		DB:       sqlDB,
		Username: username,
		Path:     dbPath,
	}

	connPool.Store(username, conn)
	return conn, nil
}

// Close closes a specific model's database connection and removes it from
// the pool.
//
// Parameters:
//   - username: The model username whose connection to close.
//
// Returns:
//   - Error if the close fails.
func Close(username string) error {
	raw, ok := connPool.LoadAndDelete(username)
	if !ok {
		return nil
	}
	conn := raw.(*Conn)
	return conn.DB.Close()
}

// CloseAll closes all open database connections. Called during shutdown.
//
// Returns:
//   - The first error encountered, or nil.
func CloseAll() error {
	var firstErr error
	connPool.Range(func(key, value any) bool {
		conn := value.(*Conn)
		if err := conn.DB.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		connPool.Delete(key)
		return true
	})
	return firstErr
}

// GetConn retrieves a cached connection without opening a new one.
//
// Parameters:
//   - username: The model username.
//
// Returns:
//   - The cached *Conn, or nil if not open.
func GetConn(username string) *Conn {
	raw, ok := connPool.Load(username)
	if !ok {
		return nil
	}
	return raw.(*Conn)
}
