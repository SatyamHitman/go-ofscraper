// =============================================================================
// FILE: internal/db/backup.go
// PURPOSE: Database backup operations. Creates timestamped copies of model
//          databases before destructive operations like merging or migration.
//          Ports Python db/backup.py.
// =============================================================================

package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ---------------------------------------------------------------------------
// Backup
// ---------------------------------------------------------------------------

// Backup creates a timestamped backup of the given database file.
//
// Parameters:
//   - dbPath: Path to the database file to back up.
//
// Returns:
//   - The backup file path, and any error.
func Backup(dbPath string) (string, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return "", fmt.Errorf("database file does not exist: %s", dbPath)
	}

	// Build backup filename with timestamp.
	dir := filepath.Dir(dbPath)
	base := filepath.Base(dbPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	ts := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(dir, fmt.Sprintf("%s_backup_%s%s", name, ts, ext))

	if err := copyFile(dbPath, backupPath); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	// Also back up WAL and SHM files if they exist.
	for _, suffix := range []string{"-wal", "-shm"} {
		src := dbPath + suffix
		if _, err := os.Stat(src); err == nil {
			dst := backupPath + suffix
			_ = copyFile(src, dst) // Best effort for WAL/SHM.
		}
	}

	return backupPath, nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
