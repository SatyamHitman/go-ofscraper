// =============================================================================
// FILE: internal/api/messages.go
// PURPOSE: Messages API handler. Fetches direct messages with ID-based
//          pagination. Ports Python data/api/messages.py.
// =============================================================================

package api

import (
	"context"
	"fmt"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// GetMessages fetches messages for a model with pagination.
//
// Parameters:
//   - ctx: Context.
//   - modelID: The model's numeric ID.
//   - after: Pagination cursor (0 for first page).
//
// Returns:
//   - Slice of posts (messages), and any error.
func (c *Client) GetMessages(ctx context.Context, modelID int64, after float64) ([]model.Post, error) {
	var url string
	if after > 0 {
		url = MessagesNextURL(modelID, fmt.Sprintf("%.0f", after))
	} else {
		url = MessagesURL(modelID)
	}

	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetMessages: %w", err)
	}

	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetMessages: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetMessages: decode error: %w", err)
	}

	return parsePostList(raw, "messages", modelID), nil
}
