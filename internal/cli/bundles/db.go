// =============================================================================
// FILE: internal/cli/bundles/db.go
// PURPOSE: DBBundle registers flags for the database command.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterDBBundle registers flags for the database management command.
func RegisterDBBundle(cmd *cobra.Command) {
	RegisterCommonBundle(cmd)
	flags.RegisterDownloadFlags(cmd)
	flags.RegisterUserSelectFlags(cmd)
}
