// =============================================================================
// FILE: internal/cli/hide_args.go
// PURPOSE: HideArgs redacts sensitive CLI arguments from log output so that
//          private keys, tokens, and other secrets are not leaked.
// =============================================================================

package cli

import (
	"strings"
)

// sensitiveFlags lists flag names whose values should be redacted in log output.
var sensitiveFlags = map[string]bool{
	"--private-key": true,
	"--discord":     true,
	"--config":      true,
}

// HideArgs returns a copy of the argument slice with sensitive flag values
// replaced by "***REDACTED***".
func HideArgs(args []string) []string {
	out := make([]string, len(args))
	copy(out, args)

	for i := 0; i < len(out); i++ {
		// Handle --flag=value form.
		if idx := strings.Index(out[i], "="); idx > 0 {
			name := out[i][:idx]
			if sensitiveFlags[name] {
				out[i] = name + "=***REDACTED***"
			}
			continue
		}
		// Handle --flag value form.
		if sensitiveFlags[out[i]] && i+1 < len(out) {
			i++
			out[i] = "***REDACTED***"
		}
	}
	return out
}
