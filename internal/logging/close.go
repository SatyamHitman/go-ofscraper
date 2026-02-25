// =============================================================================
// FILE: internal/logging/close.go
// PURPOSE: Graceful shutdown for the logging subsystem. Flushes buffered log
//          entries and closes open file handles. Must be called before process
//          exit to avoid losing final log messages. Ports Python
//          utils/logs/close.py.
// =============================================================================

package logging

import (
	"log/slog"
	"sync"
)

// ---------------------------------------------------------------------------
// Closeable handler tracking
// ---------------------------------------------------------------------------

var (
	// closeables holds handlers that need explicit cleanup on shutdown.
	closeables []Closeable

	// closeMu protects the closeables slice.
	closeMu sync.Mutex
)

// Closeable is implemented by log handlers that hold resources (file handles,
// network connections) requiring explicit release.
type Closeable interface {
	// Close flushes pending data and releases resources. It should be safe
	// to call multiple times.
	Close() error
}

// registerCloseable adds a handler to the shutdown list.
//
// Parameters:
//   - c: The closeable handler to track.
func registerCloseable(c Closeable) {
	closeMu.Lock()
	defer closeMu.Unlock()
	closeables = append(closeables, c)
}

// ---------------------------------------------------------------------------
// Shutdown
// ---------------------------------------------------------------------------

// Close shuts down the logging subsystem. It flushes and closes all registered
// handlers (file writers, discord buffers, etc.). Should be called once during
// application shutdown, typically via defer in main().
//
// Returns:
//   - The first error encountered during close, or nil.
func Close() error {
	closeMu.Lock()
	defer closeMu.Unlock()

	var firstErr error
	for _, c := range closeables {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	// Reset state so Init can be called again if needed (e.g. tests).
	closeables = nil
	logMu.Lock()
	rootLogger = nil
	initialised = false
	logMu.Unlock()

	slog.Info("logging subsystem shut down")

	return firstErr
}
