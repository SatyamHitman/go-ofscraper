// =============================================================================
// FILE: internal/cli/flags/user_list.go
// PURPOSE: User list flag definitions: user-list, blacklist.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterUserListFlags adds user list flags to the given command.
func RegisterUserListFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSlice("user-list", nil, "Paths to files containing usernames (one per line)")
	f.StringSlice("blacklist", nil, "Paths to files containing blacklisted usernames")
}
