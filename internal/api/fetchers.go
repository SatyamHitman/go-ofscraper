// =============================================================================
// FILE: internal/api/fetchers.go
// PURPOSE: Remaining ContentFetcher interface implementations. Provides
//          GetStories, GetHighlights, GetPinned, GetArchived, GetStreams,
//          GetLabels, GetPurchased, GetSubscriptions, GetProfile, and
//          PostFavorite. Ports corresponding Python data/api/ modules.
// =============================================================================

package api

import (
	"context"
	"fmt"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Stories
// ---------------------------------------------------------------------------

// GetStories fetches story posts for a model.
func (c *Client) GetStories(ctx context.Context, modelID int64) ([]model.Post, error) {
	url := HighlightsURL(modelID)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetStories: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetStories: status %d", resp.StatusCode)
	}

	var raw []any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetStories: decode error: %w", err)
	}

	var posts []model.Post
	for _, item := range raw {
		if storyMap, ok := item.(map[string]any); ok {
			post := parsePost(storyMap, "stories", modelID)
			posts = append(posts, post)
		}
	}
	return posts, nil
}

// ---------------------------------------------------------------------------
// Highlights
// ---------------------------------------------------------------------------

// GetHighlights fetches highlight posts for a model.
func (c *Client) GetHighlights(ctx context.Context, modelID int64) ([]model.Post, error) {
	url := HighlightsURL(modelID)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetHighlights: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetHighlights: status %d", resp.StatusCode)
	}

	var raw []any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetHighlights: decode error: %w", err)
	}

	var posts []model.Post
	for _, item := range raw {
		if hlMap, ok := item.(map[string]any); ok {
			// Highlights contain nested stories.
			if stories, ok := hlMap["stories"].([]any); ok {
				for _, s := range stories {
					if storyMap, ok := s.(map[string]any); ok {
						post := parsePost(storyMap, "highlights", modelID)
						posts = append(posts, post)
					}
				}
			}
		}
	}
	return posts, nil
}

// ---------------------------------------------------------------------------
// Pinned
// ---------------------------------------------------------------------------

// GetPinned fetches pinned posts for a model.
func (c *Client) GetPinned(ctx context.Context, modelID int64) ([]model.Post, error) {
	url := PinnedURL(modelID)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetPinned: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetPinned: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetPinned: decode error: %w", err)
	}

	return parsePostList(raw, "pinned", modelID), nil
}

// ---------------------------------------------------------------------------
// Archived
// ---------------------------------------------------------------------------

// GetArchived fetches archived posts for a model.
func (c *Client) GetArchived(ctx context.Context, modelID int64, after float64) ([]model.Post, error) {
	var url string
	if after > 0 {
		url = ArchivedNextURL(modelID, after)
	} else {
		url = ArchivedURL(modelID)
	}

	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetArchived: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetArchived: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetArchived: decode error: %w", err)
	}

	return parsePostList(raw, "archived", modelID), nil
}

// ---------------------------------------------------------------------------
// Streams
// ---------------------------------------------------------------------------

// GetStreams fetches stream posts for a model.
func (c *Client) GetStreams(ctx context.Context, modelID int64, after float64) ([]model.Post, error) {
	var url string
	if after > 0 {
		url = StreamsNextURL(modelID, after)
	} else {
		url = StreamsURL(modelID)
	}

	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetStreams: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetStreams: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetStreams: decode error: %w", err)
	}

	return parsePostList(raw, "streams", modelID), nil
}

// ---------------------------------------------------------------------------
// Labels
// ---------------------------------------------------------------------------

// GetLabels fetches labelled posts for a model.
func (c *Client) GetLabels(ctx context.Context, modelID int64) ([]model.Post, error) {
	url := LabelsURL(modelID)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetLabels: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetLabels: status %d", resp.StatusCode)
	}

	var raw []any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetLabels: decode error: %w", err)
	}

	var posts []model.Post
	for _, item := range raw {
		if labelMap, ok := item.(map[string]any); ok {
			post := parsePost(labelMap, "labels", modelID)
			posts = append(posts, post)
		}
	}
	return posts, nil
}

// ---------------------------------------------------------------------------
// Purchased
// ---------------------------------------------------------------------------

// GetPurchased fetches purchased content for a model.
func (c *Client) GetPurchased(ctx context.Context, modelID int64) ([]model.Post, error) {
	url := PurchasedURL(modelID)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetPurchased: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetPurchased: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetPurchased: decode error: %w", err)
	}

	return parsePostList(raw, "purchased", modelID), nil
}

// ---------------------------------------------------------------------------
// Subscriptions
// ---------------------------------------------------------------------------

// GetSubscriptions fetches the user's active or expired subscriptions.
func (c *Client) GetSubscriptions(ctx context.Context, subType string) ([]model.User, error) {
	url := SubscriptionsURL(subType)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return nil, fmt.Errorf("GetSubscriptions: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return nil, fmt.Errorf("GetSubscriptions: status %d", resp.StatusCode)
	}

	var raw []any
	if err := resp.JSON(&raw); err != nil {
		return nil, fmt.Errorf("GetSubscriptions: decode error: %w", err)
	}

	var users []model.User
	for _, item := range raw {
		if userMap, ok := item.(map[string]any); ok {
			users = append(users, parseUserFromMap(userMap))
		}
	}
	return users, nil
}

// ---------------------------------------------------------------------------
// Profile
// ---------------------------------------------------------------------------

// GetProfile fetches a model's public profile.
func (c *Client) GetProfile(ctx context.Context, username string) (model.User, error) {
	url := ProfileURL(username)
	req := gohttp.NewRequest(url)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return model.User{}, fmt.Errorf("GetProfile: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return model.User{}, fmt.Errorf("GetProfile: status %d", resp.StatusCode)
	}

	var raw map[string]any
	if err := resp.JSON(&raw); err != nil {
		return model.User{}, fmt.Errorf("GetProfile: decode error: %w", err)
	}

	return parseUserFromMap(raw), nil
}

// ---------------------------------------------------------------------------
// Favorite (like/unlike)
// ---------------------------------------------------------------------------

// PostFavorite likes or unlikes a post.
func (c *Client) PostFavorite(ctx context.Context, postID int64, action string) error {
	url := FavoriteURL(postID, action)
	req := gohttp.NewPostRequest(url, nil)
	resp, err := gohttp.DoWithRetry(ctx, c.session, req, gohttp.DefaultRetryConfig())
	if err != nil {
		return fmt.Errorf("PostFavorite: %w", err)
	}
	if !resp.IsOK() {
		resp.Close()
		return fmt.Errorf("PostFavorite: status %d", resp.StatusCode)
	}
	resp.Close()
	return nil
}
