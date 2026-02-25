// =============================================================================
// FILE: internal/download/drm_download.go
// PURPOSE: DRM/DASH download handler. Downloads protected media requiring
//          DRM decryption via the DRM manager. Ports Python
//          managers/alt_download.py.
// =============================================================================

package download

import (
	"context"
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// DRM download
// ---------------------------------------------------------------------------

// downloadProtected handles a DRM-protected media download.
func (o *Orchestrator) downloadProtected(ctx context.Context, m *model.Media) error {
	if o.drm == nil || !o.drm.IsEnabled() {
		return fmt.Errorf("DRM decryption not configured")
	}

	if m.MpdURL == "" {
		return fmt.Errorf("no MPD URL for protected media %d", m.ID)
	}

	outputPath := m.FilePath
	if outputPath == "" {
		return fmt.Errorf("no output path set for media %d", m.ID)
	}

	// Build license URL from media info.
	licenseURL := fmt.Sprintf("https://onlyfans.com/api2/v2/users/media/%d/drm/post/widevine", m.ID)

	result, err := o.drm.Decrypt(ctx, m.MpdURL, licenseURL, outputPath)
	if err != nil {
		return fmt.Errorf("DRM decrypt: %w", err)
	}

	if o.logger != nil {
		o.logger.Info("DRM download complete",
			"media_id", m.ID,
			"output", result.OutputPath,
			"kid", result.KID,
		)
	}

	return nil
}
