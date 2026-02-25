// =============================================================================
// FILE: internal/cli/bundles/story_check.go
// PURPOSE: StoryCheckBundle registers story-specific check flags.
// =============================================================================

package bundles

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/flags"
)

// RegisterStoryCheckBundle registers flags for the story check command.
func RegisterStoryCheckBundle(cmd *cobra.Command) {
	RegisterCheckBundle(cmd)
	flags.RegisterMediaFilterFlags(cmd)
}
