// =============================================================================
// FILE: internal/cli/flags/automatic.go
// PURPOSE: Automatic flag definitions: action, daemon.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterAutomaticFlags adds automation-related flags to the given command.
func RegisterAutomaticFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSliceP("action", "a", []string{"download"}, "Actions to perform (download, like, unlike)")
	f.BoolP("daemon", "d", false, "Run in daemon mode with periodic execution")
}
