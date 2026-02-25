// =============================================================================
// FILE: internal/commands/utils/scrape_context.go
// PURPOSE: ScrapeContext holds shared state during a scrape session. Aggregates
//          user lists, post collections, and result counters used across the
//          scrape pipeline stages. Ports Python data/models/selector.py state.
// =============================================================================

package cmdutils

import (
	"sync"
	"sync/atomic"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// ScrapeContext
// ---------------------------------------------------------------------------

// ScrapeContext holds shared mutable state for a single scrape session.
// All counters are safe for concurrent access.
type ScrapeContext struct {
	// Users is the resolved list of users to process.
	Users []*model.User

	// Posts accumulates posts found across all users.
	Posts []*model.Post
	postsMu sync.Mutex

	// Counters track progress for summary reporting.
	UsersProcessed atomic.Int64
	PostsFound     atomic.Int64
	MediaFound     atomic.Int64
	MediaDownloaded atomic.Int64
	MediaSkipped   atomic.Int64
	MediaFailed    atomic.Int64
	LikesAttempted atomic.Int64
	LikesSucceeded atomic.Int64
	Errors         atomic.Int64
}

// NewScrapeContext creates a new empty ScrapeContext.
//
// Returns:
//   - An initialized ScrapeContext ready for use.
func NewScrapeContext() *ScrapeContext {
	return &ScrapeContext{}
}

// AddPosts appends posts to the context's collection in a thread-safe manner.
//
// Parameters:
//   - posts: The posts to add.
func (sc *ScrapeContext) AddPosts(posts []*model.Post) {
	sc.postsMu.Lock()
	defer sc.postsMu.Unlock()
	sc.Posts = append(sc.Posts, posts...)
	sc.PostsFound.Add(int64(len(posts)))
}

// AllPosts returns a snapshot of all collected posts.
//
// Returns:
//   - A slice of all posts collected so far.
func (sc *ScrapeContext) AllPosts() []*model.Post {
	sc.postsMu.Lock()
	defer sc.postsMu.Unlock()
	result := make([]*model.Post, len(sc.Posts))
	copy(result, sc.Posts)
	return result
}

// RecordMediaResult increments the appropriate counter based on download status.
//
// Parameters:
//   - status: The download status to record.
func (sc *ScrapeContext) RecordMediaResult(status model.DownloadStatus) {
	switch status {
	case model.DownloadStatusSucceeded:
		sc.MediaDownloaded.Add(1)
	case model.DownloadStatusSkipped:
		sc.MediaSkipped.Add(1)
	case model.DownloadStatusFailed:
		sc.MediaFailed.Add(1)
		sc.Errors.Add(1)
	}
}
