// =============================================================================
// FILE: internal/cli/flags/download.go
// PURPOSE: Download flag definitions: arrow, database, save-dir, download-limit.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterDownloadFlags adds download-related flags to the given command.
func RegisterDownloadFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.Bool("arrow", false, "Use arrow (columnar) storage for metadata")
	f.String("database", "", "Path to the database file")
	f.String("save-dir", "", "Directory where downloaded content is saved")
	f.Int("download-limit", 0, "Maximum number of concurrent downloads (0 = unlimited)")
}
