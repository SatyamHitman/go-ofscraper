// =============================================================================
// FILE: internal/api/timeline.go
// PURPOSE: Timeline posts API handler. Fetches timeline posts with pagination
//          using timestamp-based cursors. Ports Python data/api/timeline.py.
// =============================================================================

package api

import (
	"context"
	"fmt"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Timeline
// ---------------------------------------------------------------------------

// GetTimeline fetches timeline posts for a model with pagination.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - modelID: The model's numeric ID.
//   - after: Pagination cursor (0 for first page).
//
// Returns:
//   - Slice of posts, and any error.
func (c *Client) GetTimeline(ctx context.Context, modelID int64, after float64) ([]model.Post, error) {
	var url string
	if after > 0 {
		url = TimelineNextURL(modelID, after)
	} else {
		url = TimelineURL(modelID)
	}

	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetTimeline: %w", err)
	}

	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetTimeline: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetTimeline: decode error: %w", err)
	}

	return parsePostList(raw, "timeline", modelID), nil
}
