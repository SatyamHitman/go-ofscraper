// =============================================================================
// FILE: internal/config/env/url_other.go
// PURPOSE: Other URL environment variable defaults not covered by the main
//          endpoint or dynamic URL files. Ports Python of_env/values/url/other_url.py.
// =============================================================================

package env

// CDRMServiceURL returns the CDRM service URL for remote DRM key extraction.
func CDRMServiceURL() string {
	return GetString("OF_CDRM_SERVICE_URL", "")
}
