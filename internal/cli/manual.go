// =============================================================================
// FILE: internal/cli/manual.go
// PURPOSE: Manual subcommand. Downloads content from direct URLs.
//          Ports Python parse/commands/manual.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var manualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Download from direct URLs",
	Long:  `Downloads content from specific URLs provided as arguments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Manual download mode...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(manualCmd)
	manualCmd.Flags().StringSlice("url", nil, "URLs to download")
}
