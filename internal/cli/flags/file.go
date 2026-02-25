// =============================================================================
// FILE: internal/cli/flags/file.go
// PURPOSE: File flag definitions: file-format, text-type, space-replacer,
//          text-length.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterFileFlags adds file-related flags to the given command.
func RegisterFileFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.String("file-format", "{model_username}/{responsetype}/{value}/{mediatype}", "Template for output file paths")
	f.String("text-type", "letter", "Text filename type (letter, word)")
	f.String("space-replacer", "_", "Character to replace spaces in filenames")
	f.Int("text-length", 0, "Maximum text length in filenames (0 = unlimited)")
}
