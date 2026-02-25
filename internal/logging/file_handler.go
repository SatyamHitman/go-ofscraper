// =============================================================================
// FILE: internal/logging/file_handler.go
// PURPOSE: File-based log handler. Writes JSON-structured log lines to a
//          rotating log file. Supports automatic rotation by size/age.
//          Ports Python utils/logs/classes/handlers/file.py.
// =============================================================================

package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// File handler
// ---------------------------------------------------------------------------

// fileHandler wraps a slog.JSONHandler writing to a rotatable file.
type fileHandler struct {
	handler slog.Handler
	file    *os.File
	logDir  string
	rotate  bool
}

// newFileHandler creates a file-based log handler that writes JSON lines
// to a timestamped file inside logDir. If rotate is true, old log files
// are pruned during creation.
//
// Parameters:
//   - logDir: Directory where log files are stored.
//   - level: Minimum log level.
//   - rotate: Whether to remove old log files.
//
// Returns:
//   - A slog.Handler backed by the log file, and any error.
func newFileHandler(logDir string, level slog.Level, rotate bool) (slog.Handler, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}

	// Prune old logs before opening a new one.
	if rotate {
		pruneOldLogs(logDir)
	}

	// Build a timestamped filename: gofscraper_20060102_150405.log
	ts := time.Now().Format("20060102_150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("gofscraper_%s.log", ts))

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", logPath, err)
	}

	h := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: level,
	})

	fh := &fileHandler{
		handler: h,
		file:    f,
		logDir:  logDir,
		rotate:  rotate,
	}

	// Register for shutdown cleanup.
	registerCloseable(fh)

	return fh, nil
}

// Enabled implements slog.Handler.
func (fh *fileHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return fh.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (fh *fileHandler) Handle(ctx context.Context, record slog.Record) error {
	return fh.handler.Handle(ctx, record)
}

// WithAttrs implements slog.Handler.
func (fh *fileHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &fileHandler{
		handler: fh.handler.WithAttrs(attrs),
		file:    fh.file,
		logDir:  fh.logDir,
		rotate:  fh.rotate,
	}
}

// WithGroup implements slog.Handler.
func (fh *fileHandler) WithGroup(name string) slog.Handler {
	return &fileHandler{
		handler: fh.handler.WithGroup(name),
		file:    fh.file,
		logDir:  fh.logDir,
		rotate:  fh.rotate,
	}
}

// Close flushes and closes the underlying log file.
//
// Returns:
//   - Error if the file cannot be synced or closed.
func (fh *fileHandler) Close() error {
	if fh.file == nil {
		return nil
	}
	_ = fh.file.Sync()
	return fh.file.Close()
}

// ---------------------------------------------------------------------------
// Log rotation
// ---------------------------------------------------------------------------

// maxLogFiles is the maximum number of log files to retain.
const maxLogFiles = 5

// pruneOldLogs removes the oldest log files when more than maxLogFiles exist
// in the given directory.
//
// Parameters:
//   - logDir: The directory containing log files.
func pruneOldLogs(logDir string) {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}

	// Collect log files matching our naming pattern.
	var logFiles []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "gofscraper_") && strings.HasSuffix(e.Name(), ".log") {
			logFiles = append(logFiles, e)
		}
	}

	if len(logFiles) <= maxLogFiles {
		return
	}

	// Sort by name (which sorts by timestamp since format is fixed-width).
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Name() < logFiles[j].Name()
	})

	// Remove the oldest files to bring count down to maxLogFiles.
	toRemove := len(logFiles) - maxLogFiles
	for i := 0; i < toRemove; i++ {
		_ = os.Remove(filepath.Join(logDir, logFiles[i].Name()))
	}
}
