// =============================================================================
// FILE: internal/paths/sanitize.go
// PURPOSE: Path sanitisation and truncation. Ensures generated file and
//          directory names are safe across all platforms (Windows, macOS,
//          Linux) and respect filesystem length limits. Ports Python
//          utils/paths/ sanitisation logic + pathvalidate integration.
// =============================================================================

package paths

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Filename sanitisation
// ---------------------------------------------------------------------------

// Reserved filenames on Windows that are invalid regardless of extension.
var windowsReserved = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true,
	"COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
	"LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

// unsafeCharsRe matches characters that are unsafe in filenames across all
// major operating systems.
var unsafeCharsRe = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

// SanitizeFilename makes a filename safe for all platforms.
//
// Parameters:
//   - name: The raw filename.
//   - spacer: Replacement for spaces ("" to keep spaces, "_" etc.).
//
// Returns:
//   - A sanitised filename.
func SanitizeFilename(name, spacer string) string {
	if name == "" {
		return "_"
	}

	// Replace unsafe characters.
	name = unsafeCharsRe.ReplaceAllString(name, "_")

	// Replace spaces if a spacer is specified.
	if spacer != "" {
		name = strings.ReplaceAll(name, " ", spacer)
	}

	// Remove trailing dots and spaces (Windows issue).
	name = strings.TrimRight(name, ". ")

	// Handle Windows reserved names.
	upper := strings.ToUpper(name)
	// Check without extension.
	base := upper
	if idx := strings.IndexByte(base, '.'); idx >= 0 {
		base = base[:idx]
	}
	if windowsReserved[base] {
		name = "_" + name
	}

	if name == "" {
		return "_"
	}

	return name
}

// SanitizeDirName makes a directory name safe for all platforms.
//
// Parameters:
//   - name: The raw directory name.
//   - spacer: Replacement for spaces.
//
// Returns:
//   - A sanitised directory name.
func SanitizeDirName(name, spacer string) string {
	// Directory names follow the same rules as filenames.
	return SanitizeFilename(name, spacer)
}

// ---------------------------------------------------------------------------
// Path truncation
// ---------------------------------------------------------------------------

// TruncatePath ensures a full file path does not exceed the OS maximum path
// length. Shortens the filename portion while preserving the directory and
// file extension.
//
// Parameters:
//   - fullPath: The full file path.
//
// Returns:
//   - The potentially truncated path.
func TruncatePath(fullPath string) string {
	maxPath := env.MaxPathLength()
	if maxPath <= 0 {
		maxPath = 260 // Windows default
	}

	if len(fullPath) <= maxPath {
		return fullPath
	}

	dir := filepath.Dir(fullPath)
	base := filepath.Base(fullPath)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)

	// Calculate how much space we have for the filename.
	available := maxPath - len(dir) - 1 - len(ext) // -1 for separator
	if available <= 0 {
		// Directory itself is too long; can't truncate meaningfully.
		return fullPath
	}

	nameOnly = truncateUTF8(nameOnly, available)
	return filepath.Join(dir, nameOnly+ext)
}

// TruncateFilename ensures a filename (without directory) does not exceed
// the maximum filename length for the filesystem.
//
// Parameters:
//   - name: The filename to check.
//
// Returns:
//   - The potentially truncated filename.
func TruncateFilename(name string) string {
	maxName := env.MaxFilenameLength()
	if maxName <= 0 {
		maxName = 255 // Common filesystem limit
	}

	if len(name) <= maxName {
		return name
	}

	ext := filepath.Ext(name)
	nameOnly := strings.TrimSuffix(name, ext)

	available := maxName - len(ext)
	if available <= 0 {
		return name[:maxName]
	}

	nameOnly = truncateUTF8(nameOnly, available)
	return nameOnly + ext
}

// truncateUTF8 truncates a string to at most maxBytes bytes on a valid UTF-8
// boundary.
func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	for maxBytes > 0 && !utf8.RuneStart(s[maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
}
