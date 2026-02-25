// =============================================================================
// FILE: internal/cli/root.go
// PURPOSE: Root cobra command. Defines the top-level CLI command, persistent
//          flags, and the command tree structure.
//          Ports Python utils/args/main.py.
// =============================================================================

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gofscraper/pkg/version"
)

// ---------------------------------------------------------------------------
// Root command
// ---------------------------------------------------------------------------

var rootCmd = &cobra.Command{
	Use:   "gofscraper",
	Short: "GoFScraper — OnlyFans content downloader",
	Long:  `GoFScraper is a command-line tool for downloading and managing OnlyFans content.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags (global).
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "Profile name")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level (trace, debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("no-interactive", "n", false, "Disable interactive prompts")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// Version flag.
	rootCmd.Version = version.String()
}

// Root returns the root cobra command for adding sub-commands.
func Root() *cobra.Command {
	return rootCmd
}
