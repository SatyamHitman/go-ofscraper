// =============================================================================
// FILE: internal/cli/flags/user_select.go
// PURPOSE: User select flag definitions: usernames, excluded-users.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterUserSelectFlags adds user selection flags to the given command.
func RegisterUserSelectFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSliceP("usernames", "u", nil, "Usernames to process")
	f.StringSlice("excluded-users", nil, "Usernames to exclude from processing")
}
