// =============================================================================
// FILE: internal/cli/flags/logging.go
// PURPOSE: Logging flag definitions: log-level, discord, profile-log, no-log.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterLoggingFlags adds logging-related flags to the given command.
func RegisterLoggingFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("log-level", "info", "Log level (trace, debug, info, warn, error)")
	f.String("discord", "", "Discord webhook URL for log output")
	f.Bool("profile-log", false, "Enable per-profile log files")
	f.Bool("no-log", false, "Disable file logging entirely")
}
