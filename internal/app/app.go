// =============================================================================
// FILE: internal/app/app.go
// PURPOSE: App lifecycle orchestrator. Manages the full application lifecycle:
//          initialization, configuration loading, auth setup, command dispatch,
//          and graceful shutdown. Ports Python main/open/load.py and run.py.
// =============================================================================

package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gofscraper/internal/auth"
	"gofscraper/internal/config"
	gohttp "gofscraper/internal/http"
	"gofscraper/internal/logging"
)

// ---------------------------------------------------------------------------
// App
// ---------------------------------------------------------------------------

// App is the main application instance.
type App struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *config.AppConfig
	session *gohttp.SessionManager
	logger  *slog.Logger
}

// New creates a new App instance.
func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Init initializes all application subsystems.
//
// Returns:
//   - Error if any subsystem fails to initialize.
func (a *App) Init() error {
	// Step 1: Load config.
	if err := config.Init(""); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	a.cfg = config.Get()

	// Step 2: Initialize logging.
	if err := logging.Init(&logging.Options{
		Level: "info",
	}); err != nil {
		return fmt.Errorf("init logging: %w", err)
	}
	a.logger = logging.Logger()

	a.logger.Info("gofscraper starting")

	// Step 3: Load auth.
	authData, err := auth.Load("")
	if err != nil {
		a.logger.Warn("auth not loaded", "error", err)
	}

	// Step 4: Create HTTP session.
	a.session = gohttp.New(authData)

	// Step 5: Setup signal handling.
	a.setupSignals()

	return nil
}

// setupSignals configures graceful shutdown on SIGINT/SIGTERM.
func (a *App) setupSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		a.logger.Info("received signal, shutting down", "signal", sig)
		a.cancel()
	}()
}

// Context returns the app's context.
func (a *App) Context() context.Context {
	return a.ctx
}

// Session returns the HTTP session manager.
func (a *App) Session() *gohttp.SessionManager {
	return a.session
}

// Config returns the app config.
func (a *App) Config() *config.AppConfig {
	return a.cfg
}

// Logger returns the app logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}
