// =============================================================================
// FILE: internal/scripts/scripts.go
// PURPOSE: Script runner interface and manager. Defines the ScriptRunner
//          interface and provides the default implementation for executing
//          user-defined scripts at various lifecycle points.
//          Ports Python scripts/ module.
// =============================================================================

package scripts

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// ScriptRunner interface
// ---------------------------------------------------------------------------

// ScriptRunner defines the interface for executing user scripts at lifecycle points.
type ScriptRunner interface {
	RunAfterDownloadAction(ctx context.Context, username string, media []*model.Media, action string) error
	RunAfterDownload(ctx context.Context, filePath string) error
	RunAfterLikeAction(ctx context.Context, username string, posts []model.Post) error
	RunNaming(ctx context.Context, media *model.Media) (string, error)
	RunSkipCheck(ctx context.Context, total int64, media *model.Media) (bool, error)
	RunFinal(ctx context.Context) error
}

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

// Config holds script configuration.
type Config struct {
	AfterDownloadAction string // Script path for after download action
	AfterDownload       string // Script path for after each download
	AfterLikeAction     string // Script path for after like action
	Naming              string // Script path for custom naming
	SkipDownload        string // Script path for skip download check
	Final               string // Script path for final cleanup
	Logger              *slog.Logger
}

// ---------------------------------------------------------------------------
// Manager
// ---------------------------------------------------------------------------

// Manager implements ScriptRunner using shell script execution.
type Manager struct {
	cfg    Config
	logger *slog.Logger
}

// NewManager creates a script manager with the given config.
func NewManager(cfg Config) *Manager {
	return &Manager{
		cfg:    cfg,
		logger: cfg.Logger,
	}
}

// ---------------------------------------------------------------------------
// Script execution helper
// ---------------------------------------------------------------------------

// runScript executes a script with the given arguments and environment.
func (m *Manager) runScript(ctx context.Context, scriptPath string, args []string, env map[string]string) (string, error) {
	if scriptPath == "" {
		return "", nil
	}

	cmd := exec.CommandContext(ctx, scriptPath, args...)

	// Set environment variables.
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("script %s failed: %w\noutput: %s", scriptPath, err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
