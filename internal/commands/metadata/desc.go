// =============================================================================
// FILE: internal/commands/metadata/desc.go
// PURPOSE: Description/text metadata extraction for media items. Extracts
//          and normalizes text metadata from media and their parent posts.
//          Ports Python metadata/desc.py.
// =============================================================================

package metadata

import (
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// MetadataDesc holds extracted description metadata for a media item.
// ---------------------------------------------------------------------------

// MetadataDesc contains normalized description fields for a media item.
type MetadataDesc struct {
	// PostText is the cleaned text from the parent post.
	PostText string

	// MediaType is the normalized media type string.
	MediaType string

	// Filename is the extracted filename from the media URL.
	Filename string

	// ResponseType is the content area this media belongs to.
	ResponseType string

	// Label is the label associated with the parent post.
	Label string

	// Value indicates "free" or "paid" status.
	Value string

	// Date is the formatted date string.
	Date string
}

// ---------------------------------------------------------------------------
// Extraction
// ---------------------------------------------------------------------------

// ExtractDesc extracts description metadata from a media item and its
// parent post.
//
// Parameters:
//   - media: The media item to extract metadata from.
//
// Returns:
//   - A MetadataDesc with all available fields populated.
func ExtractDesc(media *model.Media) MetadataDesc {
	desc := MetadataDesc{
		MediaType:    string(media.MediaType()),
		Filename:     media.Filename(),
		ResponseType: media.ResponseType,
		Label:        media.Label,
		Value:        media.Value,
		Date:         media.FormattedDate(),
	}

	// Extract post text if parent post is available.
	if media.Post != nil {
		desc.PostText = normalizeText(media.Post.RawText)
	} else if media.Text != "" {
		desc.PostText = normalizeText(media.Text)
	}

	return desc
}

// HasText returns true if the description has non-empty post text.
func (d MetadataDesc) HasText() bool {
	return d.PostText != ""
}

// Summary returns a brief one-line summary of the metadata.
func (d MetadataDesc) Summary() string {
	parts := []string{d.ResponseType, d.MediaType}
	if d.Filename != "" {
		parts = append(parts, d.Filename)
	}
	if d.Date != "" {
		parts = append(parts, d.Date)
	}
	return strings.Join(parts, " | ")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// normalizeText cleans text for metadata storage by trimming whitespace
// and collapsing runs of whitespace.
func normalizeText(text string) string {
	text = strings.TrimSpace(text)
	// Collapse multiple whitespace characters.
	fields := strings.Fields(text)
	return strings.Join(fields, " ")
}
