// =============================================================================
// FILE: internal/cli/flags/advanced_program.go
// PURPOSE: Advanced program flag definitions: no-cache, no-cache-api,
//          key-mode, private-key, check-interval, auto-like, dynamic-rules,
//          update-profile.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterAdvancedProgramFlags adds advanced program flags to the given command.
func RegisterAdvancedProgramFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.Bool("no-cache", false, "Disable local file cache")
	f.Bool("no-cache-api", false, "Disable API response caching")
	f.String("key-mode", "auto", "Key resolution mode (auto, manual, keydb)")
	f.String("private-key", "", "Path to private key file for decryption")
	f.Int("check-interval", 0, "Interval in seconds between daemon checks (0 = default)")
	f.Bool("auto-like", false, "Automatically like posts when downloading")
	f.String("dynamic-rules", "", "URL or path to dynamic signing rules")
	f.Bool("update-profile", false, "Update user profile metadata during scraping")
}
