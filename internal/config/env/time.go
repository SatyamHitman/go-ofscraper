// =============================================================================
// FILE: internal/config/env/time.go
// PURPOSE: Defines all time-related and cache expiry constants. Controls how
//          long cached data remains valid, database operation intervals, and
//          profile data expiry. Ports Python of_env/values/time.py.
// =============================================================================

package env

import "time"

// ---------------------------------------------------------------------------
// Cache expiry durations
// ---------------------------------------------------------------------------

// HourlyExpiry returns the hourly cache expiry duration (1 hour).
func HourlyExpiry() time.Duration {
	return GetDuration("OF_HOURLY_EXPIRY", 3600)
}

// ThirtyExpiry returns the 30-minute cache expiry duration.
func ThirtyExpiry() time.Duration {
	return GetDuration("OF_THIRTY_EXPIRY", 1800)
}

// SizeTimeout returns the cache timeout for file size lookups (2 weeks).
func SizeTimeout() time.Duration {
	return GetDuration("OF_SIZE_TIMEOUT", 1209600)
}

// DatabaseTimeout returns the database operation timeout (5 minutes).
func DatabaseTimeout() time.Duration {
	return GetDuration("OF_DATABASE_TIMEOUT", 300)
}

// KeyExpiry returns the DRM key cache expiry (0 = no expiry).
func KeyExpiry() time.Duration {
	return GetDuration("OF_KEY_EXPIRY", 0)
}

// ResponseExpiry returns the API response cache expiry duration.
func ResponseExpiry() time.Duration {
	return GetDuration("OF_RESPONSE_EXPIRY", 5000000)
}

// ---------------------------------------------------------------------------
// Interval durations
// ---------------------------------------------------------------------------

// DBInterval returns the database maintenance interval (1 day).
func DBInterval() time.Duration {
	return GetDuration("OF_DB_INTERVAL", 86400)
}

// DaySeconds is the number of seconds in one day.
const DaySeconds = 86400

// ThreeDaySeconds is the number of seconds in three days.
const ThreeDaySeconds = 259200

// ---------------------------------------------------------------------------
// Profile data expiry
// ---------------------------------------------------------------------------

// ProfileDataExpiry returns the profile data cache expiry (1 day).
func ProfileDataExpiry() time.Duration {
	return GetDuration("OF_PROFILE_DATA_EXPIRY", 86400)
}

// ProfileDataExpiryAsync returns the async profile data cache expiry (1 day).
func ProfileDataExpiryAsync() time.Duration {
	return GetDuration("OF_PROFILE_DATA_EXPIRY_ASYNC", 86400)
}
