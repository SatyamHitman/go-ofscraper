// =============================================================================
// FILE: internal/api/me.go
// PURPOSE: /me endpoint handler. Fetches the authenticated user's profile
//          to verify auth and get the current user ID. Ports Python
//          data/api/me.py.
// =============================================================================

package api

import (
	"context"
	"fmt"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// /me endpoint
// ---------------------------------------------------------------------------

// GetMe fetches the authenticated user's profile.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - The current user as model.User, and any error.
func (c *Client) GetMe(ctx context.Context) (model.User, error) {
	req := gohttp.NewRequest(MeURL())
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return model.User{}, fmt.Errorf("GetMe failed: %w", err)
	}

	if !resp.IsOK() {
		resp.Close()
		return model.User{}, fmt.Errorf("GetMe: unexpected status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return model.User{}, fmt.Errorf("GetMe: failed to decode: %w", err)
	}

	user := parseUserFromMap(raw)
	c.log.Info("authenticated user loaded", "id", user.ID, "name", user.Name)

	return user, nil
}

// parseUserFromMap extracts a model.User from the /me response JSON.
func parseUserFromMap(raw map[string]any) model.User {
	u := model.User{}

	if id, ok := raw["id"].(float64); ok {
		u.ID = int64(id)
	}
	if name, ok := raw["name"].(string); ok {
		u.Name = name
	}

	return u
}
