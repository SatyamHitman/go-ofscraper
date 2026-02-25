// =============================================================================
// FILE: internal/cli/story_check.go
// PURPOSE: Story check subcommand. Lists and inspects stories.
//          Ports Python parse/commands/story.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var storyCheckCmd = &cobra.Command{
	Use:   "story_check",
	Short: "Check and list stories",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Story check mode...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storyCheckCmd)
}
