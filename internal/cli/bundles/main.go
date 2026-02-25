// =============================================================================
// FILE: internal/cli/bundles/main.go
// PURPOSE: MainBundle registers all flags for the main scraper command.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterMainBundle registers every flag group needed by the main scraper
// command.
func RegisterMainBundle(cmd *cobra.Command) {
	RegisterAdvancedCommonBundle(cmd)
	flags.RegisterDownloadFlags(cmd)
	flags.RegisterMediaFilterFlags(cmd)
	flags.RegisterPostFilterFlags(cmd)
	flags.RegisterFileFlags(cmd)
	flags.RegisterAutomaticFlags(cmd)
	flags.RegisterUserSelectFlags(cmd)
	flags.RegisterUserListFlags(cmd)
	flags.RegisterUserSortFlags(cmd)
	flags.RegisterScriptFlags(cmd)
	flags.RegisterMetadataFilterFlags(cmd)
	flags.RegisterPostsAreaFlags(cmd)
}
