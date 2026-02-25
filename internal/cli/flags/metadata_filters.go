// =============================================================================
// FILE: internal/cli/flags/metadata_filters.go
// PURPOSE: Metadata filter flag definitions: update-cache-list, download-list,
//          protected-filter, protected-bypass, cache-filter, cache-bypass,
//          preview.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterMetadataFilterFlags adds metadata filtering flags to the given command.
func RegisterMetadataFilterFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSlice("update-cache-list", nil, "Cache lists to update")
	f.StringSlice("download-list", nil, "Specific download lists to process")
	f.String("protected-filter", "", "Filter for protected content (include, exclude, only)")
	f.Bool("protected-bypass", false, "Bypass protected content restrictions")
	f.String("cache-filter", "", "Filter for cached content (include, exclude, only)")
	f.Bool("cache-bypass", false, "Bypass cache filtering restrictions")
	f.Bool("preview", false, "Preview mode: show what would be processed without acting")
}
