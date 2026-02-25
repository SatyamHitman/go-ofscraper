// =============================================================================
// FILE: internal/logging/classes.go
// PURPOSE: Log class definitions and named logger registry. Provides
//          pre-configured loggers for each subsystem (api, download, auth,
//          db, etc.) with consistent attribute tagging. Ports Python
//          utils/logs/classes/classes.py.
// =============================================================================

package logging

import (
	"log/slog"
	"sync"
)

// ---------------------------------------------------------------------------
// Named logger registry
// ---------------------------------------------------------------------------

var (
	// namedLoggers caches subsystem loggers keyed by component name.
	namedLoggers sync.Map
)

// Named returns (or creates) a logger tagged with the given component name.
// The logger is cached so subsequent calls with the same name are cheap.
//
// Parameters:
//   - component: Subsystem identifier (e.g. "api", "download", "auth").
//
// Returns:
//   - A *slog.Logger with the "component" attribute set.
func Named(component string) *slog.Logger {
	if cached, ok := namedLoggers.Load(component); ok {
		return cached.(*slog.Logger)
	}

	l := Logger().With(slog.String("component", component))
	namedLoggers.Store(component, l)
	return l
}

// ---------------------------------------------------------------------------
// Well-known component loggers
// ---------------------------------------------------------------------------

// Pre-defined component names for consistent use across the codebase.
const (
	ComponentAPI      = "api"
	ComponentAuth     = "auth"
	ComponentConfig   = "config"
	ComponentDB       = "database"
	ComponentDownload = "download"
	ComponentDRM      = "drm"
	ComponentFilter   = "filter"
	ComponentHTTP     = "http"
	ComponentModel    = "model"
	ComponentScript   = "script"
	ComponentTUI      = "tui"
	ComponentCache    = "cache"
	ComponentWorker   = "worker"
)

// Convenience functions that return pre-tagged loggers for each subsystem.

// API returns the API subsystem logger.
func API() *slog.Logger { return Named(ComponentAPI) }

// Auth returns the auth subsystem logger.
func Auth() *slog.Logger { return Named(ComponentAuth) }

// Config returns the config subsystem logger.
func Config() *slog.Logger { return Named(ComponentConfig) }

// DB returns the database subsystem logger.
func DB() *slog.Logger { return Named(ComponentDB) }

// Download returns the download subsystem logger.
func Download() *slog.Logger { return Named(ComponentDownload) }

// DRM returns the DRM subsystem logger.
func DRM() *slog.Logger { return Named(ComponentDRM) }

// Filter returns the filter subsystem logger.
func Filter() *slog.Logger { return Named(ComponentFilter) }

// HTTP returns the HTTP client subsystem logger.
func HTTP() *slog.Logger { return Named(ComponentHTTP) }

// Model returns the model subsystem logger.
func Model() *slog.Logger { return Named(ComponentModel) }

// Script returns the script subsystem logger.
func Script() *slog.Logger { return Named(ComponentScript) }

// TUI returns the TUI subsystem logger.
func TUI() *slog.Logger { return Named(ComponentTUI) }

// Cache returns the cache subsystem logger.
func Cache() *slog.Logger { return Named(ComponentCache) }

// Worker returns the worker subsystem logger.
func Worker() *slog.Logger { return Named(ComponentWorker) }
