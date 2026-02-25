// =============================================================================
// FILE: internal/cli/flags/check.go
// PURPOSE: Check flag definitions: check-area, user, file, force,
//          table-progress, table-name.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterCheckFlags adds check-mode flags to the given command.
func RegisterCheckFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringSlice("check-area", nil, "Content areas to check (timeline, messages, archived, etc.)")
	f.StringSlice("user", nil, "Users to check")
	f.String("file", "", "Path to file for check input")
	f.Bool("force", false, "Force check even if recently completed")
	f.Bool("table-progress", true, "Show progress in table format")
	f.String("table-name", "", "Custom table name for progress display")
}
