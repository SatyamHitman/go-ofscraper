// =============================================================================
// FILE: internal/config/env/env.go
// PURPOSE: Central environment variable loader. Provides functions to read
//          configuration values from environment variables with typed defaults.
//          Ports Python utils/of_env/of_env.py and load.py. All env var
//          reading for the application routes through this package.
// =============================================================================

package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// String getters
// ---------------------------------------------------------------------------

// GetString reads an environment variable as a string. Returns the default
// value if the variable is not set or is empty.
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback value if the variable is not set.
//
// Returns:
//   - The environment variable value, or defaultVal.
func GetString(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// ---------------------------------------------------------------------------
// Integer getters
// ---------------------------------------------------------------------------

// GetInt reads an environment variable as an integer. Returns the default
// value if the variable is not set, empty, or not a valid integer.
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback value.
//
// Returns:
//   - The parsed integer value, or defaultVal.
func GetInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// GetInt64 reads an environment variable as an int64.
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback value.
//
// Returns:
//   - The parsed int64 value, or defaultVal.
func GetInt64(key string, defaultVal int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// ---------------------------------------------------------------------------
// Float getters
// ---------------------------------------------------------------------------

// GetFloat64 reads an environment variable as a float64.
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback value.
//
// Returns:
//   - The parsed float64 value, or defaultVal.
func GetFloat64(key string, defaultVal float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// ---------------------------------------------------------------------------
// Boolean getters
// ---------------------------------------------------------------------------

// GetBool reads an environment variable as a boolean.
// Recognized true values: "true", "1", "yes", "on" (case-insensitive).
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback value.
//
// Returns:
//   - The parsed boolean value, or defaultVal.
func GetBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	switch strings.ToLower(val) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	return defaultVal
}

// ---------------------------------------------------------------------------
// Duration getters
// ---------------------------------------------------------------------------

// GetDuration reads an environment variable as a time.Duration (from seconds).
//
// Parameters:
//   - key: The environment variable name.
//   - defaultSeconds: The fallback value in seconds.
//
// Returns:
//   - The parsed duration, or default duration from defaultSeconds.
func GetDuration(key string, defaultSeconds float64) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return time.Duration(defaultSeconds * float64(time.Second))
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return time.Duration(defaultSeconds * float64(time.Second))
	}
	return time.Duration(parsed * float64(time.Second))
}

// ---------------------------------------------------------------------------
// List getters
// ---------------------------------------------------------------------------

// GetStringList reads an environment variable as a comma-separated string list.
//
// Parameters:
//   - key: The environment variable name.
//   - defaultVal: The fallback list.
//
// Returns:
//   - The parsed string slice, or defaultVal.
func GetStringList(key string, defaultVal []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultVal
	}
	return result
}
