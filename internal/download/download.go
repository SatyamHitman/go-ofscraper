// =============================================================================
// FILE: internal/download/download.go
// PURPOSE: Download orchestrator. Coordinates the entire download pipeline:
//          collects media items, dispatches to workers, tracks progress, and
//          reports results. Ports Python commands/scraper/actions/download/.
// =============================================================================

package download

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"gofscraper/internal/drm"
	gohttp "gofscraper/internal/http"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

// Config holds download configuration.
type Config struct {
	Workers       int     // Number of concurrent download workers
	ChunkSize     int64   // Download chunk size in bytes
	MaxRetries    int     // Max retries per download
	SpeedLimit    int64   // Bandwidth limit in bytes/sec (0 = unlimited)
	TempDir       string  // Temp directory for in-progress downloads
	SkipPrevious  bool    // Skip previously downloaded media
	ResumeEnabled bool    // Enable resume for interrupted downloads
	FFmpegPath    string  // Path to FFmpeg binary
	Logger        *slog.Logger
}

// DefaultConfig returns sensible download defaults.
func DefaultConfig() Config {
	return Config{
		Workers:       5,
		ChunkSize:     1024 * 1024, // 1 MB
		MaxRetries:    3,
		SpeedLimit:    0,
		TempDir:       "",
		SkipPrevious:  true,
		ResumeEnabled: true,
		FFmpegPath:    "ffmpeg",
	}
}

// ---------------------------------------------------------------------------
// Result
// ---------------------------------------------------------------------------

// Result tracks the outcome of the download batch.
type Result struct {
	mu        sync.Mutex
	Total     int
	Succeeded int
	Failed    int
	Skipped   int
	Errors    []DownloadError
}

// DownloadError records a single download failure.
type DownloadError struct {
	MediaID int64
	URL     string
	Err     error
}

// AddSuccess increments the success counter.
func (r *Result) AddSuccess() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Succeeded++
}

// AddFailure records a failure.
func (r *Result) AddFailure(mediaID int64, url string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Failed++
	r.Errors = append(r.Errors, DownloadError{
		MediaID: mediaID,
		URL:     url,
		Err:     err,
	})
}

// AddSkipped increments the skip counter.
func (r *Result) AddSkipped() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Skipped++
}

// Summary returns a human-readable summary.
func (r *Result) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fmt.Sprintf("Total: %d | Succeeded: %d | Failed: %d | Skipped: %d",
		r.Total, r.Succeeded, r.Failed, r.Skipped)
}

// ---------------------------------------------------------------------------
// Orchestrator
// ---------------------------------------------------------------------------

// Orchestrator manages the full download pipeline.
type Orchestrator struct {
	cfg     Config
	session *gohttp.SessionManager
	drm     *drm.Manager
	logger  *slog.Logger
}

// NewOrchestrator creates a download orchestrator.
//
// Parameters:
//   - cfg: Download configuration.
//   - session: HTTP session for making requests.
//
// Returns:
//   - Configured Orchestrator.
func NewOrchestrator(cfg Config, session *gohttp.SessionManager) *Orchestrator {
	return &Orchestrator{
		cfg:     cfg,
		session: session,
		logger:  cfg.Logger,
	}
}

// Run executes the download pipeline for the given media items.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - media: Media items to download.
//
// Returns:
//   - Result with download outcomes, or error for pipeline-level failures.
func (o *Orchestrator) Run(ctx context.Context, media []*model.Media) (*Result, error) {
	result := &Result{Total: len(media)}

	if len(media) == 0 {
		return result, nil
	}

	// Create job channel and dispatch workers.
	jobs := make(chan *model.Media, len(media))
	var wg sync.WaitGroup

	workers := o.cfg.Workers
	if workers <= 0 {
		workers = 1
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for m := range jobs {
				if ctx.Err() != nil {
					result.AddSkipped()
					continue
				}
				o.downloadOne(ctx, m, result)
			}
		}()
	}

	// Enqueue all media.
	for _, m := range media {
		jobs <- m
	}
	close(jobs)

	wg.Wait()
	return result, nil
}

// downloadOne handles a single media download.
func (o *Orchestrator) downloadOne(ctx context.Context, m *model.Media, result *Result) {
	if !m.IsLinked() {
		result.AddSkipped()
		m.MarkDownloadSkipped()
		return
	}

	var err error
	if m.IsProtected() {
		err = o.downloadProtected(ctx, m)
	} else {
		err = o.downloadNormal(ctx, m)
	}

	if err != nil {
		result.AddFailure(m.ID, m.Link(), err)
		m.MarkDownloadFailed()
		if o.logger != nil {
			o.logger.Error("download failed",
				"media_id", m.ID,
				"error", err,
			)
		}
	} else {
		result.AddSuccess()
		m.MarkDownloadSucceeded()
	}
}
