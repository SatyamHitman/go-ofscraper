// =============================================================================
// FILE: internal/cli/bundles/metadata.go
// PURPOSE: MetadataBundle registers flags for the metadata command.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterMetadataBundle registers flags for the metadata command.
func RegisterMetadataBundle(cmd *cobra.Command) {
	RegisterCommonBundle(cmd)
	flags.RegisterMetadataFilterFlags(cmd)
	flags.RegisterUserSelectFlags(cmd)
	flags.RegisterUserListFlags(cmd)
}
