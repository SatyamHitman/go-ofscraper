// =============================================================================
// FILE: internal/config/profile.go
// PURPOSE: Profile management for multi-profile configurations. Handles
//          creating, switching, listing, and deleting named config profiles.
//          Ports Python utils/profiles/data.py, manage.py, tools.py.
// =============================================================================

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// Profile operations
// ---------------------------------------------------------------------------

// ProfileDir returns the directory path for a named profile.
//
// Parameters:
//   - name: The profile name.
//
// Returns:
//   - The absolute path to the profile directory.
func ProfileDir(name string) string {
	return filepath.Join(env.ConfigDir(), name)
}

// CurrentProfileDir returns the directory path for the active profile.
//
// Returns:
//   - The absolute path to the current profile directory.
func CurrentProfileDir() string {
	return ProfileDir(GetMainProfile())
}

// ListProfiles returns all available profile names by scanning the config
// directory for subdirectories.
//
// Returns:
//   - A string slice of profile names, and any error.
func ListProfiles() ([]string, error) {
	configDir := env.ConfigDir()
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{DefaultProfile}, nil
		}
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	var profiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			profiles = append(profiles, entry.Name())
		}
	}

	// Ensure default profile is always listed
	if len(profiles) == 0 {
		profiles = append(profiles, DefaultProfile)
	}

	return profiles, nil
}

// CreateProfile creates a new profile directory. Does not change the active
// profile — use SwitchProfile for that.
//
// Parameters:
//   - name: The profile name to create.
//
// Returns:
//   - Error if the directory cannot be created.
func CreateProfile(name string) error {
	dir := ProfileDir(name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory %s: %w", dir, err)
	}
	return nil
}

// SwitchProfile changes the active profile to the named profile. Updates the
// global config's MainProfile field and writes the config to disk.
//
// Parameters:
//   - name: The profile name to switch to.
//
// Returns:
//   - Error if the profile doesn't exist or the config write fails.
func SwitchProfile(name string) error {
	dir := ProfileDir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", name)
	}

	cfg := Get()
	cfg.MainProfile = name
	return Update(cfg)
}

// DeleteProfile removes a profile directory. Cannot delete the active profile
// or the default profile.
//
// Parameters:
//   - name: The profile name to delete.
//
// Returns:
//   - Error if the profile is active, is the default, or cannot be removed.
func DeleteProfile(name string) error {
	if name == GetMainProfile() {
		return fmt.Errorf("cannot delete the active profile %q", name)
	}
	if name == DefaultProfile {
		return fmt.Errorf("cannot delete the default profile %q", name)
	}

	dir := ProfileDir(name)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to delete profile directory %s: %w", dir, err)
	}

	return nil
}
