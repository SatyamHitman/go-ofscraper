// =============================================================================
// FILE: internal/cli/accessors/read.go
// PURPOSE: Read various flag values from cobra commands: usernames,
//          excluded users, config, profile, and more.
// =============================================================================

package accessors

import (
	"github.com/spf13/cobra"
)

// GetUsernames returns the usernames flag values.
func GetUsernames(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("usernames")
	return v
}

// GetExcludedUsers returns the excluded-users flag values.
func GetExcludedUsers(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("excluded-users")
	return v
}

// GetConfig returns the config flag value from the command or its parents.
func GetConfig(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("config")
	if v == "" {
		// Try the persistent flag from parent.
		if p := cmd.Parent(); p != nil {
			v, _ = p.PersistentFlags().GetString("config")
		}
	}
	return v
}

// GetProfile returns the profile flag value from the command or its parents.
func GetProfile(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("profile")
	if v == "" {
		if p := cmd.Parent(); p != nil {
			v, _ = p.PersistentFlags().GetString("profile")
		}
	}
	return v
}

// GetUserList returns the user-list flag values.
func GetUserList(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("user-list")
	return v
}

// GetBlacklist returns the blacklist flag values.
func GetBlacklist(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("blacklist")
	return v
}

// GetSortType returns the sort-type flag value.
func GetSortType(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("sort-type")
	return v
}

// GetDescendingSort returns true if descending-sort is set.
func GetDescendingSort(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("descending-sort")
	return v
}

// GetConfigGroup returns the config-group flag value.
func GetConfigGroup(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("config-group")
	return v
}

// GetEnvFiles returns the env-files flag values.
func GetEnvFiles(cmd *cobra.Command) []string {
	v, _ := cmd.Flags().GetStringSlice("env-files")
	return v
}
