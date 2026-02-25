// =============================================================================
// FILE: internal/cli/flags/post_filter.go
// PURPOSE: Post filter flag definitions: only, date-after, label-after,
//          filter-text, neg-filter, before-after, timed, skip-pinned,
//          mass-only, mass-exclude, text-search, text-exclude.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterPostFilterFlags adds post filtering flags to the given command.
func RegisterPostFilterFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSlice("only", nil, "Only include these post types (free, paid)")
	f.String("date-after", "", "Only include posts after this date (YYYY-MM-DD)")
	f.String("label-after", "", "Only include posts with labels after this date")
	f.String("filter-text", "", "Only include posts matching this text pattern")
	f.String("neg-filter", "", "Exclude posts matching this text pattern")
	f.String("before-after", "", "Only include posts before this date (YYYY-MM-DD)")
	f.Bool("timed", false, "Only include timed/expiring posts")
	f.Bool("skip-pinned", false, "Skip pinned posts")
	f.StringSlice("mass-only", nil, "Only include these mass message types")
	f.StringSlice("mass-exclude", nil, "Exclude these mass message types")
	f.String("text-search", "", "Search for posts containing this text")
	f.String("text-exclude", "", "Exclude posts containing this text")
}
