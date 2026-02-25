// =============================================================================
// FILE: internal/drm/dash.go
// PURPOSE: DASH MPD manifest parser. Parses MPEG-DASH manifest XML to extract
//          stream URLs, content protection info, and key IDs.
//          Ports Python mpegdash usage + custom parsing.
// =============================================================================

package drm

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ---------------------------------------------------------------------------
// MPD XML structures
// ---------------------------------------------------------------------------

// MPD represents the top-level DASH manifest.
type MPD struct {
	XMLName xml.Name `xml:"MPD"`
	Periods []Period `xml:"Period"`
}

// Period represents a time period in the manifest.
type Period struct {
	AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

// AdaptationSet represents a group of related representations.
type AdaptationSet struct {
	MimeType          string             `xml:"mimeType,attr"`
	ContentType       string             `xml:"contentType,attr"`
	ContentProtection []ContentProtection `xml:"ContentProtection"`
	Representations   []Representation   `xml:"Representation"`
}

// ContentProtection holds DRM protection information.
type ContentProtection struct {
	SchemeIDURI string `xml:"schemeIdUri,attr"`
	DefaultKID  string `xml:"default_KID,attr"`
	CencDefaultKID string `xml:"cenc:default_KID,attr"`
	PSSH        string `xml:"pssh"`
}

// Representation represents a single stream variant (quality level).
type Representation struct {
	ID        string      `xml:"id,attr"`
	Bandwidth int         `xml:"bandwidth,attr"`
	Width     int         `xml:"width,attr"`
	Height    int         `xml:"height,attr"`
	Codecs    string      `xml:"codecs,attr"`
	BaseURL   string      `xml:"BaseURL"`
	SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

// SegmentTemplate defines segment URL templates.
type SegmentTemplate struct {
	Media          string `xml:"media,attr"`
	Initialization string `xml:"initialization,attr"`
	Timescale      int    `xml:"timescale,attr"`
	Duration       int    `xml:"duration,attr"`
}

// ---------------------------------------------------------------------------
// Manifest wrapper
// ---------------------------------------------------------------------------

// Manifest wraps a parsed MPD with convenience methods.
type Manifest struct {
	MPD     *MPD
	BaseURL string // Base URL for resolving relative URLs
}

// ParseMPD fetches and parses an MPD manifest from a URL.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - mpdURL: The manifest URL to fetch.
//
// Returns:
//   - Parsed Manifest, or error.
func ParseMPD(ctx context.Context, mpdURL string) (*Manifest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mpdURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch MPD: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MPD fetch status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read MPD body: %w", err)
	}

	return ParseMPDBytes(body, mpdURL)
}

// ParseMPDBytes parses MPD manifest from raw XML bytes.
//
// Parameters:
//   - data: Raw XML bytes.
//   - baseURL: Base URL for resolving relative paths.
//
// Returns:
//   - Parsed Manifest, or error.
func ParseMPDBytes(data []byte, baseURL string) (*Manifest, error) {
	var mpd MPD
	if err := xml.Unmarshal(data, &mpd); err != nil {
		return nil, fmt.Errorf("parse MPD XML: %w", err)
	}

	// Derive base URL from MPD URL.
	base := baseURL
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[:idx+1]
	}

	return &Manifest{
		MPD:     &mpd,
		BaseURL: base,
	}, nil
}

// KeyID extracts the content encryption key ID from the manifest.
// Returns empty string if no key ID is found.
func (m *Manifest) KeyID() string {
	for _, period := range m.MPD.Periods {
		for _, as := range period.AdaptationSets {
			for _, cp := range as.ContentProtection {
				kid := cp.DefaultKID
				if kid == "" {
					kid = cp.CencDefaultKID
				}
				if kid != "" {
					// Normalize: remove dashes, lowercase.
					kid = strings.ReplaceAll(kid, "-", "")
					return strings.ToLower(kid)
				}
			}
		}
	}
	return ""
}

// PSSH extracts the Protection System Specific Header data.
// Returns empty string if not found.
func (m *Manifest) PSSH() string {
	for _, period := range m.MPD.Periods {
		for _, as := range period.AdaptationSets {
			for _, cp := range as.ContentProtection {
				if cp.PSSH != "" {
					return strings.TrimSpace(cp.PSSH)
				}
			}
		}
	}
	return ""
}

// VideoRepresentations returns all video representations sorted by bandwidth.
func (m *Manifest) VideoRepresentations() []Representation {
	var reps []Representation
	for _, period := range m.MPD.Periods {
		for _, as := range period.AdaptationSets {
			if strings.Contains(as.MimeType, "video") || as.ContentType == "video" {
				reps = append(reps, as.Representations...)
			}
		}
	}
	return reps
}

// AudioRepresentations returns all audio representations.
func (m *Manifest) AudioRepresentations() []Representation {
	var reps []Representation
	for _, period := range m.MPD.Periods {
		for _, as := range period.AdaptationSets {
			if strings.Contains(as.MimeType, "audio") || as.ContentType == "audio" {
				reps = append(reps, as.Representations...)
			}
		}
	}
	return reps
}
