// =============================================================================
// FILE: internal/cli/bundles/post_check.go
// PURPOSE: PostCheckBundle registers post-specific check flags.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterPostCheckBundle registers flags for the post check command.
func RegisterPostCheckBundle(cmd *cobra.Command) {
	RegisterCheckBundle(cmd)
	flags.RegisterMediaFilterFlags(cmd)
	flags.RegisterPostFilterFlags(cmd)
	flags.RegisterDownloadFlags(cmd)
}
