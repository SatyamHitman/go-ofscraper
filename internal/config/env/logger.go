// =============================================================================
// FILE: internal/config/env/logger.go
// PURPOSE: Logger environment variable defaults.
//          Ports Python of_env/values/logger.py.
// =============================================================================

package env

// LogLevel returns the configured log level string.
func LogLevel() string {
	return GetString("OF_LOG_LEVEL", "DEBUG")
}

// LogExpireTime returns the log file expiry time in seconds.
func LogExpireTime() int {
	return GetInt("OF_LOG_EXPIRE_TIME", 0)
}
