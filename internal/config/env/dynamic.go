// =============================================================================
// FILE: internal/config/env/dynamic.go
// PURPOSE: Dynamic rule provider environment variable defaults.
//          Ports Python of_env/values/dynamic.py.
// =============================================================================

package env

// DynamicRuleDefault returns the default dynamic rule provider name.
func DynamicRuleDefault() string {
	return GetString("OF_DYNAMIC_RULE_DEFAULT", "digital")
}

// DynamicRuleManual returns the manual dynamic rule JSON string (if set).
func DynamicRuleManual() string {
	return GetString("OF_DYNAMIC_RULE_MANUAL", "")
}
