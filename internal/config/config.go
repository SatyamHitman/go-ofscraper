// =============================================================================
// FILE: internal/config/config.go
// PURPOSE: Core configuration management using Viper. Handles loading config
//          from JSON files, merging with defaults, and providing a global
//          config accessor. Ports Python utils/config/config.py.
// =============================================================================

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Global config singleton
// ---------------------------------------------------------------------------

var (
	// globalConfig holds the current active configuration.
	globalConfig *AppConfig

	// configMu protects concurrent access to globalConfig.
	configMu sync.RWMutex

	// configPath stores the resolved config file path.
	configPath string
)

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

// Init loads the configuration from the config file, merging with defaults.
// Must be called before any other config functions.
//
// Parameters:
//   - customPath: Optional custom config file path. If empty, uses the default
//     location at ConfigDir()/config.json.
//
// Returns:
//   - Error if the config file cannot be read or parsed.
func Init(customPath string) error {
	configMu.Lock()
	defer configMu.Unlock()

	// Start with defaults
	cfg := DefaultConfig()

	// Resolve config file path
	if customPath != "" {
		configPath = customPath
	} else {
		envPath := env.ConfigFilePath()
		if envPath != "" {
			configPath = envPath
		} else {
			configPath = filepath.Join(env.ConfigDir(), "config.json")
		}
	}

	// Read config file if it exists
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file — use defaults
			globalConfig = &cfg
			return nil
		}
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse JSON — the file may have a wrapper "config" key
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Check for "config" wrapper key (Python format compatibility)
	configData := data
	if inner, ok := raw["config"]; ok {
		configData = inner
	}

	// Merge file values over defaults
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to merge config values: %w", err)
	}

	globalConfig = &cfg
	return nil
}

// ---------------------------------------------------------------------------
// Accessors
// ---------------------------------------------------------------------------

// Get returns the current active configuration. Panics if Init has not been called.
//
// Returns:
//   - Pointer to the current AppConfig.
func Get() *AppConfig {
	configMu.RLock()
	defer configMu.RUnlock()

	if globalConfig == nil {
		// Return defaults if not initialized (safe fallback)
		cfg := DefaultConfig()
		return &cfg
	}
	return globalConfig
}

// ConfigPath returns the resolved config file path.
//
// Returns:
//   - The absolute path to the active config file.
func ConfigPath() string {
	configMu.RLock()
	defer configMu.RUnlock()
	return configPath
}

// ConfigDir returns the directory containing the config file.
//
// Returns:
//   - The directory path.
func ConfigDirPath() string {
	return filepath.Dir(ConfigPath())
}

// ---------------------------------------------------------------------------
// Update operations
// ---------------------------------------------------------------------------

// Update replaces the global configuration with the provided config and
// writes it to disk.
//
// Parameters:
//   - cfg: The new configuration to save.
//
// Returns:
//   - Error if the config file cannot be written.
func Update(cfg *AppConfig) error {
	configMu.Lock()
	defer configMu.Unlock()

	globalConfig = cfg
	return writeConfigFile(cfg)
}

// UpdateField updates a single top-level configuration field by key.
// The config is re-serialized and written to disk.
//
// Parameters:
//   - field: The JSON field name to update (e.g., "discord").
//   - value: The new value for the field.
//
// Returns:
//   - Error if the update or file write fails.
func UpdateField(field string, value interface{}) error {
	configMu.Lock()
	defer configMu.Unlock()

	if globalConfig == nil {
		cfg := DefaultConfig()
		globalConfig = &cfg
	}

	// Serialize current config to map, update field, deserialize back
	data, err := json.Marshal(globalConfig)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	var cfgMap map[string]interface{}
	if err := json.Unmarshal(data, &cfgMap); err != nil {
		return fmt.Errorf("failed to deserialize config map: %w", err)
	}

	cfgMap[field] = value

	updatedData, err := json.Marshal(cfgMap)
	if err != nil {
		return fmt.Errorf("failed to serialize updated config: %w", err)
	}

	var updatedCfg AppConfig
	if err := json.Unmarshal(updatedData, &updatedCfg); err != nil {
		return fmt.Errorf("failed to parse updated config: %w", err)
	}

	globalConfig = &updatedCfg
	return writeConfigFile(&updatedCfg)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// writeConfigFile serializes the config to JSON and writes it to the config path.
//
// Parameters:
//   - cfg: The configuration to write.
//
// Returns:
//   - Error if directory creation or file write fails.
func writeConfigFile(cfg *AppConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	return nil
}
