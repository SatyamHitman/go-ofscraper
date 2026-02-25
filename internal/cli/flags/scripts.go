// =============================================================================
// FILE: internal/cli/flags/scripts.go
// PURPOSE: Script flag definitions: after-action-script, post-script,
//          naming-script, after-dl-script, skip-dl-script.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterScriptFlags adds script hook flags to the given command.
func RegisterScriptFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("after-action-script", "", "Script to run after each action completes")
	f.String("post-script", "", "Script to run after all processing finishes")
	f.String("naming-script", "", "Script that generates custom file names")
	f.String("after-dl-script", "", "Script to run after each file download")
	f.String("skip-dl-script", "", "Script that decides whether to skip a download")
}
