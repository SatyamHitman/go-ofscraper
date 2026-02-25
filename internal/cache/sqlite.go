// =============================================================================
// FILE: internal/cache/sqlite.go
// PURPOSE: SQLite-backed cache implementation. Stores cache entries in a
//          SQLite database for durable, high-performance caching. Used as the
//          default cache mode for production. Ports the SQLite cache path
//          in Python utils/cache.py.
// =============================================================================

package cache

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite" // Pure-Go SQLite driver.
)

// ---------------------------------------------------------------------------
// SQLite cache
// ---------------------------------------------------------------------------

// sqliteCache stores cache entries in a SQLite database.
type sqliteCache struct {
	db *sql.DB
}

// newSQLiteCache creates a SQLite-backed cache in the given directory.
//
// Parameters:
//   - dir: Directory where the cache DB file is stored.
//
// Returns:
//   - A Cache backed by SQLite, and any error.
func newSQLiteCache(dir string) (Cache, error) {
	cacheDir := filepath.Join(dir, "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache dir %s: %w", cacheDir, err)
	}

	dbPath := filepath.Join(cacheDir, "cache.db")
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open cache DB: %w", err)
	}

	// Create the cache table if it doesn't exist.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cache (
			key        TEXT PRIMARY KEY,
			data       BLOB NOT NULL,
			expires_at INTEGER,
			created_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create cache table: %w", err)
	}

	// Create index on expires_at for efficient cleanup.
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_cache_expires ON cache(expires_at)`)

	sc := &sqliteCache{db: db}

	// Purge expired entries on startup.
	_ = sc.purgeExpired()

	return sc, nil
}

// Get retrieves a cached value if it exists and is not expired.
func (sc *sqliteCache) Get(key string) ([]byte, bool) {
	var data []byte
	var expiresAt sql.NullInt64

	err := sc.db.QueryRow(
		`SELECT data, expires_at FROM cache WHERE key = ?`, key,
	).Scan(&data, &expiresAt)

	if err != nil {
		return nil, false
	}

	// Check expiry.
	if expiresAt.Valid && expiresAt.Int64 > 0 {
		if time.Now().Unix() > expiresAt.Int64 {
			// Expired — remove and return miss.
			_ = sc.Delete(key)
			return nil, false
		}
	}

	return data, true
}

// Set stores a value with an optional expiry.
func (sc *sqliteCache) Set(key string, value []byte, expiry time.Duration) error {
	var expiresAt *int64
	if expiry > 0 {
		t := time.Now().Add(expiry).Unix()
		expiresAt = &t
	}

	_, err := sc.db.Exec(
		`INSERT OR REPLACE INTO cache (key, data, expires_at, created_at)
		 VALUES (?, ?, ?, ?)`,
		key, value, expiresAt, time.Now().Unix(),
	)
	return err
}

// Delete removes a cache entry.
func (sc *sqliteCache) Delete(key string) error {
	_, err := sc.db.Exec(`DELETE FROM cache WHERE key = ?`, key)
	return err
}

// Has checks if a key exists and is not expired.
func (sc *sqliteCache) Has(key string) bool {
	_, ok := sc.Get(key)
	return ok
}

// Clear removes all cache entries.
func (sc *sqliteCache) Clear() error {
	_, err := sc.db.Exec(`DELETE FROM cache`)
	return err
}

// Close closes the underlying database connection.
func (sc *sqliteCache) Close() error {
	if sc.db != nil {
		return sc.db.Close()
	}
	return nil
}

// purgeExpired removes all expired entries from the cache.
func (sc *sqliteCache) purgeExpired() error {
	_, err := sc.db.Exec(
		`DELETE FROM cache WHERE expires_at IS NOT NULL AND expires_at < ?`,
		time.Now().Unix(),
	)
	return err
}
