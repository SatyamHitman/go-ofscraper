// =============================================================================
// FILE: internal/drm/license.go
// PURPOSE: License URL construction and key caching. Builds license server
//          URLs from templates and caches decryption keys to avoid repeated
//          license requests for the same content.
// =============================================================================

package drm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// ---------------------------------------------------------------------------
// Key cache
// ---------------------------------------------------------------------------

// KeyCache caches decryption keys by key ID to avoid duplicate license requests.
type KeyCache struct {
	mu    sync.RWMutex
	cache map[string]string // kid -> key (hex)
}

// NewKeyCache creates an empty key cache.
func NewKeyCache() *KeyCache {
	return &KeyCache{
		cache: make(map[string]string),
	}
}

// Get retrieves a cached key by KID.
//
// Parameters:
//   - kid: The key ID (hex).
//
// Returns:
//   - The cached key (hex) and true, or empty string and false if not cached.
func (kc *KeyCache) Get(kid string) (string, bool) {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	key, ok := kc.cache[strings.ToLower(kid)]
	return key, ok
}

// Set stores a key in the cache.
//
// Parameters:
//   - kid: The key ID (hex).
//   - key: The decryption key (hex).
func (kc *KeyCache) Set(kid, key string) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	kc.cache[strings.ToLower(kid)] = key
}

// Size returns the number of cached keys.
func (kc *KeyCache) Size() int {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	return len(kc.cache)
}

// ---------------------------------------------------------------------------
// License URL construction
// ---------------------------------------------------------------------------

// BuildLicenseURL constructs a license server URL from a template and media info.
//
// Parameters:
//   - template: URL template with {media_id} and {drm_type} placeholders.
//   - mediaID: The media identifier.
//   - drmType: The DRM type string (e.g. "widevine").
//
// Returns:
//   - The constructed license URL.
func BuildLicenseURL(template string, mediaID int64, drmType string) string {
	url := strings.ReplaceAll(template, "{media_id}", fmt.Sprintf("%d", mediaID))
	url = strings.ReplaceAll(url, "{drm_type}", drmType)
	return url
}

// ---------------------------------------------------------------------------
// FFmpeg decryption
// ---------------------------------------------------------------------------

// FFmpegDecrypt downloads and decrypts DRM content using FFmpeg with the
// given decryption key.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - ffmpegPath: Path to the FFmpeg binary.
//   - inputURL: The DASH manifest or encrypted content URL.
//   - key: The decryption key (hex).
//   - kid: The key ID (hex).
//   - outputPath: Path for the decrypted output file.
//
// Returns:
//   - Error if the operation fails.
func FFmpegDecrypt(ctx context.Context, ffmpegPath, inputURL, key, kid, outputPath string) error {
	decryptionKey := fmt.Sprintf("%s:%s", kid, key)

	args := []string{
		"-y",
		"-decryption_key", decryptionKey,
		"-i", inputURL,
		"-c", "copy",
		outputPath,
	}

	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\noutput: %s", err, string(output))
	}

	return nil
}
