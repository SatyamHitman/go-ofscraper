// =============================================================================
// FILE: internal/drm/drm.go
// PURPOSE: DRM coordinator. Manages the DRM decryption pipeline: detects
//          protected content, fetches manifests, obtains decryption keys,
//          and delegates to FFmpeg for final decryption. Ports Python
//          alt_download + keyhelpers logic.
// =============================================================================

package drm

import (
	"context"
	"fmt"
	"log/slog"
)

// ---------------------------------------------------------------------------
// DRM mode
// ---------------------------------------------------------------------------

// Mode selects the DRM key acquisition strategy.
type Mode string

const (
	ModeManual Mode = "manual" // Local Widevine CDM device files
	ModeCDRM   Mode = "cdrm"   // Remote CDRM service
	ModeNone   Mode = "none"   // DRM disabled
)

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

// Config holds DRM configuration.
type Config struct {
	Mode       Mode   // Key acquisition mode
	CDRMServer string // CDRM service base URL (for ModeCDRM)
	DevicePath string // Path to Widevine device directory (for ModeManual)
	FFmpegPath string // Path to FFmpeg binary
	Logger     *slog.Logger
}

// DefaultConfig returns a DRM config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Mode:       ModeNone,
		FFmpegPath: "ffmpeg",
	}
}

// ---------------------------------------------------------------------------
// Manager
// ---------------------------------------------------------------------------

// Manager coordinates DRM decryption operations.
type Manager struct {
	cfg    Config
	cdrm   *CDRMClient
	logger *slog.Logger
}

// NewManager creates a DRM manager with the given config.
//
// Parameters:
//   - cfg: DRM configuration.
//
// Returns:
//   - A configured Manager.
func NewManager(cfg Config) *Manager {
	m := &Manager{
		cfg:    cfg,
		logger: cfg.Logger,
	}

	if cfg.Mode == ModeCDRM && cfg.CDRMServer != "" {
		m.cdrm = NewCDRMClient(cfg.CDRMServer)
	}

	return m
}

// ---------------------------------------------------------------------------
// Decryption
// ---------------------------------------------------------------------------

// DecryptResult holds the outcome of a DRM decryption operation.
type DecryptResult struct {
	OutputPath string // Path to the decrypted file
	Key        string // Decryption key used (hex)
	KID        string // Key ID (hex)
}

// Decrypt handles the full DRM decryption pipeline for a media item.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - mpdURL: The DASH manifest URL.
//   - licenseURL: The DRM license URL.
//   - outputPath: Desired output file path.
//
// Returns:
//   - DecryptResult with the output path and key info, or error.
func (m *Manager) Decrypt(ctx context.Context, mpdURL, licenseURL, outputPath string) (*DecryptResult, error) {
	if m.cfg.Mode == ModeNone {
		return nil, fmt.Errorf("DRM decryption is disabled")
	}

	// Step 1: Parse MPD manifest.
	manifest, err := ParseMPD(ctx, mpdURL)
	if err != nil {
		return nil, fmt.Errorf("parse MPD: %w", err)
	}

	// Step 2: Extract key ID from manifest.
	kid := manifest.KeyID()
	if kid == "" {
		return nil, fmt.Errorf("no key ID found in manifest")
	}

	// Step 3: Obtain decryption key.
	var key string
	switch m.cfg.Mode {
	case ModeCDRM:
		if m.cdrm == nil {
			return nil, fmt.Errorf("CDRM client not configured")
		}
		key, err = m.cdrm.GetKey(ctx, licenseURL, kid)
	case ModeManual:
		key, err = ManualDecrypt(ctx, m.cfg.DevicePath, licenseURL, kid)
	default:
		return nil, fmt.Errorf("unsupported DRM mode: %s", m.cfg.Mode)
	}

	if err != nil {
		return nil, fmt.Errorf("obtain key: %w", err)
	}

	// Step 4: Download and decrypt with FFmpeg.
	err = FFmpegDecrypt(ctx, m.cfg.FFmpegPath, mpdURL, key, kid, outputPath)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg decrypt: %w", err)
	}

	return &DecryptResult{
		OutputPath: outputPath,
		Key:        key,
		KID:        kid,
	}, nil
}

// IsEnabled reports whether DRM decryption is configured and available.
func (m *Manager) IsEnabled() bool {
	return m.cfg.Mode != ModeNone
}
