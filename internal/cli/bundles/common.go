// =============================================================================
// FILE: internal/cli/bundles/common.go
// PURPOSE: CommonBundle registers program + logging flags onto a command.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterCommonBundle registers the common set of flags (program + logging)
// onto the given command.
func RegisterCommonBundle(cmd *cobra.Command) {
	flags.RegisterProgramFlags(cmd)
	flags.RegisterLoggingFlags(cmd)
}
