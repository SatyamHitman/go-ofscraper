// =============================================================================
// FILE: internal/cli/bundles/msg_check.go
// PURPOSE: MsgCheckBundle registers message-specific check flags.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterMsgCheckBundle registers flags for the message check command.
func RegisterMsgCheckBundle(cmd *cobra.Command) {
	RegisterCheckBundle(cmd)
	flags.RegisterMediaFilterFlags(cmd)
	flags.RegisterPostFilterFlags(cmd)
}
