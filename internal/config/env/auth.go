// =============================================================================
// FILE: internal/config/env/auth.go
// PURPOSE: Auth-related environment variable defaults. Controls auth warning
//          timeout and related thresholds. Ports Python of_env/values/req/auth.py.
// =============================================================================

package env

import "time"

// AuthWarningTimeout returns the duration before showing auth expiry warnings.
// Default: 30 minutes (1800 seconds).
func AuthWarningTimeout() time.Duration {
	return GetDuration("OF_AUTH_WARNING_TIMEOUT", 1800.0)
}
