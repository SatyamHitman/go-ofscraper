// =============================================================================
// FILE: internal/cli/metadata.go
// PURPOSE: Metadata subcommand. Updates metadata for downloaded content.
//          Ports Python parse/commands/metadata.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Update metadata for downloaded content",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Metadata mode...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
