// =============================================================================
// FILE: internal/cli/accessors/areas.go
// PURPOSE: Read area and post-area flag values from cobra commands.
// =============================================================================

package accessors

import (
	"github.com/spf13/cobra"
)

// GetAreas returns the check-area flag values from the command.
func GetAreas(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("check-area")
	return v
}

// GetPostsAreas returns the posts/content-area flag values from the command.
func GetPostsAreas(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("posts")
	return v
}
