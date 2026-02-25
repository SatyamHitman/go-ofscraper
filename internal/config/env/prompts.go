// =============================================================================
// FILE: internal/config/env/prompts.go
// PURPOSE: Prompt-related environment variable defaults.
//          Ports Python of_env/values/prompts.py.
// =============================================================================

package env

// PromptsEnabled returns whether interactive prompts are enabled.
func PromptsEnabled() bool {
	return GetBool("OF_PROMPTS_ENABLED", true)
}
