// =============================================================================
// FILE: internal/api/client.go
// PURPOSE: API client struct and ContentFetcher interface. The API client
//          wraps the HTTP session manager with OF-specific endpoint logic for
//          fetching timelines, messages, stories, subscriptions, etc. Ports
//          Python data/api/ module.
// =============================================================================

package api

import (
	"context"
	"log/slog"

	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// ContentFetcher interface
// ---------------------------------------------------------------------------

// ContentFetcher defines the contract for fetching content from the OF API.
type ContentFetcher interface {
	GetTimeline(ctx context.Context, modelID int64, after float64) ([]model.Post, error)
	GetMessages(ctx context.Context, modelID int64, after float64) ([]model.Post, error)
	GetStories(ctx context.Context, modelID int64) ([]model.Post, error)
	GetHighlights(ctx context.Context, modelID int64) ([]model.Post, error)
	GetPinned(ctx context.Context, modelID int64) ([]model.Post, error)
	GetArchived(ctx context.Context, modelID int64, after float64) ([]model.Post, error)
	GetStreams(ctx context.Context, modelID int64, after float64) ([]model.Post, error)
	GetLabels(ctx context.Context, modelID int64) ([]model.Post, error)
	GetPurchased(ctx context.Context, modelID int64) ([]model.Post, error)
	GetSubscriptions(ctx context.Context, subType string) ([]model.User, error)
	GetProfile(ctx context.Context, username string) (model.User, error)
	GetMe(ctx context.Context) (model.User, error)
	PostFavorite(ctx context.Context, postID int64, action string) error
}

// ---------------------------------------------------------------------------
// Client
// ---------------------------------------------------------------------------

// Client implements ContentFetcher using the HTTP SessionManager.
type Client struct {
	session *gohttp.SessionManager
	log     *slog.Logger
}

// NewClient creates an API client backed by the given session manager.
//
// Parameters:
//   - session: The HTTP session manager with auth configured.
//
// Returns:
//   - A new Client implementing ContentFetcher.
func NewClient(session *gohttp.SessionManager) *Client {
	return &Client{
		session: session,
		log:     slog.Default().With("component", "api"),
	}
}

// ---------------------------------------------------------------------------
// Ensure Client implements ContentFetcher
// ---------------------------------------------------------------------------

var _ ContentFetcher = (*Client)(nil)
