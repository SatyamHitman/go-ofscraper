// =============================================================================
// FILE: internal/app/shutdown.go
// PURPOSE: Graceful shutdown. Coordinates cleanup of all subsystems in the
//          correct order. Ports Python main/close/exit.py.
// =============================================================================

package app

import (
	"gofscraper/internal/db"
	"gofscraper/internal/logging"
)

// ---------------------------------------------------------------------------
// Shutdown
// ---------------------------------------------------------------------------

// Shutdown performs graceful cleanup of all subsystems.
func (a *App) Shutdown() {
	if a.logger != nil {
		a.logger.Info("shutting down")
	}

	// Close HTTP session.
	if a.session != nil {
		a.session.Close()
	}

	// Close all database connections.
	db.CloseAll()

	// Flush logs.
	logging.Close()

	// Cancel context.
	a.cancel()
}
