// =============================================================================
// FILE: internal/cache/cache.go
// PURPOSE: Cache interface and factory. Defines the Cache abstraction used for
//          caching API responses, download state, and other data. Supports
//          multiple backends (SQLite, JSON file, in-memory). Ports Python
//          utils/cache.py interface.
// =============================================================================

package cache

import (
	"time"
)

// ---------------------------------------------------------------------------
// Cache interface
// ---------------------------------------------------------------------------

// Cache provides a key-value store with optional expiry for caching API
// responses, file hashes, and other transient data.
type Cache interface {
	// Get retrieves a value by key. Returns the raw bytes and true if found
	// and not expired, or nil and false otherwise.
	Get(key string) ([]byte, bool)

	// Set stores a value with an optional expiry duration. Use 0 for no expiry.
	Set(key string, value []byte, expiry time.Duration) error

	// Delete removes a key from the cache.
	Delete(key string) error

	// Has reports whether a key exists and is not expired.
	Has(key string) bool

	// Clear removes all entries from the cache.
	Clear() error

	// Close releases any resources held by the cache (file handles, DB conns).
	Close() error
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

// Mode identifies which cache backend to use.
type Mode string

const (
	ModeSQLite  Mode = "sqlite"
	ModeJSON    Mode = "json"
	ModeMemory  Mode = "memory"
	ModeDisable Mode = "disabled"
)

// New creates a cache instance based on the given mode.
//
// Parameters:
//   - mode: The cache backend to use.
//   - dir: Directory for persistent caches (ignored for memory/disabled).
//
// Returns:
//   - A Cache implementation, and any error.
func New(mode Mode, dir string) (Cache, error) {
	switch mode {
	case ModeSQLite:
		return newSQLiteCache(dir)
	case ModeJSON:
		return newJSONCache(dir)
	case ModeMemory:
		return newMemoryCache(), nil
	case ModeDisable:
		return newMemoryCache(), nil // No-op cache that never returns hits.
	default:
		return newMemoryCache(), nil
	}
}
