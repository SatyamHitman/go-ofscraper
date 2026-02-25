// =============================================================================
// FILE: internal/cli/bundles/check.go
// PURPOSE: CheckBundle registers flags for check commands.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterCheckBundle registers common flags plus check-specific flags.
func RegisterCheckBundle(cmd *cobra.Command) {
	RegisterCommonBundle(cmd)
	flags.RegisterCheckFlags(cmd)
	flags.RegisterUserSelectFlags(cmd)
	flags.RegisterUserListFlags(cmd)
	flags.RegisterUserSortFlags(cmd)
	flags.RegisterAdvancedUserFilterFlags(cmd)
}
