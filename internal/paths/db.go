// =============================================================================
// FILE: internal/paths/db.go
// PURPOSE: Database path resolution. Provides functions for locating the
//          per-model SQLite database files based on the current profile and
//          save location settings. Ports Python utils/paths/db.py.
// =============================================================================

package paths

import (
	"path/filepath"
	"strings"

	"gofscraper/internal/config"
)

// ---------------------------------------------------------------------------
// Database paths
// ---------------------------------------------------------------------------

// DBDir returns the directory where the database for a given model is stored.
// This follows the same directory structure as the download location, but
// with the DB inside the model's directory.
//
// Parameters:
//   - modelUsername: The OF model username.
//
// Returns:
//   - Absolute path to the model's database directory.
func DBDir(modelUsername string) string {
	saveLocation := config.GetSaveLocation()
	return filepath.Join(saveLocation, sanitizeComponent(modelUsername))
}

// DBPath returns the full path to the SQLite database file for a model.
//
// Parameters:
//   - modelUsername: The OF model username.
//
// Returns:
//   - Absolute path to the .db file.
func DBPath(modelUsername string) string {
	return filepath.Join(DBDir(modelUsername), "user_data.db")
}

// BackupDBPath returns the path for a database backup file.
//
// Parameters:
//   - modelUsername: The OF model username.
//   - suffix: Backup suffix (e.g. timestamp string).
//
// Returns:
//   - Absolute path to the backup .db file.
func BackupDBPath(modelUsername string, suffix string) string {
	return filepath.Join(DBDir(modelUsername), "user_data_"+suffix+".db")
}

// AllDBPaths scans the save location for all model directories containing
// database files.
//
// Returns:
//   - A map of model username to DB file path, and any error.
func AllDBPaths() (map[string]string, error) {
	saveLocation := config.GetSaveLocation()
	result := make(map[string]string)

	entries, err := filepath.Glob(filepath.Join(saveLocation, "*", "user_data.db"))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		dir := filepath.Dir(entry)
		model := filepath.Base(dir)
		result[model] = entry
	}

	return result, nil
}

// sanitizeComponent removes characters unsafe for directory names.
func sanitizeComponent(s string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(s)
}
