// =============================================================================
// FILE: internal/cli/bundles/check_utils.go
// PURPOSE: Shared check bundle utilities.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"
)

// DefaultCheckAreas returns the default content areas for check commands.
func DefaultCheckAreas() []string {
	return []string{"timeline", "messages", "archived", "stories", "highlights"}
}

// SetCheckDefaults applies default flag values for check commands if the user
// did not provide them explicitly.
func SetCheckDefaults(cmd *cobra.Command) {
	if areas, _ := cmd.Flags().GetStringSlice("check-area"); len(areas) == 0 {
		_ = cmd.Flags().Set("check-area", "timeline")
	}
}

// IsCheckForced returns true if the --force flag is set on the command.
func IsCheckForced(cmd *cobra.Command) bool {
	forced, _ := cmd.Flags().GetBool("force")
	return forced
}
