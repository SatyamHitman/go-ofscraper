// =============================================================================
// FILE: internal/api/endpoints.go
// PURPOSE: All OF API endpoint URL templates. Provides functions that build
//          full API URLs with model IDs, pagination cursors, and query params.
//          Ports Python of_env/values/url/url.py.
// =============================================================================

package api

import (
	"fmt"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Endpoint builders
// ---------------------------------------------------------------------------

// base returns the OF API base URL.
func base() string {
	return env.BaseURL()
}

// MeURL returns the /me endpoint URL.
func MeURL() string {
	return base() + env.MeEP()
}

// ProfileURL returns the user profile endpoint.
//
// Parameters:
//   - username: The model's username.
func ProfileURL(username string) string {
	return base() + fmt.Sprintf(env.ProfileEP(), username)
}

// TimelineURL returns the timeline posts endpoint.
//
// Parameters:
//   - modelID: The model's numeric ID.
func TimelineURL(modelID int64) string {
	return base() + fmt.Sprintf(env.TimelineEP(), modelID)
}

// TimelineNextURL returns the paginated timeline endpoint.
//
// Parameters:
//   - modelID: The model's numeric ID.
//   - after: The pagination cursor timestamp.
func TimelineNextURL(modelID int64, after float64) string {
	return base() + fmt.Sprintf(env.TimelineNextEP(), modelID, after)
}

// PinnedURL returns the pinned posts endpoint.
//
// Parameters:
//   - modelID: The model's numeric ID.
func PinnedURL(modelID int64) string {
	return base() + fmt.Sprintf(env.TimelinePinnedEP(), modelID)
}

// ArchivedURL returns the archived posts endpoint.
//
// Parameters:
//   - modelID: The model's numeric ID.
func ArchivedURL(modelID int64) string {
	return base() + fmt.Sprintf(env.ArchivedEP(), modelID)
}

// ArchivedNextURL returns the paginated archived endpoint.
func ArchivedNextURL(modelID int64, after float64) string {
	return base() + fmt.Sprintf(env.ArchivedNextEP(), modelID, after)
}

// StreamsURL returns the streams endpoint.
func StreamsURL(modelID int64) string {
	return base() + fmt.Sprintf(env.StreamsEP(), modelID)
}

// StreamsNextURL returns the paginated streams endpoint.
func StreamsNextURL(modelID int64, after float64) string {
	return base() + fmt.Sprintf(env.StreamsNextEP(), modelID, after)
}

// MessagesURL returns the messages endpoint.
func MessagesURL(modelID int64) string {
	return base() + fmt.Sprintf(env.MessagesEP(), modelID)
}

// MessagesNextURL returns the paginated messages endpoint.
func MessagesNextURL(modelID int64, after string) string {
	return base() + fmt.Sprintf(env.MessagesNextEP(), modelID, after)
}

// HighlightsURL returns the highlights + stories endpoint.
func HighlightsURL(modelID int64) string {
	return base() + fmt.Sprintf(env.HighlightsWithStoriesEP(), modelID)
}

// StoryURL returns a specific story endpoint.
func StoryURL(storyID int64) string {
	return base() + fmt.Sprintf(env.StoryEP(), storyID)
}

// LabelsURL returns the labels list endpoint.
func LabelsURL(modelID int64) string {
	return base() + fmt.Sprintf(env.LabelsEP(), modelID)
}

// LabelledPostsURL returns the posts-by-label endpoint.
func LabelledPostsURL(modelID, labelID int64) string {
	return base() + fmt.Sprintf(env.LabelledPostsEP(), modelID, labelID)
}

// SubscriptionsURL returns the subscriptions list endpoint.
func SubscriptionsURL(subType string) string {
	return base() + fmt.Sprintf(env.SubscriptionsEP(), subType)
}

// PurchasedURL returns the purchased content endpoint.
func PurchasedURL(modelID int64) string {
	return base() + fmt.Sprintf(env.PurchasedContentEP(), modelID)
}

// FavoriteURL returns the like/unlike endpoint.
func FavoriteURL(postID int64, action string) string {
	return base() + fmt.Sprintf(env.FavoriteEP(), postID, action)
}

// InitURL returns the init/startup endpoint.
func InitURL() string {
	return base() + env.InitEP()
}
