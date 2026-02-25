// =============================================================================
// FILE: internal/cli/flags/program.go
// PURPOSE: Program-level flag definitions: version-check, config-group, run,
//          env-files. Registers flags on cobra commands.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterProgramFlags adds program-level flags to the given command.
func RegisterProgramFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.Bool("version-check", true, "Check for new versions on startup")
	f.String("config-group", "", "Configuration group to use")
	f.StringSlice("run", nil, "Specific run targets")
	f.StringSlice("env-files", nil, "Paths to .env files to load")
}
