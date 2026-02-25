// =============================================================================
// FILE: internal/download/paths.go
// PURPOSE: Download path resolution. Constructs output file paths for media
//          items using configured templates and placeholder expansion.
//          Ports Python actions/utils/paths.py.
// =============================================================================

package download

import (
	"fmt"
	"path/filepath"
	"strings"

	"gofscraper/internal/model"
	"gofscraper/internal/paths"
)

// ---------------------------------------------------------------------------
// Path resolution
// ---------------------------------------------------------------------------

// PathConfig holds path template configuration.
type PathConfig struct {
	SaveLocation  string // Base save directory
	DirFormat     string // Directory template (e.g. "{model_username}/{response_type}")
	FileFormat    string // File name template (e.g. "{filename}.{ext}")
	TextFormat    string // Text file name template
	DateFormat    string // Date format string for placeholders
	TruncateLength int   // Max filename length
}

// DefaultPathConfig returns default path configuration.
func DefaultPathConfig() PathConfig {
	return PathConfig{
		SaveLocation:   "",
		DirFormat:      "{model_username}/{response_type}",
		FileFormat:     "{filename}.{ext}",
		TextFormat:     "{model_username}_{post_id}.txt",
		DateFormat:     "2006-01-02",
		TruncateLength: 200,
	}
}

// ResolvePath builds the full output path for a media item.
//
// Parameters:
//   - m: The media item.
//   - cfg: Path configuration.
//
// Returns:
//   - The resolved absolute file path, or error.
func ResolvePath(m *model.Media, cfg PathConfig) (string, error) {
	if cfg.SaveLocation == "" {
		return "", fmt.Errorf("save location not configured")
	}

	// Build directory path from template.
	dir := expandPlaceholders(cfg.DirFormat, m)
	dir = paths.SanitizeDirName(dir, "_")

	// Build filename from template.
	filename := expandPlaceholders(cfg.FileFormat, m)
	filename = paths.SanitizeFilename(filename, "_")

	if cfg.TruncateLength > 0 {
		filename = paths.TruncateFilename(filename)
	}

	return filepath.Join(cfg.SaveLocation, dir, filename), nil
}

// expandPlaceholders replaces template placeholders with media values.
func expandPlaceholders(template string, m *model.Media) string {
	replacer := strings.NewReplacer(
		"{model_username}", safeStr(m.Username),
		"{model_id}", fmt.Sprintf("%d", m.ModelID),
		"{post_id}", fmt.Sprintf("%d", m.PostID),
		"{media_id}", fmt.Sprintf("%d", m.ID),
		"{media_type}", safeStr(m.Type),
		"{response_type}", safeStr(m.ResponseType),
		"{filename}", safeStr(m.Filename()),
		"{ext}", safeStr(m.ContentTypeExt()),
		"{label}", safeStr(m.Label),
		"{value}", safeStr(m.Value),
		"{date}", safeStr(m.FormattedDate()),
	)
	return replacer.Replace(template)
}

// safeStr returns a safe default for empty strings.
func safeStr(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
