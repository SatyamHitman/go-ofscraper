// =============================================================================
// FILE: internal/hash/hash.go
// PURPOSE: File hashing and deduplication using XXHash128. Provides fast,
//          non-cryptographic hashing for detecting duplicate media files and
//          verifying download integrity. Ports Python utils/hash.py.
// =============================================================================

package hash

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/zeebo/xxh3"
)

// ---------------------------------------------------------------------------
// File hashing
// ---------------------------------------------------------------------------

// bufSize is the read buffer size for hashing (1 MB).
const bufSize = 1024 * 1024

// File computes the XXHash128 of a file at the given path.
//
// Parameters:
//   - path: Absolute path to the file.
//
// Returns:
//   - The hex-encoded 128-bit hash string, and any error.
func File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer f.Close()

	return Reader(f)
}

// Reader computes the XXHash128 from an io.Reader.
//
// Parameters:
//   - r: The data source to hash.
//
// Returns:
//   - The hex-encoded 128-bit hash string, and any error.
func Reader(r io.Reader) (string, error) {
	h := xxh3.New()
	buf := make([]byte, bufSize)

	for {
		n, err := r.Read(buf)
		if n > 0 {
			if _, wErr := h.Write(buf[:n]); wErr != nil {
				return "", fmt.Errorf("hash write error: %w", wErr)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read error during hashing: %w", err)
		}
	}

	sum := h.Sum128()
	// Encode as 32-char hex string (128 bits = 16 bytes).
	var hashBytes [16]byte
	hashBytes[0] = byte(sum.Lo >> 56)
	hashBytes[1] = byte(sum.Lo >> 48)
	hashBytes[2] = byte(sum.Lo >> 40)
	hashBytes[3] = byte(sum.Lo >> 32)
	hashBytes[4] = byte(sum.Lo >> 24)
	hashBytes[5] = byte(sum.Lo >> 16)
	hashBytes[6] = byte(sum.Lo >> 8)
	hashBytes[7] = byte(sum.Lo)
	hashBytes[8] = byte(sum.Hi >> 56)
	hashBytes[9] = byte(sum.Hi >> 48)
	hashBytes[10] = byte(sum.Hi >> 40)
	hashBytes[11] = byte(sum.Hi >> 32)
	hashBytes[12] = byte(sum.Hi >> 24)
	hashBytes[13] = byte(sum.Hi >> 16)
	hashBytes[14] = byte(sum.Hi >> 8)
	hashBytes[15] = byte(sum.Hi)

	return hex.EncodeToString(hashBytes[:]), nil
}

// Bytes computes the XXHash128 of a byte slice.
//
// Parameters:
//   - data: The data to hash.
//
// Returns:
//   - The hex-encoded 128-bit hash string.
func Bytes(data []byte) string {
	h := xxh3.New()
	_, _ = h.Write(data)
	sum := h.Sum128()

	var hashBytes [16]byte
	hashBytes[0] = byte(sum.Lo >> 56)
	hashBytes[1] = byte(sum.Lo >> 48)
	hashBytes[2] = byte(sum.Lo >> 40)
	hashBytes[3] = byte(sum.Lo >> 32)
	hashBytes[4] = byte(sum.Lo >> 24)
	hashBytes[5] = byte(sum.Lo >> 16)
	hashBytes[6] = byte(sum.Lo >> 8)
	hashBytes[7] = byte(sum.Lo)
	hashBytes[8] = byte(sum.Hi >> 56)
	hashBytes[9] = byte(sum.Hi >> 48)
	hashBytes[10] = byte(sum.Hi >> 40)
	hashBytes[11] = byte(sum.Hi >> 32)
	hashBytes[12] = byte(sum.Hi >> 24)
	hashBytes[13] = byte(sum.Hi >> 16)
	hashBytes[14] = byte(sum.Hi >> 8)
	hashBytes[15] = byte(sum.Hi)

	return hex.EncodeToString(hashBytes[:])
}

// ---------------------------------------------------------------------------
// Deduplication
// ---------------------------------------------------------------------------

// HashSet is a set of file hashes for deduplication.
type HashSet struct {
	hashes map[string]struct{}
}

// NewHashSet creates an empty hash set.
//
// Returns:
//   - A new HashSet.
func NewHashSet() *HashSet {
	return &HashSet{hashes: make(map[string]struct{})}
}

// Add adds a hash to the set.
//
// Parameters:
//   - hash: The hash string to add.
//
// Returns:
//   - true if the hash was new (not a duplicate).
func (hs *HashSet) Add(hash string) bool {
	if _, exists := hs.hashes[hash]; exists {
		return false
	}
	hs.hashes[hash] = struct{}{}
	return true
}

// Contains checks if a hash is in the set.
//
// Parameters:
//   - hash: The hash to check.
//
// Returns:
//   - true if the hash exists in the set.
func (hs *HashSet) Contains(hash string) bool {
	_, exists := hs.hashes[hash]
	return exists
}

// Len returns the number of unique hashes in the set.
func (hs *HashSet) Len() int {
	return len(hs.hashes)
}
