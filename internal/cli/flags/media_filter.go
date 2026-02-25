// =============================================================================
// FILE: internal/cli/flags/media_filter.go
// PURPOSE: Media filter flag definitions: quality, media-type, size-max,
//          size-min, media-id, length-max, length-min, max-count, media-sort,
//          media-desc, excluded, excluded-quality, before-epoch, after,
//          remove-duplicates.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterMediaFilterFlags adds media filtering flags to the given command.
func RegisterMediaFilterFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("quality", "source", "Preferred media quality (source, high, medium, low)")
	f.StringSlice("media-type", nil, "Media types to include (images, videos, audio)")
	f.Int64("size-max", 0, "Maximum file size in bytes (0 = no limit)")
	f.Int64("size-min", 0, "Minimum file size in bytes")
	f.StringSlice("media-id", nil, "Specific media IDs to download")
	f.Int("length-max", 0, "Maximum media duration in seconds (0 = no limit)")
	f.Int("length-min", 0, "Minimum media duration in seconds")
	f.Int("max-count", 0, "Maximum number of media items to process (0 = unlimited)")
	f.String("media-sort", "date", "Sort media by field (date, size, type)")
	f.Bool("media-desc", false, "Sort media in descending order")
	f.StringSlice("excluded", nil, "Media IDs to exclude")
	f.StringSlice("excluded-quality", nil, "Quality levels to exclude")
	f.Int64("before-epoch", 0, "Only include media before this Unix epoch timestamp")
	f.String("after", "", "Only include media after this date (YYYY-MM-DD)")
	f.Bool("remove-duplicates", false, "Remove duplicate media entries")
}
