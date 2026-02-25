// =============================================================================
// FILE: internal/db/transition.go
// PURPOSE: Schema migration for SQLite databases. Runs table creation and
//          ALTER TABLE statements to bring databases up to the current schema
//          version. Ports Python db/transition.py.
// =============================================================================

package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

// ---------------------------------------------------------------------------
// Schema version
// ---------------------------------------------------------------------------

// currentSchemaVersion is the latest schema version.
const currentSchemaVersion = 1

// ---------------------------------------------------------------------------
// Migration
// ---------------------------------------------------------------------------

// migrate runs all necessary schema migrations on the database.
//
// Parameters:
//   - db: The database connection.
//
// Returns:
//   - Error if any migration step fails.
func migrate(db *sql.DB) error {
	// Create schema_flags table if it doesn't exist (used to track version).
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_flags (
			flag_name  TEXT PRIMARY KEY,
			flag_value TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_flags table: %w", err)
	}

	version := getSchemaVersion(db)
	slog.Debug("database schema version", "current", version, "target", currentSchemaVersion)

	if version < 1 {
		if err := migrateV1(db); err != nil {
			return err
		}
		setSchemaVersion(db, 1)
	}

	return nil
}

// ---------------------------------------------------------------------------
// V1 migration: Create all tables
// ---------------------------------------------------------------------------

func migrateV1(db *sql.DB) error {
	statements := []string{
		// Posts table.
		`CREATE TABLE IF NOT EXISTS posts (
			id         INTEGER PRIMARY KEY,
			post_id    INTEGER NOT NULL,
			text       TEXT,
			price      REAL DEFAULT 0,
			paid       INTEGER DEFAULT 0,
			archived   INTEGER DEFAULT 0,
			created_at TEXT,
			model_id   INTEGER,
			UNIQUE(post_id)
		)`,

		// Messages table.
		`CREATE TABLE IF NOT EXISTS messages (
			id         INTEGER PRIMARY KEY,
			post_id    INTEGER NOT NULL,
			text       TEXT,
			price      REAL DEFAULT 0,
			paid       INTEGER DEFAULT 0,
			archived   INTEGER DEFAULT 0,
			created_at TEXT,
			model_id   INTEGER,
			UNIQUE(post_id)
		)`,

		// Media table.
		`CREATE TABLE IF NOT EXISTS medias (
			id            INTEGER PRIMARY KEY,
			media_id      INTEGER NOT NULL,
			post_id       INTEGER NOT NULL,
			link          TEXT,
			directory     TEXT,
			filename      TEXT,
			size          INTEGER DEFAULT 0,
			api_type      TEXT,
			media_type    TEXT,
			preview       INTEGER DEFAULT 0,
			linked        TEXT,
			downloaded    INTEGER DEFAULT 0,
			created_at    TEXT,
			posted_at     TEXT,
			hash          TEXT,
			model_id      INTEGER,
			UNIQUE(media_id)
		)`,

		// Stories table.
		`CREATE TABLE IF NOT EXISTS stories (
			id         INTEGER PRIMARY KEY,
			post_id    INTEGER NOT NULL,
			text       TEXT,
			price      REAL DEFAULT 0,
			paid       INTEGER DEFAULT 0,
			archived   INTEGER DEFAULT 0,
			created_at TEXT,
			model_id   INTEGER,
			UNIQUE(post_id)
		)`,

		// Labels table.
		`CREATE TABLE IF NOT EXISTS labels (
			id         INTEGER PRIMARY KEY,
			label_id   INTEGER NOT NULL,
			name       TEXT,
			type       TEXT,
			post_id    INTEGER,
			model_id   INTEGER,
			UNIQUE(label_id, post_id)
		)`,

		// Others table.
		`CREATE TABLE IF NOT EXISTS others (
			id         INTEGER PRIMARY KEY,
			post_id    INTEGER NOT NULL,
			text       TEXT,
			price      REAL DEFAULT 0,
			paid       INTEGER DEFAULT 0,
			archived   INTEGER DEFAULT 0,
			created_at TEXT,
			model_id   INTEGER,
			UNIQUE(post_id)
		)`,

		// Products table.
		`CREATE TABLE IF NOT EXISTS products (
			id         INTEGER PRIMARY KEY,
			post_id    INTEGER NOT NULL,
			text       TEXT,
			price      REAL DEFAULT 0,
			paid       INTEGER DEFAULT 0,
			archived   INTEGER DEFAULT 0,
			created_at TEXT,
			model_id   INTEGER,
			UNIQUE(post_id)
		)`,

		// Profiles table.
		`CREATE TABLE IF NOT EXISTS profiles (
			id         INTEGER PRIMARY KEY,
			user_id    INTEGER NOT NULL,
			username   TEXT,
			UNIQUE(user_id)
		)`,

		// Models table.
		`CREATE TABLE IF NOT EXISTS models (
			id         INTEGER PRIMARY KEY,
			model_id   INTEGER NOT NULL,
			username   TEXT,
			UNIQUE(model_id)
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration v1 failed: %w\nStatement: %s", err, stmt)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Schema version helpers
// ---------------------------------------------------------------------------

func getSchemaVersion(db *sql.DB) int {
	var value string
	err := db.QueryRow(
		`SELECT flag_value FROM schema_flags WHERE flag_name = 'schema_version'`,
	).Scan(&value)
	if err != nil {
		return 0
	}
	var v int
	fmt.Sscanf(value, "%d", &v)
	return v
}

func setSchemaVersion(db *sql.DB, version int) {
	_, _ = db.Exec(
		`INSERT OR REPLACE INTO schema_flags (flag_name, flag_value) VALUES ('schema_version', ?)`,
		fmt.Sprintf("%d", version),
	)
}
