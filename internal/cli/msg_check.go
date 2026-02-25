// =============================================================================
// FILE: internal/cli/msg_check.go
// PURPOSE: Message check subcommand. Lists and inspects messages.
//          Ports Python parse/commands/message.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var msgCheckCmd = &cobra.Command{
	Use:   "msg_check",
	Short: "Check and list messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Message check mode...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(msgCheckCmd)
}
