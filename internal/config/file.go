// =============================================================================
// FILE: internal/config/file.go
// PURPOSE: Config file I/O operations. Handles reading, writing, and creating
//          config files on disk. Ports Python utils/config/file.py.
// =============================================================================

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// File operations
// ---------------------------------------------------------------------------

// ReadConfigFile reads and parses a config JSON file from the given path.
//
// Parameters:
//   - path: Absolute path to the config JSON file.
//
// Returns:
//   - The parsed AppConfig, and any error.
func ReadConfigFile(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := DefaultConfig()

	// Handle potential "config" wrapper key
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		if inner, ok := raw["config"]; ok {
			data = inner
		}
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return &cfg, nil
}

// WriteConfigFile writes a config struct to a JSON file at the given path.
// Creates parent directories if they don't exist.
//
// Parameters:
//   - path: Absolute path to write the config file.
//   - cfg: The configuration to serialize.
//
// Returns:
//   - Error if directory creation or file write fails.
func WriteConfigFile(path string, cfg *AppConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	return nil
}

// EnsureConfigExists creates the config file with defaults if it doesn't exist.
//
// Parameters:
//   - path: The config file path to check/create.
//
// Returns:
//   - true if the file was created, false if it already existed, and any error.
func EnsureConfigExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return false, nil // Already exists
	}

	cfg := DefaultConfig()
	if err := WriteConfigFile(path, &cfg); err != nil {
		return false, err
	}

	return true, nil
}
