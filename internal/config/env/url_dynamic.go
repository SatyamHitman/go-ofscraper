// =============================================================================
// FILE: internal/config/env/url_dynamic.go
// PURPOSE: Defines URLs for dynamic rule sources used in auth signing.
//          Each provider hosts a JSON file containing signing parameters.
//          Ports Python of_env/values/url/dynamic.py.
// =============================================================================

package env

// ---------------------------------------------------------------------------
// Dynamic rule source URLs
// ---------------------------------------------------------------------------

// XaglerURL returns the Xagler dynamic rules URL.
func XaglerURL() string {
	return GetString("OF_XAGLER_URL",
		"https://raw.githubusercontent.com/xagler/dynamic-rules/main/onlyfans.json")
}

// RafaURL returns the Rafa dynamic rules URL.
func RafaURL() string {
	return GetString("OF_RAFA_URL",
		"https://raw.githubusercontent.com/rafa-9/dynamic-rules/main/rules.json")
}

// DigitalCriminalsURL returns the DigitalCriminals (DATAHOARDERS) dynamic rules URL.
func DigitalCriminalsURL() string {
	return GetString("OF_DIGITALCRIMINALS_URL",
		"https://raw.githubusercontent.com/DATAHOARDERS/dynamic-rules/main/onlyfans.json")
}

// DatawhoresURL returns the Datawhores dynamic rules URL.
func DatawhoresURL() string {
	return GetString("OF_DATAWHORES_URL",
		"https://raw.githubusercontent.com/datawhores/onlyfans-dynamic-rules/main/dynamicRules.json")
}

// DeviintURL returns the Deviint dynamic rules URL.
func DeviintURL() string {
	return GetString("OF_DEVIINT_URL",
		"https://raw.githubusercontent.com/deviint/onlyfans-dynamic-rules/main/dynamicRules.json")
}

// RileyURL returns the Riley Access Labs dynamic rules URL.
func RileyURL() string {
	return GetString("OF_RILEY_URL",
		"https://raw.githubusercontent.com/riley-access-labs/onlyfans-dynamic-rules-1/refs/heads/main/dynamicRules.json")
}

// DynamicGenericURL returns the user-configured generic dynamic rules URL.
// Empty by default — must be set via environment variable for "generic" mode.
func DynamicGenericURL() string {
	return GetString("OF_DYNAMIC_GENERIC_URL", "")
}
