// =============================================================================
// FILE: internal/cli/mutators/write.go
// PURPOSE: Write-related argument mutations. Normalises output paths and
//          file format strings before command execution.
// =============================================================================

package mutators

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/callbacks"
)

// MutateWriteArgs normalises output directory and file format flags.
func MutateWriteArgs(cmd *cobra.Command) {
	// Normalise save-dir to an absolute path.
	if dir, _ := cmd.Flags().GetString("save-dir"); dir != "" {
		_ = cmd.Flags().Set("save-dir", callbacks.NormalizePath(dir))
	}

	// Normalise database path.
	if db, _ := cmd.Flags().GetString("database"); db != "" {
		_ = cmd.Flags().Set("database", callbacks.NormalizePath(db))
	}
}
