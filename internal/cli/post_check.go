// =============================================================================
// FILE: internal/cli/post_check.go
// PURPOSE: Post check subcommand. Lists and inspects posts.
//          Ports Python parse/commands/post.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var postCheckCmd = &cobra.Command{
	Use:   "post_check",
	Short: "Check and list posts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Post check mode...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCheckCmd)
}
