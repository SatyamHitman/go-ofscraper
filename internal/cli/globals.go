// =============================================================================
// FILE: internal/cli/globals.go
// PURPOSE: Shared CLI state. Global variables and accessors for CLI flag
//          values that need to be shared across commands.
//          Ports Python utils/args/globals.py.
// =============================================================================

package cli

import (
	"sync"
)

// ---------------------------------------------------------------------------
// Global CLI state
// ---------------------------------------------------------------------------

var (
	cliMu         sync.RWMutex
	configPath    string
	profileName   string
	logLevel      string
	noInteractive bool
	verbose       bool
)

// SetGlobals stores parsed CLI flag values for global access.
func SetGlobals(config, profile, level string, noPrompts, verb bool) {
	cliMu.Lock()
	defer cliMu.Unlock()
	configPath = config
	profileName = profile
	logLevel = level
	noInteractive = noPrompts
	verbose = verb
}

// ConfigPath returns the configured config file path.
func ConfigPath() string {
	cliMu.RLock()
	defer cliMu.RUnlock()
	return configPath
}

// ProfileName returns the active profile name.
func ProfileName() string {
	cliMu.RLock()
	defer cliMu.RUnlock()
	return profileName
}

// LogLevel returns the configured log level.
func LogLevel() string {
	cliMu.RLock()
	defer cliMu.RUnlock()
	return logLevel
}

// IsInteractive returns true if interactive mode is enabled.
func IsInteractive() bool {
	cliMu.RLock()
	defer cliMu.RUnlock()
	return !noInteractive
}

// IsVerbose returns true if verbose output is enabled.
func IsVerbose() bool {
	cliMu.RLock()
	defer cliMu.RUnlock()
	return verbose
}
