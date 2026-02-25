// =============================================================================
// FILE: internal/config/wrapper.go
// PURPOSE: Provides a config reader wrapper that adds fallback handling for
//          missing or invalid config values. Used by config accessor functions
//          to safely read values with defaults. Ports Python
//          utils/config/utils/wrapper.py config_reader decorator.
// =============================================================================

package config

// ---------------------------------------------------------------------------
// ConfigReader wraps typed access to config fields with fallback logic.
// ---------------------------------------------------------------------------

// ConfigReader provides safe access to AppConfig fields with automatic
// fallback to default values when fields are missing or zero-valued.
type ConfigReader struct {
	cfg *AppConfig
}

// NewConfigReader creates a ConfigReader for the given config.
//
// Parameters:
//   - cfg: The AppConfig to wrap. If nil, uses the global config.
//
// Returns:
//   - A new ConfigReader instance.
func NewConfigReader(cfg *AppConfig) *ConfigReader {
	if cfg == nil {
		cfg = Get()
	}
	return &ConfigReader{cfg: cfg}
}

// String reads a string config value with a fallback default.
//
// Parameters:
//   - value: The config value to check.
//   - fallback: The default value if the config value is empty.
//
// Returns:
//   - The config value, or fallback if empty.
func (cr *ConfigReader) String(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// Int reads an integer config value with a fallback default.
//
// Parameters:
//   - value: The config value to check.
//   - fallback: The default value if the config value is zero.
//
// Returns:
//   - The config value, or fallback if zero.
func (cr *ConfigReader) Int(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

// Int64 reads an int64 config value with a fallback default.
//
// Parameters:
//   - value: The config value to check.
//   - fallback: The default value if the config value is zero.
//
// Returns:
//   - The config value, or fallback if zero.
func (cr *ConfigReader) Int64(value, fallback int64) int64 {
	if value == 0 {
		return fallback
	}
	return value
}

// Bool reads a boolean config value. Since Go booleans default to false,
// this returns the config value directly — there's no way to distinguish
// "not set" from "set to false" without using *bool.
//
// Parameters:
//   - value: The config boolean value.
//
// Returns:
//   - The boolean value.
func (cr *ConfigReader) Bool(value bool) bool {
	return value
}

// StringSlice reads a string slice config value with a fallback default.
//
// Parameters:
//   - value: The config slice to check.
//   - fallback: The default slice if the config slice is empty.
//
// Returns:
//   - The config slice, or fallback if empty.
func (cr *ConfigReader) StringSlice(value, fallback []string) []string {
	if len(value) == 0 {
		return fallback
	}
	return value
}
