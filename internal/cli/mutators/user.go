// =============================================================================
// FILE: internal/cli/mutators/user.go
// PURPOSE: User-related argument mutations. Parses and cleans username inputs.
// =============================================================================

package mutators

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/callbacks"
)

// MutateUsers parses and normalises the username and excluded-users flag
// values on the command.
func MutateUsers(cmd *cobra.Command) {
	// Parse usernames.
	if raw, _ := cmd.Flags().GetStringSlice("usernames"); len(raw) > 0 {
		parsed := callbacks.ParseUsernames(raw)
		_ = cmd.Flags().Set("usernames", sliceToCSV(parsed))
	}

	// Parse excluded users.
	if raw, _ := cmd.Flags().GetStringSlice("excluded-users"); len(raw) > 0 {
		parsed := callbacks.ParseUsernames(raw)
		_ = cmd.Flags().Set("excluded-users", sliceToCSV(parsed))
	}
}

// sliceToCSV joins a string slice with commas for use with cmd.Flags().Set().
func sliceToCSV(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += "," + s
	}
	return out
}
