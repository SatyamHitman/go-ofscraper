// =============================================================================
// FILE: internal/filter/filter.go
// PURPOSE: Filter interfaces and chain builder. Defines the typed filter
//          function signatures for media, posts, and models. Provides a
//          chain builder that composes multiple filters into a pipeline.
//          Ports Python filters/ base module.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Filter function types
// ---------------------------------------------------------------------------

// MediaFilter transforms a media slice, typically removing items that
// don't match criteria. Returns a new slice (non-destructive).
type MediaFilter func([]*model.Media) []*model.Media

// PostFilter transforms a post slice.
type PostFilter func([]model.Post) []model.Post

// ModelFilter transforms a user/model slice.
type ModelFilter func([]model.User) []model.User

// ---------------------------------------------------------------------------
// Chain builders
// ---------------------------------------------------------------------------

// ChainMedia composes multiple media filters into a single pipeline.
// Filters are applied in order. An empty chain is a no-op.
//
// Parameters:
//   - filters: The filters to chain.
//
// Returns:
//   - A single MediaFilter that applies all filters in sequence.
func ChainMedia(filters ...MediaFilter) MediaFilter {
	return func(media []*model.Media) []*model.Media {
		for _, f := range filters {
			if f == nil {
				continue
			}
			media = f(media)
		}
		return media
	}
}

// ChainPosts composes multiple post filters into a single pipeline.
//
// Parameters:
//   - filters: The filters to chain.
//
// Returns:
//   - A single PostFilter that applies all filters in sequence.
func ChainPosts(filters ...PostFilter) PostFilter {
	return func(posts []model.Post) []model.Post {
		for _, f := range filters {
			if f == nil {
				continue
			}
			posts = f(posts)
		}
		return posts
	}
}

// ChainModels composes multiple model filters into a single pipeline.
//
// Parameters:
//   - filters: The filters to chain.
//
// Returns:
//   - A single ModelFilter that applies all filters in sequence.
func ChainModels(filters ...ModelFilter) ModelFilter {
	return func(users []model.User) []model.User {
		for _, f := range filters {
			if f == nil {
				continue
			}
			users = f(users)
		}
		return users
	}
}

// ---------------------------------------------------------------------------
// Helper: generic filter
// ---------------------------------------------------------------------------

// Where returns items where the predicate returns true.
func Where[T any](items []T, pred func(T) bool) []T {
	var result []T
	for _, item := range items {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}
