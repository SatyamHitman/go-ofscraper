// =============================================================================
// FILE: internal/auth/providers/init.go
// PURPOSE: Provider initialisation. Registers all known signature providers
//          with their URLs and priorities. Called during application startup.
// =============================================================================

package providers

import (
	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

// Init registers all known signature providers. Should be called once during
// application startup after environment variables are loaded.
func Init() {
	// Register providers in priority order (lower = tried first).
	Register(NewGenericProvider("digitalcriminals", env.DigitalCriminalsURL(), 1))
	Register(NewGenericProvider("datawhores", env.DatawhoresURL(), 2))
	Register(NewGenericProvider("xagler", env.XaglerURL(), 3))
	Register(NewGenericProvider("rafa", env.RafaURL(), 4))
	Register(NewGenericProvider("deviint", env.DeviintURL(), 5))

	// Generic fallback from env override.
	if url := env.DynamicGenericURL(); url != "" {
		Register(NewGenericProvider("custom", url, 10))
	}
}
