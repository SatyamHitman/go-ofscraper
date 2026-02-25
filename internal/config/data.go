// =============================================================================
// FILE: internal/config/data.go
// PURPOSE: Provides typed accessor functions for reading individual config
//          values. Each function reads from the global AppConfig with fallback
//          to defaults. Ports Python utils/config/data.py getter functions.
// =============================================================================

package config

import "gofscraper/internal/config/env"

// ---------------------------------------------------------------------------
// Profile & Metadata accessors
// ---------------------------------------------------------------------------

// GetMainProfile returns the active profile name from config.
//
// Returns:
//   - The profile name string.
func GetMainProfile() string {
	cfg := Get()
	if cfg.MainProfile == "" {
		return DefaultProfile
	}
	return cfg.MainProfile
}

// GetMetadata returns the metadata/database path format template.
//
// Returns:
//   - The metadata path template string.
func GetMetadata() string {
	cfg := Get()
	if cfg.Metadata == "" {
		return DefaultMetadataFormat
	}
	return cfg.Metadata
}

// GetDiscord returns the Discord webhook URL.
//
// Returns:
//   - The webhook URL string, or empty if disabled.
func GetDiscord() string {
	return Get().Discord
}

// ---------------------------------------------------------------------------
// File options accessors
// ---------------------------------------------------------------------------

// GetSaveLocation returns the base save location for downloaded content.
//
// Returns:
//   - The save location path string.
func GetSaveLocation() string {
	cfg := Get()
	if cfg.File.SaveLocation == "" {
		return env.DefaultSavePath()
	}
	return cfg.File.SaveLocation
}

// GetDirFormat returns the directory path format template.
//
// Returns:
//   - The directory format string.
func GetDirFormat() string {
	cfg := Get()
	if cfg.File.DirFormat == "" {
		return DefaultDirFormat
	}
	return cfg.File.DirFormat
}

// GetFileFormat returns the filename format template.
//
// Returns:
//   - The filename format string.
func GetFileFormat() string {
	cfg := Get()
	if cfg.File.FileFormat == "" {
		return DefaultFileFormat
	}
	return cfg.File.FileFormat
}

// GetTextLength returns the text truncation length limit.
//
// Returns:
//   - The text length integer (0 = no truncation).
func GetTextLength() int {
	return Get().File.TextLength
}

// GetDateFormat returns the date display format string.
//
// Returns:
//   - The date format string.
func GetDateFormat() string {
	cfg := Get()
	if cfg.File.DateFormat == "" {
		return DefaultDateFormat
	}
	return cfg.File.DateFormat
}

// GetSpaceReplacer returns the character used to replace spaces in paths.
//
// Returns:
//   - The space replacer string.
func GetSpaceReplacer() string {
	return Get().File.SpaceReplacer
}

// GetTextType returns the text truncation mode ("letter" or "word").
//
// Returns:
//   - The text type string.
func GetTextType() string {
	cfg := Get()
	if cfg.File.TextType == "" {
		return DefaultTextType
	}
	return cfg.File.TextType
}

// GetTruncation returns whether path truncation is enabled.
//
// Returns:
//   - true if truncation is enabled.
func GetTruncation() bool {
	return Get().File.Truncation
}

// ---------------------------------------------------------------------------
// Download options accessors
// ---------------------------------------------------------------------------

// GetFilter returns the media type filter list (e.g., ["Images", "Videos"]).
//
// Returns:
//   - The filter string slice.
func GetFilter() []string {
	cfg := Get()
	if len(cfg.Download.Filter) == 0 {
		return MediaFilterOptions
	}
	return cfg.Download.Filter
}

// GetAutoResume returns whether download auto-resume is enabled.
//
// Returns:
//   - true if auto-resume is enabled.
func GetAutoResume() bool {
	return Get().Download.AutoResume
}

// GetSystemFreeSize returns the minimum required free disk space in bytes.
//
// Returns:
//   - The free space threshold in bytes (0 = no check).
func GetSystemFreeSize() int64 {
	return Get().Download.SystemFreeMin
}

// GetMaxPostCount returns the maximum number of posts to process.
//
// Returns:
//   - The max post count (0 = unlimited).
func GetMaxPostCount() int {
	return Get().Download.MaxPostCount
}

// ---------------------------------------------------------------------------
// CDM options accessors
// ---------------------------------------------------------------------------

// GetPrivateKey returns the Widevine private key file path.
//
// Returns:
//   - The file path string, or empty if not configured.
func GetPrivateKey() string {
	return Get().CDM.PrivateKey
}

// GetClientID returns the Widevine client ID file path.
//
// Returns:
//   - The file path string, or empty if not configured.
func GetClientID() string {
	return Get().CDM.ClientID
}

// GetKeyMode returns the CDM key mode ("cdrm" or "manual").
//
// Returns:
//   - The key mode string.
func GetKeyMode() string {
	cfg := Get()
	if cfg.CDM.KeyMode == "" {
		return DefaultKeyMode
	}
	return cfg.CDM.KeyMode
}

// ---------------------------------------------------------------------------
// Performance accessors
// ---------------------------------------------------------------------------

// GetDownloadSemaphores returns the download concurrency semaphore count.
//
// Returns:
//   - The semaphore count integer.
func GetDownloadSemaphores() int {
	cfg := Get()
	if cfg.Performance.DownloadSems <= 0 {
		return DefaultDownloadSem
	}
	return cfg.Performance.DownloadSems
}

// GetDownloadLimit returns the bandwidth download limit in bytes.
//
// Returns:
//   - The limit in bytes (0 = unlimited).
func GetDownloadLimit() int64 {
	return Get().Performance.DownloadLimit
}

// ---------------------------------------------------------------------------
// Content filter accessors
// ---------------------------------------------------------------------------

// GetBlockAds returns whether ad blocking is enabled.
//
// Returns:
//   - true if ad blocking is enabled.
func GetBlockAds() bool {
	return Get().Content.BlockAds
}

// GetFileSizeMax returns the maximum file size filter in bytes.
//
// Returns:
//   - The max size in bytes (0 = no limit).
func GetFileSizeMax() int64 {
	return Get().Content.FileSizeMax
}

// GetFileSizeMin returns the minimum file size filter in bytes.
//
// Returns:
//   - The min size in bytes (0 = no limit).
func GetFileSizeMin() int64 {
	return Get().Content.FileSizeMin
}

// GetMaxMediaLength returns the maximum media duration filter.
//
// Returns:
//   - The max length (0 = no limit).
func GetMaxMediaLength() int {
	return Get().Content.LengthMax
}

// GetMinMediaLength returns the minimum media duration filter.
//
// Returns:
//   - The min length (0 = no limit).
func GetMinMediaLength() int {
	return Get().Content.LengthMin
}

// ---------------------------------------------------------------------------
// Advanced options accessors
// ---------------------------------------------------------------------------

// GetDynamicMode returns the dynamic rule provider name.
//
// Returns:
//   - The provider name string.
func GetDynamicMode() string {
	cfg := Get()
	if cfg.Advanced.DynamicMode == "" {
		return DefaultDynamicRule
	}
	return cfg.Advanced.DynamicMode
}

// GetDownloadBars returns whether download progress bars are shown.
//
// Returns:
//   - true if progress bars are enabled.
func GetDownloadBars() bool {
	return Get().Advanced.DownloadBars
}

// GetCacheMode returns the cache backend mode ("sqlite", "json", "disabled").
//
// Returns:
//   - The cache mode string.
func GetCacheMode() string {
	cfg := Get()
	if cfg.Advanced.CacheMode == "" {
		return DefaultCacheMode
	}
	return cfg.Advanced.CacheMode
}

// GetRotateLogs returns whether log rotation is enabled.
//
// Returns:
//   - true if log rotation is enabled.
func GetRotateLogs() bool {
	return Get().Advanced.RotateLogs
}

// GetSanitizeDB returns whether DB text sanitization is enabled.
//
// Returns:
//   - true if sanitization is enabled.
func GetSanitizeDB() bool {
	return Get().Advanced.SanitizeText
}

// GetTempDir returns the temporary directory override.
//
// Returns:
//   - The temp dir path, or empty for system default.
func GetTempDir() string {
	return Get().Advanced.TempDir
}

// GetHashEnabled returns whether file hash deduplication is enabled.
//
// Returns:
//   - true if hashing is enabled.
func GetHashEnabled() bool {
	return Get().Advanced.RemoveHashMatch
}

// GetEnableAutoAfter returns whether auto-after mode is enabled.
//
// Returns:
//   - true if auto-after is enabled.
func GetEnableAutoAfter() bool {
	return Get().Advanced.EnableAutoAfter
}

// GetDefaultUserList returns the default user list names.
//
// Returns:
//   - The user list string slice.
func GetDefaultUserList() []string {
	cfg := Get()
	if len(cfg.Advanced.DefaultUserList) == 0 {
		return []string{DefaultUserList}
	}
	return cfg.Advanced.DefaultUserList
}

// GetDefaultBlackList returns the default blacklist names.
//
// Returns:
//   - The blacklist string slice.
func GetDefaultBlackList() []string {
	return Get().Advanced.DefaultBlackList
}

// GetSSLVerify returns whether SSL certificate validation is enabled.
//
// Returns:
//   - true if SSL verification is enabled.
func GetSSLVerify() bool {
	return Get().Advanced.SSLVerify
}

// GetFFmpeg returns the FFmpeg binary path.
//
// Returns:
//   - The path string, or empty for auto-detect.
func GetFFmpeg() string {
	return Get().Binary.FFmpeg
}

// ---------------------------------------------------------------------------
// Script options accessors
// ---------------------------------------------------------------------------

// GetAfterActionScript returns the after-action script path.
func GetAfterActionScript() string {
	return Get().Scripts.AfterActionScript
}

// GetPostScript returns the post-processing script path.
func GetPostScript() string {
	return Get().Scripts.PostScript
}

// GetNamingScript returns the custom naming script path.
func GetNamingScript() string {
	return Get().Scripts.NamingScript
}

// GetAfterDownloadScript returns the after-download script path.
func GetAfterDownloadScript() string {
	return Get().Scripts.AfterDownloadScript
}

// GetSkipDownloadScript returns the skip-download check script path.
func GetSkipDownloadScript() string {
	return Get().Scripts.SkipDownloadScript
}

// ---------------------------------------------------------------------------
// Response type accessors
// ---------------------------------------------------------------------------

// GetResponseType returns the display name for a given API response type.
//
// Parameters:
//   - apiType: The API response type string (e.g., "timeline", "messages").
//
// Returns:
//   - The configured display name, or the default mapping.
func GetResponseType(apiType string) string {
	cfg := Get()

	switch apiType {
	case "timeline":
		if cfg.Response.Timeline != "" {
			return cfg.Response.Timeline
		}
	case "message", "messages":
		if cfg.Response.Message != "" {
			return cfg.Response.Message
		}
	case "archived":
		if cfg.Response.Archived != "" {
			return cfg.Response.Archived
		}
	case "paid":
		if cfg.Response.Paid != "" {
			return cfg.Response.Paid
		}
	case "stories":
		if cfg.Response.Stories != "" {
			return cfg.Response.Stories
		}
	case "highlights":
		if cfg.Response.Highlights != "" {
			return cfg.Response.Highlights
		}
	case "profile":
		if cfg.Response.Profile != "" {
			return cfg.Response.Profile
		}
	case "pinned":
		if cfg.Response.Pinned != "" {
			return cfg.Response.Pinned
		}
	case "streams":
		if cfg.Response.Streams != "" {
			return cfg.Response.Streams
		}
	}

	// Fall back to default map
	if name, ok := DefaultResponseTypeMap[apiType]; ok {
		return name
	}
	return apiType
}
