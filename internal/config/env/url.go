// =============================================================================
// FILE: internal/config/env/url.go
// PURPOSE: Defines all OnlyFans API endpoint URL templates. Each endpoint is
//          loaded from an environment variable with a hardcoded default.
//          Ports Python utils/of_env/values/url/url.py.
// =============================================================================

package env

// ---------------------------------------------------------------------------
// API Endpoint URL templates
// ---------------------------------------------------------------------------

// InitEP returns the session initialization endpoint.
func InitEP() string {
	return GetString("OF_INIT_EP", "/api2/v2/init")
}

// MeEP returns the current user profile endpoint.
func MeEP() string {
	return GetString("OF_ME_EP", "/api2/v2/users/me")
}

// ProfileEP returns the user profile endpoint template.
// Format placeholder: user_id or username.
func ProfileEP() string {
	return GetString("OF_PROFILE_EP", "/api2/v2/users/%s")
}

// TimelineEP returns the timeline posts endpoint template.
// Format placeholder: user_id.
func TimelineEP() string {
	return GetString("OF_TIMELINE_EP",
		"/api2/v2/users/%d/posts?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&pinned=0&format=infinite")
}

// TimelineNextEP returns the paginated timeline posts endpoint template.
// Format placeholders: user_id, afterPublishTime timestamp.
func TimelineNextEP() string {
	return GetString("OF_TIMELINE_NEXT_EP",
		"/api2/v2/users/%d/posts?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&afterPublishTime=%s&pinned=0&format=infinite")
}

// TimelinePinnedEP returns the pinned posts endpoint template.
// Format placeholders: user_id, counters.
func TimelinePinnedEP() string {
	return GetString("OF_TIMELINE_PINNED_EP",
		"/api2/v2/users/%d/posts?skip_users=all&pinned=1&counters=%s&format=infinite")
}

// ArchivedEP returns the archived posts endpoint template.
// Format placeholder: user_id.
func ArchivedEP() string {
	return GetString("OF_ARCHIVED_EP",
		"/api2/v2/users/%d/posts/archived?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&format=infinite")
}

// ArchivedNextEP returns the paginated archived posts endpoint template.
// Format placeholders: user_id, afterPublishTime timestamp.
func ArchivedNextEP() string {
	return GetString("OF_ARCHIVED_NEXT_EP",
		"/api2/v2/users/%d/posts/archived?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&afterPublishTime=%s&format=infinite")
}

// StreamsEP returns the streams endpoint template.
// Format placeholder: user_id.
func StreamsEP() string {
	return GetString("OF_STREAMS_EP",
		"/api2/v2/users/%d/posts/streams?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&format=infinite")
}

// StreamsNextEP returns the paginated streams endpoint template.
// Format placeholders: user_id, afterPublishTime timestamp.
func StreamsNextEP() string {
	return GetString("OF_STREAMS_NEXT_EP",
		"/api2/v2/users/%d/posts/streams?limit=100&order=publish_date_asc&skip_users=all&skip_users_dups=1&afterPublishTime=%s&format=infinite")
}

// MessagesEP returns the messages endpoint template.
// Format placeholder: chat_id (which is the model's user_id).
func MessagesEP() string {
	return GetString("OF_MESSAGES_EP",
		"/api2/v2/chats/%d/messages?limit=100&order=desc&skip_users=all&skip_users_dups=1")
}

// MessagesNextEP returns the paginated messages endpoint template.
// Format placeholders: chat_id, message_id for pagination cursor.
func MessagesNextEP() string {
	return GetString("OF_MESSAGES_NEXT_EP",
		"/api2/v2/chats/%d/messages?limit=100&id=%d&order=desc&skip_users=all&skip_users_dups=1")
}

// SubscriptionsEP returns the all-subscriptions endpoint template.
// Format placeholder: offset.
func SubscriptionsEP() string {
	return GetString("OF_SUBSCRIPTIONS_EP",
		"/api2/v2/subscriptions/subscribes?offset=%d&limit=10&type=all&format=infinite")
}

// SubscriptionsActiveEP returns the active subscriptions endpoint template.
// Format placeholder: offset.
func SubscriptionsActiveEP() string {
	return GetString("OF_SUBSCRIPTIONS_ACTIVE_EP",
		"/api2/v2/subscriptions/subscribes?offset=%d&limit=10&type=active&format=infinite")
}

// SubscriptionsExpiredEP returns the expired subscriptions endpoint template.
// Format placeholder: offset.
func SubscriptionsExpiredEP() string {
	return GetString("OF_SUBSCRIPTIONS_EXPIRED_EP",
		"/api2/v2/subscriptions/subscribes?offset=%d&limit=10&type=expired&format=infinite")
}

// LabelsEP returns the labels list endpoint template.
// Format placeholders: user_id, offset.
func LabelsEP() string {
	return GetString("OF_LABELS_EP",
		"/api2/v2/users/%d/labels?limit=100&offset=%d&order=desc&non-empty=1")
}

// LabelledPostsEP returns the labelled posts endpoint template.
// Format placeholders: user_id, offset, label_id.
func LabelledPostsEP() string {
	return GetString("OF_LABELLED_POSTS_EP",
		"/api2/v2/users/%d/posts?limit=100&offset=%d&order=publish_date_desc&skip_users=all&counters=0&format=infinite&label=%d")
}

// HighlightsWithStoriesEP returns the highlights-with-stories endpoint template.
// Format placeholders: user_id, offset.
func HighlightsWithStoriesEP() string {
	return GetString("OF_HIGHLIGHTS_STORIES_EP",
		"/api2/v2/users/%d/stories/highlights?limit=5&offset=%d&unf=1")
}

// HighlightsWithAStoryEP returns the single-story highlights endpoint.
// Format placeholder: user_id.
func HighlightsWithAStoryEP() string {
	return GetString("OF_HIGHLIGHTS_A_STORY_EP",
		"/api2/v2/users/%d/stories?unf=1")
}

// StoryEP returns the individual story endpoint template.
// Format placeholder: story_id.
func StoryEP() string {
	return GetString("OF_STORY_EP",
		"/api2/v2/stories/highlights/%d?unf=1")
}

// PurchasedContentEP returns the purchased content by user endpoint.
// Format placeholders: offset, author_id.
func PurchasedContentEP() string {
	return GetString("OF_PURCHASED_EP",
		"/api2/v2/posts/paid/all?limit=100&skip_users=all&format=infinite&offset=%d&author=%d")
}

// PurchasedContentAllEP returns the all purchased content endpoint.
// Format placeholder: offset.
func PurchasedContentAllEP() string {
	return GetString("OF_PURCHASED_ALL_EP",
		"/api2/v2/posts/paid/all?limit=100&skip_users=all&format=infinite&offset=%d")
}

// LicenceURL returns the DRM license endpoint template.
// Format placeholders: media_id, drm_id, drm_type.
func LicenceURL() string {
	return GetString("OF_LICENCE_URL",
		"/api2/v2/users/media/%d/drm/%s/%s?type=widevine")
}

// IndividualTimelineEP returns the single post endpoint template.
// Format placeholder: post_id.
func IndividualTimelineEP() string {
	return GetString("OF_INDIVIDUAL_TIMELINE_EP",
		"/api2/v2/posts/%d?skip_users=all")
}

// FavoriteEP returns the favorite/like endpoint template.
// Format placeholders: post_id, user_id.
func FavoriteEP() string {
	return GetString("OF_FAVORITE_EP",
		"/api2/v2/posts/%d/favorites/%d")
}

// ListEP returns the custom lists endpoint template.
// Format placeholder: offset.
func ListEP() string {
	return GetString("OF_LIST_EP",
		"/api2/v2/lists?offset=%d&skip_users=all&limit=100&format=infinite")
}

// ListUsersEP returns the list users endpoint template.
// Format placeholders: list_id, offset.
func ListUsersEP() string {
	return GetString("OF_LIST_USERS_EP",
		"/api2/v2/lists/%d/users?offset=%d&limit=100&format=infinite")
}

// BaseURL returns the OnlyFans base URL for constructing full API URLs.
func BaseURL() string {
	return GetString("OF_BASE_URL", "https://onlyfans.com")
}
