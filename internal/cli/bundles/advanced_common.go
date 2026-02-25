// =============================================================================
// FILE: internal/cli/bundles/advanced_common.go
// PURPOSE: AdvancedCommonBundle registers common + advanced flags.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterAdvancedCommonBundle registers common flags plus all advanced flag
// groups onto the given command.
func RegisterAdvancedCommonBundle(cmd *cobra.Command) {
	RegisterCommonBundle(cmd)
	flags.RegisterAdvancedProgramFlags(cmd)
	flags.RegisterAdvancedProcessingFlags(cmd)
	flags.RegisterAdvancedUserFilterFlags(cmd)
}
