// =============================================================================
// FILE: internal/cli/accessors/time.go
// PURPOSE: Read date/time flag values from cobra commands.
// =============================================================================

package accessors

import (
	"time"

	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// GetAfterDate returns the parsed "after" date flag value.
// Returns the zero time if the flag is empty or unparseable.
func GetAfterDate(cmd *cobra.Command) time.Time {
	v, _ := cmd.Flags().GetString("after")
	return flags.ParseDateOrDefault(v, time.Time{})
}

// GetBeforeDate returns a time derived from the "before-epoch" flag.
// Returns the zero time if the flag is 0.
func GetBeforeDate(cmd *cobra.Command) time.Time {
	v, _ := cmd.Flags().GetInt64("before-epoch")
	return flags.ParseEpoch(v)
}

// GetDateRange returns the after and before dates as a pair.
func GetDateRange(cmd *cobra.Command) (after, before time.Time) {
	return GetAfterDate(cmd), GetBeforeDate(cmd)
}
