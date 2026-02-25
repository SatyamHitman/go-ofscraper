// =============================================================================
// FILE: internal/cli/accessors/actions.go
// PURPOSE: Read action and daemon flag values from cobra commands.
// =============================================================================

package accessors

import (
	"github.com/spf13/cobra"
)

// GetAction returns the action flag values from the command.
func GetAction(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("action")
	return v
}

// GetDaemon returns true if daemon mode is enabled.
func GetDaemon(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("daemon")
	return v
}
