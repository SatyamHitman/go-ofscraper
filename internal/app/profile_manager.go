// =============================================================================
// FILE: internal/app/profile_manager.go
// PURPOSE: ProfileManager handles profile operations including listing,
//          creating, switching, and deleting configuration profiles.
//          Ports Python utils/profiles/manage.py.
// =============================================================================

package app

import (
	"fmt"
	"log/slog"

	"gofscraper/internal/config"
)

// ---------------------------------------------------------------------------
// ProfileManager
// ---------------------------------------------------------------------------

// ProfileManager provides profile management operations.
type ProfileManager struct {
	logger *slog.Logger
}

// NewProfileManager creates a ProfileManager.
//
// Parameters:
//   - logger: Structured logger.
//
// Returns:
//   - A configured ProfileManager.
func NewProfileManager(logger *slog.Logger) *ProfileManager {
	if logger == nil {
		logger = slog.Default()
	}
	return &ProfileManager{
		logger: logger,
	}
}

// List returns all available profile names.
//
// Returns:
//   - A slice of profile name strings, and any error.
func (pm *ProfileManager) List() ([]string, error) {
	profiles, err := config.ListProfiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	pm.logger.Info("available profiles",
		"count", len(profiles),
		"active", config.GetMainProfile(),
	)

	return profiles, nil
}

// Create creates a new configuration profile.
//
// Parameters:
//   - name: The name for the new profile.
//
// Returns:
//   - Error if the profile cannot be created.
func (pm *ProfileManager) Create(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if err := config.CreateProfile(name); err != nil {
		return fmt.Errorf("failed to create profile %q: %w", name, err)
	}

	pm.logger.Info("profile created", "name", name)
	return nil
}

// Switch changes the active profile.
//
// Parameters:
//   - name: The profile name to switch to.
//
// Returns:
//   - Error if the profile does not exist or the switch fails.
func (pm *ProfileManager) Switch(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	prev := config.GetMainProfile()
	if err := config.SwitchProfile(name); err != nil {
		return fmt.Errorf("failed to switch to profile %q: %w", name, err)
	}

	pm.logger.Info("profile switched",
		"from", prev,
		"to", name,
	)
	return nil
}

// Delete removes a configuration profile.
//
// Parameters:
//   - name: The profile name to delete.
//
// Returns:
//   - Error if the profile cannot be deleted.
func (pm *ProfileManager) Delete(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if err := config.DeleteProfile(name); err != nil {
		return fmt.Errorf("failed to delete profile %q: %w", name, err)
	}

	pm.logger.Info("profile deleted", "name", name)
	return nil
}

// Current returns the name of the active profile.
//
// Returns:
//   - The active profile name.
func (pm *ProfileManager) Current() string {
	return config.GetMainProfile()
}
