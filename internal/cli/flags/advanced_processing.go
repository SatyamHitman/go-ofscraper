// =============================================================================
// FILE: internal/cli/flags/advanced_processing.go
// PURPOSE: Advanced processing flag definitions: user-filter, force-individual.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterAdvancedProcessingFlags adds advanced processing flags to the given
// command.
func RegisterAdvancedProcessingFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("user-filter", "", "Custom expression to filter users")
	f.Bool("force-individual", false, "Force individual user processing instead of batch")
}
