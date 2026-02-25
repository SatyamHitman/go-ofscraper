// =============================================================================
// FILE: internal/app/post_collection.go
// PURPOSE: PostCollection aggregates posts from multiple content areas
//          (timeline, messages, stories, etc.) into a unified collection for
//          processing. Provides filtering and grouping utilities.
//          Ports Python data/post.py post collection logic.
// =============================================================================

package app

import (
	"sync"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// PostCollection
// ---------------------------------------------------------------------------

// PostCollection aggregates posts from multiple content areas into a single
// collection with filtering and grouping capabilities.
type PostCollection struct {
	mu    sync.RWMutex
	posts []*model.Post

	// byArea groups posts by their response type area.
	byArea map[model.ResponseType][]*model.Post

	// byUser groups posts by username.
	byUser map[string][]*model.Post
}

// NewPostCollection creates an empty PostCollection.
//
// Returns:
//   - An initialized PostCollection.
func NewPostCollection() *PostCollection {
	return &PostCollection{
		byArea: make(map[model.ResponseType][]*model.Post),
		byUser: make(map[string][]*model.Post),
	}
}

// Add appends posts to the collection and updates indexes.
//
// Parameters:
//   - posts: The posts to add.
func (pc *PostCollection) Add(posts ...*model.Post) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	for _, post := range posts {
		pc.posts = append(pc.posts, post)
		pc.byArea[post.ResponseType] = append(pc.byArea[post.ResponseType], post)
		pc.byUser[post.Username] = append(pc.byUser[post.Username], post)
	}
}

// All returns all posts in the collection.
//
// Returns:
//   - A snapshot slice of all posts.
func (pc *PostCollection) All() []*model.Post {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	result := make([]*model.Post, len(pc.posts))
	copy(result, pc.posts)
	return result
}

// ByArea returns posts filtered by response type.
//
// Parameters:
//   - area: The response type to filter by.
//
// Returns:
//   - Posts matching the given area.
func (pc *PostCollection) ByArea(area model.ResponseType) []*model.Post {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	posts := pc.byArea[area]
	result := make([]*model.Post, len(posts))
	copy(result, posts)
	return result
}

// ByUser returns posts filtered by username.
//
// Parameters:
//   - username: The username to filter by.
//
// Returns:
//   - Posts from the given user.
func (pc *PostCollection) ByUser(username string) []*model.Post {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	posts := pc.byUser[username]
	result := make([]*model.Post, len(posts))
	copy(result, posts)
	return result
}

// Len returns the total number of posts in the collection.
//
// Returns:
//   - The post count.
func (pc *PostCollection) Len() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return len(pc.posts)
}

// Areas returns the distinct content areas present in the collection.
//
// Returns:
//   - A slice of ResponseType values.
func (pc *PostCollection) Areas() []model.ResponseType {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	areas := make([]model.ResponseType, 0, len(pc.byArea))
	for area := range pc.byArea {
		areas = append(areas, area)
	}
	return areas
}

// Users returns the distinct usernames present in the collection.
//
// Returns:
//   - A slice of username strings.
func (pc *PostCollection) Users() []string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	users := make([]string, 0, len(pc.byUser))
	for user := range pc.byUser {
		users = append(users, user)
	}
	return users
}

// AllMedia extracts all viewable media from all posts in the collection.
//
// Returns:
//   - A flat slice of all viewable Media items.
func (pc *PostCollection) AllMedia() []*model.Media {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	var media []*model.Media
	for _, post := range pc.posts {
		media = append(media, post.ViewableMedia()...)
	}
	return media
}
