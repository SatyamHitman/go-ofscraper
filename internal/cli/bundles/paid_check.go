// =============================================================================
// FILE: internal/cli/bundles/paid_check.go
// PURPOSE: PaidCheckBundle registers paid-specific check flags.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterPaidCheckBundle registers flags for the paid check command.
func RegisterPaidCheckBundle(cmd *cobra.Command) {
	RegisterCheckBundle(cmd)
	flags.RegisterMediaFilterFlags(cmd)
	flags.RegisterPostFilterFlags(cmd)
}
