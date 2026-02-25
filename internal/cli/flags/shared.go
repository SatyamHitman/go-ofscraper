// =============================================================================
// FILE: internal/cli/flags/shared.go
// PURPOSE: Shared flag registration helpers used across multiple flag groups.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterPostsAreaFlags adds the posts/area selection flag shared by several
// commands.
func RegisterPostsAreaFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("posts", "o", nil, "Content areas to process (timeline, messages, archived, etc.)")
}

// RegisterAllFlags is a convenience function that registers every flag group
// on the given command. Useful for the main scraper command.
func RegisterAllFlags(cmd *cobra.Command) {
	RegisterProgramFlags(cmd)
	RegisterLoggingFlags(cmd)
	RegisterDownloadFlags(cmd)
	RegisterMediaFilterFlags(cmd)
	RegisterPostFilterFlags(cmd)
	RegisterFileFlags(cmd)
	RegisterAutomaticFlags(cmd)
	RegisterUserSelectFlags(cmd)
	RegisterUserListFlags(cmd)
	RegisterUserSortFlags(cmd)
	RegisterAdvancedUserFilterFlags(cmd)
	RegisterAdvancedProcessingFlags(cmd)
	RegisterAdvancedProgramFlags(cmd)
	RegisterScriptFlags(cmd)
	RegisterMetadataFilterFlags(cmd)
	RegisterPostsAreaFlags(cmd)
}
