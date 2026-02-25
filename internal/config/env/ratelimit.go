// =============================================================================
// FILE: internal/config/env/ratelimit.go
// PURPOSE: Rate limiting configuration defaults. Currently a placeholder for
//          rate limit env vars. Ports Python of_env/values/req/ratelimit.py.
// =============================================================================

package env

// RateLimitEnabled returns whether rate limiting is enabled.
func RateLimitEnabled() bool {
	return GetBool("OF_RATE_LIMIT_ENABLED", true)
}
