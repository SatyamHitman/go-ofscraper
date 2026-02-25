// =============================================================================
// FILE: internal/commands/scraper/scraper.go
// PURPOSE: Scraper orchestrator. Manages the full scrape pipeline: load
//          subscriptions, filter users, and for each user fetch posts, filter
//          posts, filter media, and dispatch to download/like/metadata handlers.
//          Ports Python runner/open/scrape_paid.py and runner/scraper.py.
// =============================================================================

package scraper

import (
	"context"
	"fmt"
	"log/slog"

	"gofscraper/internal/app"
	cmdutils "gofscraper/internal/commands/utils"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Scraper
// ---------------------------------------------------------------------------

// Scraper orchestrates a complete scrape run across one or more users.
type Scraper struct {
	logger  *slog.Logger
	scrCtx  *cmdutils.ScrapeContext
	actions []string
	areas   []string
}

// New creates a new Scraper with the given configuration.
//
// Parameters:
//   - logger: Structured logger for output.
//   - actions: Actions to perform (e.g., "download", "like", "metadata").
//   - areas: Content areas to scrape (e.g., "timeline", "messages", "stories").
//
// Returns:
//   - A configured Scraper.
func New(logger *slog.Logger, actions, areas []string) *Scraper {
	if logger == nil {
		logger = slog.Default()
	}
	return &Scraper{
		logger:  logger,
		scrCtx:  cmdutils.NewScrapeContext(),
		actions: actions,
		areas:   areas,
	}
}

// Run executes the full scrape pipeline.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance providing session and config.
//
// Returns:
//   - Error if any stage of the pipeline fails fatally.
func (s *Scraper) Run(ctx context.Context, a *app.App) error {
	s.logger.Info("scraper starting",
		"actions", s.actions,
		"areas", s.areas,
	)

	// Stage 1: Prepare data — resolve users and apply filters.
	users, err := PrepareData(ctx, a, s.logger)
	if err != nil {
		return fmt.Errorf("prepare data: %w", err)
	}
	if len(users) == 0 {
		s.logger.Info(cmdutils.MsgNoUsers)
		return nil
	}
	s.scrCtx.Users = users

	// Stage 2: Process each user.
	for i, user := range users {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		s.logger.Info(cmdutils.FormatUserProgress(user.Name, i+1, len(users)))
		if err := s.processUser(ctx, a, user); err != nil {
			s.logger.Error("failed to process user",
				"user", user.Name,
				"error", err,
			)
			s.scrCtx.Errors.Add(1)
			continue
		}
		s.scrCtx.UsersProcessed.Add(1)
	}

	// Stage 3: Print summary.
	PrintSummary(s.scrCtx, s.logger)

	s.logger.Info("scraper finished")
	return nil
}

// processUser handles the scrape pipeline for a single user:
// fetch posts, filter, and dispatch actions.
func (s *Scraper) processUser(ctx context.Context, a *app.App, user *model.User) error {
	// Fetch posts for all configured areas.
	var posts []*model.Post
	for _, area := range s.areas {
		s.logger.Debug(fmt.Sprintf(cmdutils.MsgFetchingPosts, area, user.Name))
		// TODO: Wire to API post fetcher per area.
		_ = area
	}

	if len(posts) == 0 {
		s.logger.Info(fmt.Sprintf(cmdutils.MsgNoPosts, user.Name))
		return nil
	}

	s.scrCtx.AddPosts(posts)

	// Collect all media from posts.
	var allMedia []*model.Media
	for _, post := range posts {
		allMedia = append(allMedia, post.ViewableMedia()...)
	}
	s.scrCtx.MediaFound.Add(int64(len(allMedia)))

	// Dispatch configured actions.
	for _, action := range s.actions {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := a.RunAction(ctx, action, s.areas, []string{user.Name}); err != nil {
			return fmt.Errorf("action %s for user %s: %w", action, user.Name, err)
		}
	}

	return nil
}

// Context returns the scrape context with accumulated results.
//
// Returns:
//   - The ScrapeContext with current counters and data.
func (s *Scraper) Context() *cmdutils.ScrapeContext {
	return s.scrCtx
}
