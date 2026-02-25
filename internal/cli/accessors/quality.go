// =============================================================================
// FILE: internal/cli/accessors/quality.go
// PURPOSE: Read quality and media type flag values from cobra commands.
// =============================================================================

package accessors

import (
	"github.com/spf13/cobra"
)

// GetQuality returns the quality flag value.
func GetQuality(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("quality")
	return v
}

// GetMediaTypes returns the media-type flag values.
func GetMediaTypes(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("media-type")
	return v
}
