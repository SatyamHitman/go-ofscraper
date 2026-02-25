// =============================================================================
// FILE: internal/cli/mutators/before.go
// PURPOSE: Pre-processing mutators that run before command execution.
//          Normalises and validates flag values.
// =============================================================================

package mutators

import (
	"github.com/spf13/cobra"

	"gofscraper/internal/cli/callbacks"
)

// RunBeforeMutators applies all pre-processing mutations to the command's
// parsed flags. Intended to be called from PersistentPreRunE or PreRunE.
func RunBeforeMutators(cmd *cobra.Command) error {
	if err := mutateActions(cmd); err != nil {
		return err
	}
	if err := mutateAreas(cmd); err != nil {
		return err
	}
	MutateUsers(cmd)
	return nil
}

// mutateActions validates and normalises the action flag values.
func mutateActions(cmd *cobra.Command) error {
	actions, _ := cmd.Flags().GetStringSlice("action")
	for _, a := range actions {
		if err := callbacks.ValidateAction(a); err != nil {
			return err
		}
	}
	return nil
}

// mutateAreas validates and normalises the posts area flag values.
func mutateAreas(cmd *cobra.Command) error {
	areas, _ := cmd.Flags().GetStringSlice("posts")
	if len(areas) == 0 {
		return nil
	}
	areas = callbacks.NormalizePostAreas(areas)
	if err := callbacks.ValidatePostAreas(areas); err != nil {
		return err
	}
	return nil
}
