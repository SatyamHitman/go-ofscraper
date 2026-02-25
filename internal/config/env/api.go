// =============================================================================
// FILE: internal/config/env/api.go
// PURPOSE: API-specific constants and configuration defaults.
//          Ports Python of_env/values/req/api.py.
// =============================================================================

package env

// APIPageLimit returns the default page size for API pagination.
func APIPageLimit() int {
	return GetInt("OF_API_PAGE_LIMIT", 100)
}

// SubscriptionPageLimit returns the page size for subscription pagination.
func SubscriptionPageLimit() int {
	return GetInt("OF_SUBSCRIPTION_PAGE_LIMIT", 10)
}

// HighlightPageLimit returns the page size for highlight pagination.
func HighlightPageLimit() int {
	return GetInt("OF_HIGHLIGHT_PAGE_LIMIT", 5)
}
