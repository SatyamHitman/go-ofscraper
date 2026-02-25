// =============================================================================
// FILE: internal/app/stats.go
// PURPOSE: Statistics tracking and reporting for scrape operations. Tracks
//          users processed, posts found, media downloaded, errors, and timing.
//          Ports Python utils/logs/stats.py.
// =============================================================================

package app

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// Stats
// ---------------------------------------------------------------------------

// Stats tracks and reports statistics for a scrape run.
type Stats struct {
	// Counters
	UsersProcessed  atomic.Int64
	PostsFound      atomic.Int64
	MediaFound      atomic.Int64
	MediaDownloaded atomic.Int64
	MediaSkipped    atomic.Int64
	MediaFailed     atomic.Int64
	LikesAttempted  atomic.Int64
	LikesSucceeded  atomic.Int64
	MetadataChanged atomic.Int64
	Errors          atomic.Int64

	// Timing
	startTime time.Time
	endTime   time.Time
	mu        sync.Mutex
}

// NewStats creates a new Stats instance with the start time set to now.
//
// Returns:
//   - An initialized Stats.
func NewStats() *Stats {
	return &Stats{
		startTime: time.Now(),
	}
}

// MarkDone records the end time for the stats period.
func (s *Stats) MarkDone() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endTime = time.Now()
}

// Duration returns the elapsed time of the tracked period.
//
// Returns:
//   - The duration from start to end (or now if not yet done).
func (s *Stats) Duration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	end := s.endTime
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(s.startTime)
}

// Log writes the current statistics to the provided logger.
//
// Parameters:
//   - logger: The structured logger for output.
func (s *Stats) Log(logger *slog.Logger) {
	logger.Info("scrape statistics",
		"duration", s.Duration().Round(time.Second),
		"users_processed", s.UsersProcessed.Load(),
		"posts_found", s.PostsFound.Load(),
		"media_found", s.MediaFound.Load(),
		"media_downloaded", s.MediaDownloaded.Load(),
		"media_skipped", s.MediaSkipped.Load(),
		"media_failed", s.MediaFailed.Load(),
		"likes_attempted", s.LikesAttempted.Load(),
		"likes_succeeded", s.LikesSucceeded.Load(),
		"metadata_changed", s.MetadataChanged.Load(),
		"errors", s.Errors.Load(),
	)
}

// Summary returns a formatted multi-line summary string.
//
// Returns:
//   - A formatted summary of all statistics.
func (s *Stats) Summary() string {
	return fmt.Sprintf(
		"Duration: %s\n"+
			"Users Processed:  %d\n"+
			"Posts Found:      %d\n"+
			"Media Found:      %d\n"+
			"Media Downloaded: %d\n"+
			"Media Skipped:    %d\n"+
			"Media Failed:     %d\n"+
			"Likes Attempted:  %d\n"+
			"Likes Succeeded:  %d\n"+
			"Metadata Changed: %d\n"+
			"Errors:           %d",
		s.Duration().Round(time.Second),
		s.UsersProcessed.Load(),
		s.PostsFound.Load(),
		s.MediaFound.Load(),
		s.MediaDownloaded.Load(),
		s.MediaSkipped.Load(),
		s.MediaFailed.Load(),
		s.LikesAttempted.Load(),
		s.LikesSucceeded.Load(),
		s.MetadataChanged.Load(),
		s.Errors.Load(),
	)
}

// HasErrors returns true if any errors were recorded.
func (s *Stats) HasErrors() bool {
	return s.Errors.Load() > 0
}

// Reset clears all counters and resets the start time.
func (s *Stats) Reset() {
	s.UsersProcessed.Store(0)
	s.PostsFound.Store(0)
	s.MediaFound.Store(0)
	s.MediaDownloaded.Store(0)
	s.MediaSkipped.Store(0)
	s.MediaFailed.Store(0)
	s.LikesAttempted.Store(0)
	s.LikesSucceeded.Store(0)
	s.MetadataChanged.Store(0)
	s.Errors.Store(0)

	s.mu.Lock()
	s.startTime = time.Now()
	s.endTime = time.Time{}
	s.mu.Unlock()
}
