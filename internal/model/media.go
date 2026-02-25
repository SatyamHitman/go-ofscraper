// =============================================================================
// FILE: internal/model/media.go
// PURPOSE: Defines the Media domain model representing a single media item
//          (image, video, audio) within a post. Ports Python classes/of/media.py
//          with all fields, computed properties, and status tracking.
// =============================================================================

package model

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Enumerations
// ---------------------------------------------------------------------------

// MediaType classifies the kind of media content.
type MediaType string

const (
	MediaTypeImages       MediaType = "images"
	MediaTypeVideos       MediaType = "videos"
	MediaTypeAudios       MediaType = "audios"
	MediaTypeTexts        MediaType = "texts"
	MediaTypeForcedSkipped MediaType = "forced_skipped"
)

// DownloadType indicates whether the media requires DRM decryption.
type DownloadType string

const (
	DownloadTypeNormal    DownloadType = "Normal"
	DownloadTypeProtected DownloadType = "Protected"
)

// DownloadStatus tracks the outcome of a download attempt.
type DownloadStatus string

const (
	DownloadStatusSucceeded DownloadStatus = "succeeded"
	DownloadStatusFailed    DownloadStatus = "failed"
	DownloadStatusSkipped   DownloadStatus = "skipped"
)

// MetadataStatus tracks the outcome of a metadata update attempt.
type MetadataStatus string

const (
	MetadataStatusChanged   MetadataStatus = "changed"
	MetadataStatusUnchanged MetadataStatus = "unchanged"
	MetadataStatusFailed    MetadataStatus = "failed"
	MetadataStatusSkipped   MetadataStatus = "skipped"
)

// ---------------------------------------------------------------------------
// Media struct
// ---------------------------------------------------------------------------

// Media represents a single downloadable media item associated with a Post.
// It holds the raw API data, computed properties for download/display, and
// mutable status fields for tracking download and metadata operations.
type Media struct {
	// --- Identification ---
	ID     int64 `json:"id"`
	PostID int64 `json:"post_id"`

	// --- Content source ---
	RawURL   string `json:"url,omitempty"`    // Direct download URL (normal content)
	MpdURL   string `json:"mpd,omitempty"`    // DASH manifest URL (DRM content)
	HlsURL   string `json:"hls,omitempty"`    // HLS manifest URL (alternative DRM)
	Type     string `json:"type,omitempty"`   // Raw type from API ("photo", "video", "audio", "gif")
	Duration string `json:"duration,omitempty"` // Duration timestamp string

	// --- Parent reference ---
	Post     *Post  `json:"-"`         // Parent post (set after construction)
	Username string `json:"username"`  // Creator's username (from parent post)
	ModelID  int64  `json:"model_id"`  // Creator's ID (from parent post)

	// --- Position ---
	Count int `json:"count"` // 1-based index within the parent post's media array

	// --- Content metadata ---
	CanView bool    `json:"canview"` // User has permission to view this media
	Preview int     `json:"preview"` // 1 if parent post is a preview, else 0
	Expires bool    `json:"expires"` // Whether the media has an expiration
	Size    float64 `json:"size"`    // File size in bytes (populated after HEAD request)

	// --- Dates ---
	CreatedAt string `json:"created_at,omitempty"` // Media creation timestamp
	PostedAt  string `json:"posted_at,omitempty"`  // Post publish timestamp

	// --- DRM / CloudFront fields ---
	Policy       string `json:"policy,omitempty"`
	KeyPair      string `json:"keypair,omitempty"`
	Signature    string `json:"signature,omitempty"`
	HlsPolicy    string `json:"hls_policy,omitempty"`
	HlsKeyPair   string `json:"hls_keypair,omitempty"`
	HlsSignature string `json:"hls_signature,omitempty"`

	// --- Post-level inherited fields ---
	ResponseType string `json:"response_type"` // "timeline", "archived", "pinned", etc.
	Label        string `json:"label,omitempty"`
	Value        string `json:"value"`  // "free" or "paid"
	Mass         bool   `json:"mass"`   // From queue/mass message
	Text         string `json:"text,omitempty"` // Parent post text

	// --- Mutable status fields ---
	DownloadAttempted bool            `json:"download_attempted"`
	DownloadSucceeded *bool           `json:"download_succeeded"` // nil = not attempted, true/false = result
	MetadataAttempted bool            `json:"metadata_attempted"`
	MetadataSucceeded *bool           `json:"metadata_succeeded"` // nil = unchanged, true = changed, false = failed

	// --- Filepath (set after path resolution) ---
	FilePath string `json:"filepath,omitempty"`

	// --- Quality selection cache ---
	SelectedQuality string `json:"selected_quality,omitempty"`

	// --- Internal sync ---
	mu sync.Mutex `json:"-"`
}

// ---------------------------------------------------------------------------
// Computed properties
// ---------------------------------------------------------------------------

// MediaType returns the normalized media type classification based on the raw
// API type string. Maps API values like "photo"/"gif" to MediaTypeImages, etc.
//
// Returns:
//   - The normalized MediaType enum value.
func (m *Media) MediaType() MediaType {
	switch strings.ToLower(m.Type) {
	case "photo", "gif":
		return MediaTypeImages
	case "video":
		return MediaTypeVideos
	case "audio":
		return MediaTypeAudios
	case "text":
		return MediaTypeTexts
	default:
		if m.Type == "" {
			return MediaTypeForcedSkipped
		}
		return MediaTypeForcedSkipped
	}
}

// DownloadKind returns whether this media requires DRM decryption or is a
// normal direct download. Protected content has an MPD manifest URL.
//
// Returns:
//   - DownloadTypeProtected if MpdURL is set, DownloadTypeNormal otherwise.
func (m *Media) DownloadKind() DownloadType {
	if m.MpdURL != "" {
		return DownloadTypeProtected
	}
	return DownloadTypeNormal
}

// IsProtected returns true if the media requires DRM decryption.
func (m *Media) IsProtected() bool {
	return m.MpdURL != ""
}

// IsLinked returns true if the media has a downloadable URL (direct or MPD).
func (m *Media) IsLinked() bool {
	return m.RawURL != "" || m.MpdURL != ""
}

// Link returns the primary download URL. Prefers direct URL over MPD.
//
// Returns:
//   - The download URL string, or empty if no URL is available.
func (m *Media) Link() string {
	if m.RawURL != "" {
		return m.RawURL
	}
	return m.MpdURL
}

// Filename extracts the filename from the download URL, stripping query params.
//
// Returns:
//   - The filename string, or empty if no URL is available.
func (m *Media) Filename() string {
	link := m.RawURL
	if link == "" {
		return ""
	}

	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	return path.Base(u.Path)
}

// ContentTypeExt returns the default file extension for this media type.
//
// Returns:
//   - File extension string: "mp4", "jpg", or "mp3".
func (m *Media) ContentTypeExt() string {
	switch m.MediaType() {
	case MediaTypeVideos:
		return "mp4"
	case MediaTypeAudios:
		return "mp3"
	default:
		return "jpg"
	}
}

// Date returns the best available date string for this media.
// Priority: CreatedAt > PostedAt > parent Post date.
//
// Returns:
//   - A date string in the format available from the API.
func (m *Media) Date() string {
	if m.CreatedAt != "" {
		return m.CreatedAt
	}
	if m.PostedAt != "" {
		return m.PostedAt
	}
	if m.Post != nil {
		return m.Post.Date()
	}
	return ""
}

// FormattedDate parses and formats the media date into "2006-01-02 15:04:05".
//
// Returns:
//   - The formatted date string, or the raw date if parsing fails.
func (m *Media) FormattedDate() string {
	raw := m.Date()
	if raw == "" {
		return ""
	}

	// Try common API date formats
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
	} {
		t, err := time.Parse(layout, raw)
		if err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	return raw
}

// LicenseURL constructs the DRM license URL for protected content.
//
// Parameters:
//   - baseTemplate: The URL template with format placeholders.
//   - mediaID: The media identifier.
//   - drmID: The DRM identifier.
//   - drmType: The DRM type string.
//
// Returns:
//   - The constructed license URL string.
func (m *Media) LicenseURL(baseTemplate string, drmID, drmType string) string {
	return fmt.Sprintf(baseTemplate, m.ID, drmID, drmType)
}

// ---------------------------------------------------------------------------
// Status mutation methods
// ---------------------------------------------------------------------------

// MarkDownloadSucceeded records a successful download.
func (m *Media) MarkDownloadSucceeded() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DownloadAttempted = true
	t := true
	m.DownloadSucceeded = &t
}

// MarkDownloadFailed records a failed download.
func (m *Media) MarkDownloadFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DownloadAttempted = true
	f := false
	m.DownloadSucceeded = &f
}

// MarkDownloadSkipped records that the download was intentionally skipped.
func (m *Media) MarkDownloadSkipped() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DownloadAttempted = false
	m.DownloadSucceeded = nil
}

// DownloadStatusString returns the human-readable download status.
//
// Returns:
//   - "succeeded", "failed", or "skipped".
func (m *Media) DownloadStatusString() DownloadStatus {
	if !m.DownloadAttempted {
		return DownloadStatusSkipped
	}
	if m.DownloadSucceeded != nil && *m.DownloadSucceeded {
		return DownloadStatusSucceeded
	}
	return DownloadStatusFailed
}

// MarkMetadataChanged records that metadata was successfully updated.
func (m *Media) MarkMetadataChanged() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetadataAttempted = true
	t := true
	m.MetadataSucceeded = &t
}

// MarkMetadataUnchanged records that metadata was checked but not changed.
func (m *Media) MarkMetadataUnchanged() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetadataAttempted = true
	m.MetadataSucceeded = nil
}

// MarkMetadataFailed records that the metadata update failed.
func (m *Media) MarkMetadataFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MetadataAttempted = true
	f := false
	m.MetadataSucceeded = &f
}

// MetadataStatusString returns the human-readable metadata status.
//
// Returns:
//   - "changed", "unchanged", "failed", or "skipped".
func (m *Media) MetadataStatusString() MetadataStatus {
	if !m.MetadataAttempted {
		return MetadataStatusSkipped
	}
	if m.MetadataSucceeded == nil {
		return MetadataStatusUnchanged
	}
	if *m.MetadataSucceeded {
		return MetadataStatusChanged
	}
	return MetadataStatusFailed
}
