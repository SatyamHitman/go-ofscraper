// =============================================================================
// FILE: internal/cli/flags/user_sort.go
// PURPOSE: User sort flag definitions: sort-type, descending-sort.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterUserSortFlags adds user sorting flags to the given command.
func RegisterUserSortFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("sort-type", "name", "Sort users by field (name, subscribed, expired)")
	f.Bool("descending-sort", false, "Sort users in descending order")
}
