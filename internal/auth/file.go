// =============================================================================
// FILE: internal/auth/file.go
// PURPOSE: Auth file read/write operations. Handles loading and saving auth
//          credentials from/to the JSON auth file on disk. Ports Python
//          utils/auth/file.py.
// =============================================================================

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// Auth file I/O
// ---------------------------------------------------------------------------

// ReadAuthFile reads auth credentials from a JSON file.
//
// Parameters:
//   - path: Absolute path to the auth JSON file.
//
// Returns:
//   - The parsed Data, and any error.
func ReadAuthFile(path string) (*Data, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth file %s: %w", path, err)
	}

	// Try parsing as a direct Data object.
	var data Data
	if err := json.Unmarshal(raw, &data); err != nil {
		// Try parsing as {"auth": {...}} wrapper.
		var wrapper struct {
			Auth Data `json:"auth"`
		}
		if err2 := json.Unmarshal(raw, &wrapper); err2 != nil {
			return nil, fmt.Errorf("failed to parse auth file %s: %w", path, err)
		}
		data = wrapper.Auth
	}

	return &data, nil
}

// WriteAuthFile saves auth credentials to a JSON file. Creates parent
// directories if needed.
//
// Parameters:
//   - path: Absolute path to write.
//   - data: The auth credentials to save.
//
// Returns:
//   - Error if the write fails.
func WriteAuthFile(path string, data *Data) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialise auth data: %w", err)
	}

	if err := os.WriteFile(path, raw, 0600); err != nil {
		return fmt.Errorf("failed to write auth file %s: %w", path, err)
	}

	return nil
}

// EnsureAuthExists creates the auth file with empty values if it doesn't exist.
//
// Parameters:
//   - path: The auth file path.
//
// Returns:
//   - true if the file was created, false if it existed, and any error.
func EnsureAuthExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return false, nil
	}

	data := &Data{}
	if err := WriteAuthFile(path, data); err != nil {
		return false, err
	}
	return true, nil
}
