// =============================================================================
// FILE: internal/cache/memory.go
// PURPOSE: In-memory cache implementation. Stores cached entries in a
//          sync.Map with TTL-based expiry. Used as the default/fallback cache
//          and for the "disabled" cache mode. Ports the no-persistence path
//          in Python utils/cache.py.
// =============================================================================

package cache

import (
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Memory cache
// ---------------------------------------------------------------------------

// memoryEntry is a single cached value with optional expiry.
type memoryEntry struct {
	data      []byte
	expiresAt time.Time // Zero value = no expiry.
}

// isExpired reports whether this entry has passed its expiry time.
func (e memoryEntry) isExpired() bool {
	if e.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.expiresAt)
}

// memoryCache is an in-memory Cache backed by sync.Map.
type memoryCache struct {
	store sync.Map
}

// newMemoryCache creates a new in-memory cache.
//
// Returns:
//   - A Cache that stores data in memory only.
func newMemoryCache() Cache {
	return &memoryCache{}
}

// Get retrieves a value if it exists and is not expired.
func (mc *memoryCache) Get(key string) ([]byte, bool) {
	raw, ok := mc.store.Load(key)
	if !ok {
		return nil, false
	}
	entry := raw.(memoryEntry)
	if entry.isExpired() {
		mc.store.Delete(key)
		return nil, false
	}
	// Return a copy to prevent mutation.
	data := make([]byte, len(entry.data))
	copy(data, entry.data)
	return data, true
}

// Set stores a value with an optional expiry.
func (mc *memoryCache) Set(key string, value []byte, expiry time.Duration) error {
	entry := memoryEntry{
		data: make([]byte, len(value)),
	}
	copy(entry.data, value)

	if expiry > 0 {
		entry.expiresAt = time.Now().Add(expiry)
	}

	mc.store.Store(key, entry)
	return nil
}

// Delete removes a key.
func (mc *memoryCache) Delete(key string) error {
	mc.store.Delete(key)
	return nil
}

// Has reports whether a key exists and is not expired.
func (mc *memoryCache) Has(key string) bool {
	_, ok := mc.Get(key)
	return ok
}

// Clear removes all entries.
func (mc *memoryCache) Clear() error {
	mc.store.Range(func(key, _ any) bool {
		mc.store.Delete(key)
		return true
	})
	return nil
}

// Close is a no-op for the memory cache.
func (mc *memoryCache) Close() error {
	return mc.Clear()
}
