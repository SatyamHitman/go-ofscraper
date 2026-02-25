// =============================================================================
// FILE: internal/api/common.go
// PURPOSE: Common API parsing helpers. Shared functions for extracting posts,
//          media, and pagination data from API response JSON. Combines logic
//          from Python data/api/common/after.py, check.py, timeline.py.
// =============================================================================

package api

import (
	"fmt"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Response parsing
// ---------------------------------------------------------------------------

// parsePostList extracts a list of posts from an API response map.
// Handles both array responses and object responses with a "list" key.
//
// Parameters:
//   - raw: The raw JSON response as a map.
//   - apiType: The API source type (e.g. "timeline", "messages").
//   - modelID: The model's numeric ID.
//
// Returns:
//   - Slice of parsed posts.
func parsePostList(raw map[string]any, apiType string, modelID int64) []model.Post {
	var posts []model.Post

	// Try "list" key first (common response format).
	list, ok := raw["list"].([]any)
	if !ok {
		// Try "hasMore" paginated format.
		list, ok = raw["list"].([]any)
		if !ok {
			return posts
		}
	}

	for _, item := range list {
		postMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		post := parsePost(postMap, apiType, modelID)
		posts = append(posts, post)
	}

	return posts
}

// parsePost extracts a single Post from a JSON map.
func parsePost(raw map[string]any, apiType string, modelID int64) model.Post {
	p := model.Post{
		ModelID: modelID,
	}

	if id, ok := raw["id"].(float64); ok {
		p.ID = int64(id)
	}
	if text, ok := raw["text"].(string); ok {
		p.RawText = text
	}
	if text, ok := raw["rawText"].(string); ok {
		if text != "" {
			p.RawText = text
		}
	}
	if price, ok := raw["price"].(float64); ok {
		p.Price = price
	}
	if paid, ok := raw["isPaid"].(bool); ok {
		p.Paid = paid
	}
	if archived, ok := raw["isArchived"].(bool); ok {
		p.Archived = archived
	}
	if created, ok := raw["postedAt"].(string); ok {
		p.CreatedAt = created
	}
	if created, ok := raw["createdAt"].(string); ok {
		if p.CreatedAt == "" {
			p.CreatedAt = created
		}
	}

	// Parse media list.
	if mediaList, ok := raw["media"].([]any); ok {
		for _, mediaItem := range mediaList {
			if mediaMap, ok := mediaItem.(map[string]any); ok {
				media := parseMedia(mediaMap, p.ID)
				p.AllMedia = append(p.AllMedia, media)
			}
		}
	}

	return p
}

// parseMedia extracts a Media struct from a JSON map.
func parseMedia(raw map[string]any, postID int64) *model.Media {
	m := &model.Media{
		PostID: postID,
	}

	if id, ok := raw["id"].(float64); ok {
		m.ID = int64(id)
	}
	if t, ok := raw["type"].(string); ok {
		m.Type = t
	}
	if created, ok := raw["createdAt"].(string); ok {
		m.CreatedAt = created
	}

	// Extract source URL (full > source > preview).
	if src, ok := raw["full"].(string); ok && src != "" {
		m.RawURL = src
	} else if src, ok := raw["source"].(map[string]any); ok {
		if srcURL, ok := src["source"].(string); ok {
			m.RawURL = srcURL
		}
	}

	// Check if DRM protected.
	if files, ok := raw["files"].(map[string]any); ok {
		if drm, ok := files["drm"].(map[string]any); ok {
			if manifest, ok := drm["manifest"].(map[string]any); ok {
				if dash, ok := manifest["dash"].(string); ok && dash != "" {
					m.MpdURL = dash
				}
			}
		}
	}

	return m
}

// extractPaginationAfter reads the "tailMarker" or "afterPublishTime"
// cursor from a paginated response.
//
// Parameters:
//   - raw: The API response map.
//
// Returns:
//   - The after cursor value as float64, or 0 if not present.
func extractPaginationAfter(raw map[string]any) float64 {
	if tail, ok := raw["tailMarker"].(string); ok && tail != "" {
		// Try to parse as number.
		var val float64
		fmt.Sscanf(tail, "%f", &val)
		return val
	}
	if after, ok := raw["afterPublishTime"].(float64); ok {
		return after
	}
	return 0
}
