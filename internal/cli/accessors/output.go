// =============================================================================
// FILE: internal/cli/accessors/output.go
// PURPOSE: Read output-related flag values from cobra commands.
// =============================================================================

package accessors

import (
	"github.com/spf13/cobra"
)

// GetOutputDir returns the save-dir flag value.
func GetOutputDir(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("save-dir")
	return v
}

// GetFileFormat returns the file-format flag value.
func GetFileFormat(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("file-format")
	return v
}
