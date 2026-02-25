// =============================================================================
// FILE: internal/cli/scraper.go
// PURPOSE: Scraper subcommand. The main command for downloading content from
//          OnlyFans. Ports Python parse/commands/main.py.
// =============================================================================

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Scraper command
// ---------------------------------------------------------------------------

var scraperCmd = &cobra.Command{
	Use:   "scraper",
	Short: "Run the scraper to download content",
	Long:  `Downloads media and text from OnlyFans creators based on configured settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting scraper...")
		// TODO: Wire to app.RunScraper()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scraperCmd)

	// Scraper-specific flags.
	scraperCmd.Flags().StringSliceP("action", "a", []string{"download"}, "Actions to perform (download, like, unlike)")
	scraperCmd.Flags().StringSliceP("posts", "o", nil, "Content areas (timeline, messages, archived, etc.)")
	scraperCmd.Flags().StringSliceP("users", "u", nil, "Usernames to process")
	scraperCmd.Flags().StringSlice("excluded-users", nil, "Usernames to exclude")
	scraperCmd.Flags().BoolP("daemon", "d", false, "Run in daemon mode")
}
