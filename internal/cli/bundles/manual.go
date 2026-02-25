// =============================================================================
// FILE: internal/cli/bundles/manual.go
// PURPOSE: ManualBundle registers flags for the manual download command.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterManualBundle registers flags for the manual download command.
func RegisterManualBundle(cmd *cobra.Command) {
	RegisterCommonBundle(cmd)
	flags.RegisterDownloadFlags(cmd)
	flags.RegisterFileFlags(cmd)
	flags.RegisterMediaFilterFlags(cmd)
}
