// =============================================================================
// FILE: internal/commands/utils/command.go
// PURPOSE: Base command interface and common setup helpers shared across all
//          command implementations. Ports Python utils/args/helpers.py command
//          scaffolding.
// =============================================================================

package cmdutils

import (
	"context"
	"log/slog"
)

// ---------------------------------------------------------------------------
// Command interface
// ---------------------------------------------------------------------------

// Command defines the interface that all application commands must implement.
type Command interface {
	// Name returns the command name for display and routing.
	Name() string

	// Run executes the command with the given context and arguments.
	Run(ctx context.Context, args []string) error
}

// ---------------------------------------------------------------------------
// CommandBase provides common fields shared by all command implementations.
// ---------------------------------------------------------------------------

// CommandBase holds shared state that every command needs.
type CommandBase struct {
	Logger *slog.Logger
}

// NewCommandBase creates a CommandBase with the provided logger.
//
// Parameters:
//   - logger: The structured logger to use for command output.
//
// Returns:
//   - A configured CommandBase.
func NewCommandBase(logger *slog.Logger) CommandBase {
	if logger == nil {
		logger = slog.Default()
	}
	return CommandBase{Logger: logger}
}

// LogStart logs the beginning of a command execution.
//
// Parameters:
//   - name: The command name.
//   - args: The command arguments.
func (cb *CommandBase) LogStart(name string, args []string) {
	cb.Logger.Info("command starting", "command", name, "args", args)
}

// LogDone logs the completion of a command execution.
//
// Parameters:
//   - name: The command name.
func (cb *CommandBase) LogDone(name string) {
	cb.Logger.Info("command completed", "command", name)
}
