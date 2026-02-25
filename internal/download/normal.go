// =============================================================================
// FILE: internal/download/normal.go
// PURPOSE: Normal HTTP download handler. Downloads unprotected media files
//          using direct HTTP GET with chunked transfer, resume support, and
//          progress tracking. Ports Python managers/main_download.py.
// =============================================================================

package download

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Normal download
// ---------------------------------------------------------------------------

// downloadNormal downloads a non-DRM media file via direct HTTP.
func (o *Orchestrator) downloadNormal(ctx context.Context, m *model.Media) error {
	if m.RawURL == "" {
		return fmt.Errorf("no download URL")
	}

	// Determine output path.
	outputPath := m.FilePath
	if outputPath == "" {
		return fmt.Errorf("no output path set for media %d", m.ID)
	}

	// Ensure parent directory exists.
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Check for resume.
	var startByte int64
	if o.cfg.ResumeEnabled {
		if info, err := os.Stat(outputPath + ".part"); err == nil {
			startByte = info.Size()
		}
	}

	// Build request.
	req := gohttp.NewRequest(m.RawURL)
	if startByte > 0 {
		req = req.WithHeader("Range", fmt.Sprintf("bytes=%d-", startByte))
	}

	resp, err := o.session.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Close()

	if !resp.IsOK() && resp.StatusCode != 206 {
		return fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	// Open output file.
	partPath := outputPath + ".part"
	flags := os.O_CREATE | os.O_WRONLY
	if startByte > 0 && resp.StatusCode == 206 {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
		startByte = 0
	}

	f, err := os.OpenFile(partPath, flags, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Copy with optional speed limit.
	var reader io.Reader = resp.Body
	if o.cfg.SpeedLimit > 0 {
		reader = NewSpeedLimitReader(resp.Body, o.cfg.SpeedLimit)
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// Promote .part to final path.
	if err := os.Rename(partPath, outputPath); err != nil {
		return fmt.Errorf("rename file: %w", err)
	}

	return nil
}
