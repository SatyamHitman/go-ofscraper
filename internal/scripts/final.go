// =============================================================================
// FILE: internal/scripts/final.go
// PURPOSE: Final script. Runs a user-defined script after all operations are
//          complete, used for cleanup or notification.
//          Ports Python scripts/final_script.py.
// =============================================================================

package scripts

import (
	"context"
)

// ---------------------------------------------------------------------------
// Final script
// ---------------------------------------------------------------------------

// RunFinal executes the final cleanup script.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - Error if the script fails.
func (m *Manager) RunFinal(ctx context.Context) error {
	if m.cfg.Final == "" {
		return nil
	}

	_, err := m.runScript(ctx, m.cfg.Final, nil, nil)
	return err
}
