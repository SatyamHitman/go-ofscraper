// =============================================================================
// FILE: internal/cache/json.go
// PURPOSE: JSON file-backed cache implementation. Stores cache entries as
//          individual JSON files in a directory. Suitable for small caches
//          that need human-readable persistence. Ports the JSON cache path
//          in Python utils/cache.py.
// =============================================================================

package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ---------------------------------------------------------------------------
// JSON cache
// ---------------------------------------------------------------------------

// jsonEntry is the on-disk format for a cached entry.
type jsonEntry struct {
	Key       string    `json:"key"`
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// jsonCache stores each entry as a separate JSON file in a directory.
type jsonCache struct {
	dir string
}

// newJSONCache creates a JSON file-backed cache in the given directory.
//
// Parameters:
//   - dir: Directory for cache files.
//
// Returns:
//   - A Cache backed by JSON files, and any error.
func newJSONCache(dir string) (Cache, error) {
	cacheDir := filepath.Join(dir, "json_cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create json cache dir %s: %w", cacheDir, err)
	}
	return &jsonCache{dir: cacheDir}, nil
}

// Get reads and returns the cached value if it exists and is not expired.
func (jc *jsonCache) Get(key string) ([]byte, bool) {
	path := jc.keyPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry jsonEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		_ = os.Remove(path)
		return nil, false
	}

	return entry.Data, true
}

// Set writes a cache entry to a JSON file.
func (jc *jsonCache) Set(key string, value []byte, expiry time.Duration) error {
	entry := jsonEntry{
		Key:       key,
		Data:      value,
		CreatedAt: time.Now(),
	}
	if expiry > 0 {
		entry.ExpiresAt = time.Now().Add(expiry)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	path := jc.keyPath(key)
	return os.WriteFile(path, data, 0644)
}

// Delete removes a cache file.
func (jc *jsonCache) Delete(key string) error {
	path := jc.keyPath(key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Has checks if a key exists and is not expired.
func (jc *jsonCache) Has(key string) bool {
	_, ok := jc.Get(key)
	return ok
}

// Clear removes all cache files.
func (jc *jsonCache) Clear() error {
	entries, err := os.ReadDir(jc.dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			_ = os.Remove(filepath.Join(jc.dir, entry.Name()))
		}
	}
	return nil
}

// Close is a no-op for the JSON cache.
func (jc *jsonCache) Close() error {
	return nil
}

// keyPath converts a cache key to a filename using SHA-256 to avoid
// filesystem-unsafe characters.
func (jc *jsonCache) keyPath(key string) string {
	h := sha256.Sum256([]byte(key))
	name := hex.EncodeToString(h[:16]) + ".json" // Use first 16 bytes.
	return filepath.Join(jc.dir, name)
}
