// =============================================================================
// FILE: internal/config/env/anon.go
// PURPOSE: Anonymous mode configuration defaults. Provides the user agent
//          string used when running without authentication. Ports Python
//          of_env/values/req/anon.py.
// =============================================================================

package env

// AnonUserAgent returns the user agent string for anonymous/unauthenticated requests.
func AnonUserAgent() string {
	return GetString("OF_ANON_USERAGENT",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
}
