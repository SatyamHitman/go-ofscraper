// =============================================================================
// FILE: internal/model/post.go
// PURPOSE: Defines the Post domain model representing a single content post
//          (timeline, archived, pinned, stream, message, story, highlight).
//          Ports Python classes/of/posts.py with all fields, status tracking,
//          and media collection methods.
// =============================================================================

package model

import (
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Enumerations
// ---------------------------------------------------------------------------

// ResponseType identifies the source/category of a post.
type ResponseType string

const (
	ResponseTimeline   ResponseType = "timeline"
	ResponseArchived   ResponseType = "archived"
	ResponsePinned     ResponseType = "pinned"
	ResponseStreams     ResponseType = "streams"
	ResponseStories    ResponseType = "stories"
	ResponseHighlights ResponseType = "highlights"
	ResponseProfile    ResponseType = "profile"
	ResponseMessages   ResponseType = "messages"
)

// ---------------------------------------------------------------------------
// Post struct
// ---------------------------------------------------------------------------

// Post represents a single content post from the OnlyFans API. It contains
// the post metadata, text content, pricing information, and associated media.
// Posts can be timeline entries, messages, stories, highlights, or other types.
type Post struct {
	// --- Identification ---
	ID      int64  `json:"id"`
	ModelID int64  `json:"model_id"` // Creator's user ID
	Username string `json:"username"` // Creator's username

	// --- Content ---
	RawText string `json:"text,omitempty"`  // Raw post text/caption
	Title   string `json:"title,omitempty"` // Post title (if present)

	// --- Pricing ---
	Price float64 `json:"price"` // Purchase price in dollars (0 = free)

	// --- Flags ---
	Paid     bool `json:"paid"`     // Has media, is open, or price > 0
	Archived bool `json:"archived"` // isArchived flag
	Pinned   bool `json:"pinned"`   // isPinned flag
	Stream   bool `json:"stream"`   // Has streamId (live stream)
	Opened   bool `json:"opened"`   // isOpened (purchased/viewed)
	Favorited bool `json:"favorited"` // isFavorite flag
	Mass     bool `json:"mass"`     // isFromQueue (mass message)
	Preview  bool `json:"preview"`  // Preview post flag
	HasExpiry bool `json:"has_expiry"` // Has expiredAt or expiresAt

	// --- Type classification ---
	RawResponseType string       `json:"raw_response_type,omitempty"` // API's responseType field
	ResponseType    ResponseType `json:"response_type"`               // Normalized type
	Label           string       `json:"label,omitempty"`             // Label name (if from a label)

	// --- Dates ---
	PostedAt  string `json:"posted_at,omitempty"`  // postedAt timestamp from API
	CreatedAt string `json:"created_at,omitempty"` // createdAt timestamp from API

	// --- Media ---
	AllMedia []*Media `json:"all_media,omitempty"` // All media items (including unviewable)

	// --- Sender (for messages) ---
	FromUser int64 `json:"from_user,omitempty"` // Sender's user ID

	// --- Download candidate flags ---
	IsDownloadCandidate bool `json:"is_download_candidate"`
	IsLikeCandidate     bool `json:"is_like_candidate"`
	IsTextCandidate     bool `json:"is_text_candidate"`
	IsMetadataCandidate bool `json:"is_metadata_candidate"`

	// --- Like action status ---
	IsActionableLike bool  `json:"is_actionable_like"`
	LikeAttempted    bool  `json:"like_attempted"`
	LikeSuccess      *bool `json:"like_success"` // nil = not attempted

	// --- Text download status ---
	TextDownloadAttempted bool  `json:"text_download_attempted"`
	TextDownloadSucceeded *bool `json:"text_download_succeeded"` // nil = not attempted

	// --- Filtered media lists (populated after filter pass) ---
	MediaToDownload  []*Media `json:"-"`
	MediaForMetadata []*Media `json:"-"`
}

// ---------------------------------------------------------------------------
// Computed properties
// ---------------------------------------------------------------------------

// Date returns the best available date string for this post.
// Priority: PostedAt > CreatedAt.
//
// Returns:
//   - A date string from the API, or empty if neither is set.
func (p *Post) Date() string {
	if p.PostedAt != "" {
		return p.PostedAt
	}
	return p.CreatedAt
}

// FormattedDate parses and formats the post date into "2006-01-02 15:04:05".
//
// Returns:
//   - The formatted date string, or the raw date if parsing fails.
func (p *Post) FormattedDate() string {
	raw := p.Date()
	if raw == "" {
		return ""
	}

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

// Value returns "free" if the post price is zero, otherwise "paid".
//
// Returns:
//   - "free" or "paid" string.
func (p *Post) Value() string {
	if p.Price == 0 {
		return "free"
	}
	return "paid"
}

// IsRegularTimeline returns true if the post is neither archived nor pinned.
//
// Returns:
//   - true for regular timeline posts.
func (p *Post) IsRegularTimeline() bool {
	return !p.Archived && !p.Pinned
}

// ViewableMedia returns media items that the user has permission to view.
// Filters AllMedia by the CanView field.
//
// Returns:
//   - Slice of viewable Media pointers.
func (p *Post) ViewableMedia() []*Media {
	var result []*Media
	for _, m := range p.AllMedia {
		if m.CanView {
			result = append(result, m)
		}
	}
	return result
}

// DBText returns the post text sanitized for database storage.
// Uses DBCleanup to strip HTML and normalize whitespace.
//
// Parameters:
//   - sanitize: Whether to apply DB sanitization. If false, returns raw text.
//
// Returns:
//   - The cleaned text string.
func (p *Post) DBText(sanitize bool) string {
	if !sanitize {
		return p.RawText
	}
	return DBCleanup(p.RawText)
}

// FileText returns the post text sanitized for use in filenames.
// Uses FileCleanup to remove unsafe characters.
//
// Returns:
//   - The filesystem-safe text string.
func (p *Post) FileText() string {
	return FileCleanup(p.RawText)
}

// LabelString returns the label name or "None" if no label is set.
//
// Returns:
//   - The label string.
func (p *Post) LabelString() string {
	if p.Label == "" {
		return "None"
	}
	return p.Label
}

// ---------------------------------------------------------------------------
// Response type helpers
// ---------------------------------------------------------------------------

// DeriveResponseType computes the normalized ResponseType from the post's
// flags and raw API response type. Checks pinned, archived, and stream flags
// first, then falls back to the raw API response type string.
//
// Returns:
//   - The computed ResponseType enum value.
func (p *Post) DeriveResponseType() ResponseType {
	if p.Pinned {
		return ResponsePinned
	}
	if p.Archived {
		return ResponseArchived
	}
	if p.Stream {
		return ResponseStreams
	}

	switch strings.ToLower(p.RawResponseType) {
	case "timeline", "post":
		return ResponseTimeline
	case "archived":
		return ResponseArchived
	case "pinned":
		return ResponsePinned
	case "stream", "streams":
		return ResponseStreams
	case "stories":
		return ResponseStories
	case "highlights":
		return ResponseHighlights
	case "profile":
		return ResponseProfile
	case "message", "messages":
		return ResponseMessages
	default:
		return ResponseTimeline
	}
}

// ---------------------------------------------------------------------------
// Like action methods
// ---------------------------------------------------------------------------

// MarkPostLiked records that the like action succeeded.
//
// Parameters:
//   - success: Whether the like API call succeeded.
func (p *Post) MarkPostLiked(success bool) {
	p.LikeAttempted = true
	p.LikeSuccess = &success
	if success {
		p.Favorited = true
	}
}

// MarkPostUnliked records that the unlike action succeeded.
//
// Parameters:
//   - success: Whether the unlike API call succeeded.
func (p *Post) MarkPostUnliked(success bool) {
	p.LikeAttempted = true
	p.LikeSuccess = &success
	if success {
		p.Favorited = false
	}
}

// ---------------------------------------------------------------------------
// Text download methods
// ---------------------------------------------------------------------------

// MarkTextDownloaded records the outcome of a text file download.
//
// Parameters:
//   - success: Whether the text download succeeded.
func (p *Post) MarkTextDownloaded(success bool) {
	p.TextDownloadAttempted = true
	p.TextDownloadSucceeded = &success
}
