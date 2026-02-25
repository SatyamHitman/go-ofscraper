// =============================================================================
// FILE: internal/logging/logging.go
// PURPOSE: Central logger setup for gofscraper. Initialises the slog-based
//          logging system with configurable handlers (stdout, file, discord).
//          Supports log rotation, coloured terminal output via tint, and
//          sensitive-data redaction. Ports Python utils/logs/logger.py.
// =============================================================================

package logging

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Package-level state
// ---------------------------------------------------------------------------

var (
	// rootLogger is the package-level logger used by all subsystems.
	rootLogger *slog.Logger

	// logMu protects rootLogger during initialisation.
	logMu sync.RWMutex

	// initialised tracks whether Init has been called.
	initialised bool
)

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

// Init sets up the logging system. It creates handlers based on the current
// configuration (log level, file output, discord webhook, colour mode) and
// installs the resulting logger as both the package default and slog default.
//
// Parameters:
//   - opts: Optional overrides. Pass nil to use defaults from env/config.
//
// Returns:
//   - Error if handler creation fails (e.g. cannot open log file).
func Init(opts *Options) error {
	logMu.Lock()
	defer logMu.Unlock()

	if opts == nil {
		opts = defaultOptions()
	}

	// Parse the configured log level.
	level := parseLevel(opts.Level)

	// Build the list of handlers.
	var handlers []slog.Handler

	// 1. Stdout handler (always present).
	stdoutHandler := newStdoutHandler(os.Stdout, level, opts.Color)
	handlers = append(handlers, stdoutHandler)

	// 2. File handler (if a log directory is configured).
	if opts.LogDir != "" {
		fh, err := newFileHandler(opts.LogDir, level, opts.RotateLogs)
		if err != nil {
			return err
		}
		handlers = append(handlers, fh)
	}

	// 3. Discord handler (if webhook URL is configured).
	webhookURL := env.DiscordWebhookURL()
	if webhookURL != "" {
		dh := newDiscordHandler(webhookURL, level)
		handlers = append(handlers, dh)
	}

	// Wrap all handlers in a sensitive-data redactor.
	var redacted []slog.Handler
	for _, h := range handlers {
		redacted = append(redacted, newSensitiveHandler(h))
	}

	// Combine into a single multi-handler.
	combined := newMultiHandler(redacted...)

	rootLogger = slog.New(combined)
	slog.SetDefault(rootLogger)
	initialised = true

	return nil
}

// ---------------------------------------------------------------------------
// Accessors
// ---------------------------------------------------------------------------

// Logger returns the package-level logger. If Init has not been called it
// returns a no-op logger that writes to io.Discard.
//
// Returns:
//   - The active *slog.Logger.
func Logger() *slog.Logger {
	logMu.RLock()
	defer logMu.RUnlock()

	if rootLogger == nil {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return rootLogger
}

// IsInitialised reports whether the logging system has been set up.
//
// Returns:
//   - true if Init completed successfully.
func IsInitialised() bool {
	logMu.RLock()
	defer logMu.RUnlock()
	return initialised
}

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

// Options configures the logging subsystem.
type Options struct {
	Level      string // slog level name: "DEBUG", "INFO", "WARN", "ERROR"
	LogDir     string // Directory for log files (empty = no file logging)
	Color      bool   // Enable coloured terminal output
	RotateLogs bool   // Enable log file rotation
}

// defaultOptions builds Options from env vars and config defaults.
func defaultOptions() *Options {
	return &Options{
		Level:      env.LogLevel(),
		LogDir:     "", // caller may set from config
		Color:      true,
		RotateLogs: true,
	}
}

// ---------------------------------------------------------------------------
// Level parsing
// ---------------------------------------------------------------------------

// parseLevel converts a level name string to slog.Level. Unrecognised values
// default to slog.LevelInfo.
//
// Parameters:
//   - name: Level name (case-insensitive match via slog).
//
// Returns:
//   - The corresponding slog.Level.
func parseLevel(name string) slog.Level {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(name)); err != nil {
		return slog.LevelInfo
	}
	return lvl
}
