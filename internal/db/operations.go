// =============================================================================
// FILE: internal/db/operations.go
// PURPOSE: Main database operations coordinator. Provides high-level CRUD
//          functions for posts, messages, media, stories, labels, and other
//          content types. Wraps raw SQL queries with proper error handling
//          and transaction management. Ports Python db/operations.py.
// =============================================================================

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Post operations
// ---------------------------------------------------------------------------

// UpsertPost inserts or updates a post record.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - conn: The model's database connection.
//   - postID: The OF post ID.
//   - text: Post text content.
//   - price: Post price.
//   - paid: Whether the post is paid.
//   - archived: Whether the post is archived.
//   - createdAt: Post creation timestamp.
//   - modelID: The model's numeric ID.
//
// Returns:
//   - Error if the upsert fails.
func UpsertPost(ctx context.Context, conn *Conn, postID int64, text string, price float64, paid, archived bool, createdAt string, modelID int64) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO posts (post_id, text, price, paid, archived, created_at, model_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(post_id) DO UPDATE SET
		   text = excluded.text,
		   price = excluded.price,
		   paid = excluded.paid,
		   archived = excluded.archived,
		   created_at = excluded.created_at`,
		postID, text, price, boolToInt(paid), boolToInt(archived), createdAt, modelID,
	)
	return err
}

// GetPost retrieves a single post by ID.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - conn: Database connection.
//   - postID: The post ID to look up.
//
// Returns:
//   - A PostRow, and any error (sql.ErrNoRows if not found).
func GetPost(ctx context.Context, conn *Conn, postID int64) (*PostRow, error) {
	row := conn.DB.QueryRowContext(ctx,
		`SELECT post_id, text, price, paid, archived, created_at, model_id FROM posts WHERE post_id = ?`,
		postID,
	)

	p := &PostRow{}
	err := row.Scan(&p.PostID, &p.Text, &p.Price, &p.Paid, &p.Archived, &p.CreatedAt, &p.ModelID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetAllPosts retrieves all posts for the model.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//
// Returns:
//   - Slice of PostRow, and any error.
func GetAllPosts(ctx context.Context, conn *Conn) ([]PostRow, error) {
	rows, err := conn.DB.QueryContext(ctx,
		`SELECT post_id, text, price, paid, archived, created_at, model_id FROM posts ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostRow
	for rows.Next() {
		var p PostRow
		if err := rows.Scan(&p.PostID, &p.Text, &p.Price, &p.Paid, &p.Archived, &p.CreatedAt, &p.ModelID); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// ---------------------------------------------------------------------------
// Media operations
// ---------------------------------------------------------------------------

// UpsertMedia inserts or updates a media record.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//   - m: The media data to upsert.
//
// Returns:
//   - Error if the upsert fails.
func UpsertMedia(ctx context.Context, conn *Conn, m MediaRow) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO medias (media_id, post_id, link, directory, filename, size, api_type, media_type, preview, linked, downloaded, created_at, posted_at, hash, model_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(media_id) DO UPDATE SET
		   link = excluded.link,
		   directory = excluded.directory,
		   filename = excluded.filename,
		   size = excluded.size,
		   api_type = excluded.api_type,
		   media_type = excluded.media_type,
		   preview = excluded.preview,
		   linked = excluded.linked,
		   downloaded = excluded.downloaded,
		   posted_at = excluded.posted_at,
		   hash = excluded.hash`,
		m.MediaID, m.PostID, m.Link, m.Directory, m.Filename, m.Size,
		m.APIType, m.MediaType, boolToInt(m.Preview), m.Linked,
		boolToInt(m.Downloaded), m.CreatedAt, m.PostedAt, m.Hash, m.ModelID,
	)
	return err
}

// GetMediaByPostID retrieves all media for a given post.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//   - postID: The post ID.
//
// Returns:
//   - Slice of MediaRow, and any error.
func GetMediaByPostID(ctx context.Context, conn *Conn, postID int64) ([]MediaRow, error) {
	rows, err := conn.DB.QueryContext(ctx,
		`SELECT media_id, post_id, link, directory, filename, size, api_type, media_type, preview, linked, downloaded, created_at, posted_at, hash, model_id
		 FROM medias WHERE post_id = ?`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMediaRows(rows)
}

// GetAllMedia retrieves all media records.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//
// Returns:
//   - Slice of MediaRow, and any error.
func GetAllMedia(ctx context.Context, conn *Conn) ([]MediaRow, error) {
	rows, err := conn.DB.QueryContext(ctx,
		`SELECT media_id, post_id, link, directory, filename, size, api_type, media_type, preview, linked, downloaded, created_at, posted_at, hash, model_id
		 FROM medias ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMediaRows(rows)
}

// GetDownloadedMedia retrieves only media that have been downloaded.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//
// Returns:
//   - Slice of MediaRow.
func GetDownloadedMedia(ctx context.Context, conn *Conn) ([]MediaRow, error) {
	rows, err := conn.DB.QueryContext(ctx,
		`SELECT media_id, post_id, link, directory, filename, size, api_type, media_type, preview, linked, downloaded, created_at, posted_at, hash, model_id
		 FROM medias WHERE downloaded = 1`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMediaRows(rows)
}

// ---------------------------------------------------------------------------
// Message operations
// ---------------------------------------------------------------------------

// UpsertMessage inserts or updates a message record.
func UpsertMessage(ctx context.Context, conn *Conn, postID int64, text string, price float64, paid, archived bool, createdAt string, modelID int64) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO messages (post_id, text, price, paid, archived, created_at, model_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(post_id) DO UPDATE SET
		   text = excluded.text,
		   price = excluded.price,
		   paid = excluded.paid,
		   archived = excluded.archived,
		   created_at = excluded.created_at`,
		postID, text, price, boolToInt(paid), boolToInt(archived), createdAt, modelID,
	)
	return err
}

// ---------------------------------------------------------------------------
// Story operations
// ---------------------------------------------------------------------------

// UpsertStory inserts or updates a story record.
func UpsertStory(ctx context.Context, conn *Conn, postID int64, text string, price float64, paid, archived bool, createdAt string, modelID int64) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO stories (post_id, text, price, paid, archived, created_at, model_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(post_id) DO UPDATE SET
		   text = excluded.text,
		   price = excluded.price,
		   paid = excluded.paid,
		   archived = excluded.archived,
		   created_at = excluded.created_at`,
		postID, text, price, boolToInt(paid), boolToInt(archived), createdAt, modelID,
	)
	return err
}

// ---------------------------------------------------------------------------
// Label operations
// ---------------------------------------------------------------------------

// UpsertLabel inserts or updates a label record.
func UpsertLabel(ctx context.Context, conn *Conn, labelID int64, name, labelType string, postID, modelID int64) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO labels (label_id, name, type, post_id, model_id)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(label_id, post_id) DO UPDATE SET
		   name = excluded.name,
		   type = excluded.type`,
		labelID, name, labelType, postID, modelID,
	)
	return err
}

// ---------------------------------------------------------------------------
// Profile operations
// ---------------------------------------------------------------------------

// UpsertProfile inserts or updates a profile record.
func UpsertProfile(ctx context.Context, conn *Conn, userID int64, username string) error {
	_, err := conn.DB.ExecContext(ctx,
		`INSERT INTO profiles (user_id, username)
		 VALUES (?, ?)
		 ON CONFLICT(user_id) DO UPDATE SET username = excluded.username`,
		userID, username,
	)
	return err
}

// ---------------------------------------------------------------------------
// Statistics
// ---------------------------------------------------------------------------

// Stats holds aggregate statistics for a model's database.
type Stats struct {
	PostCount    int
	MessageCount int
	MediaCount   int
	StoryCount   int
	LabelCount   int
	Downloaded   int
	TotalSize    int64
}

// GetStats retrieves aggregate statistics from the database.
//
// Parameters:
//   - ctx: Context.
//   - conn: Database connection.
//
// Returns:
//   - Stats struct, and any error.
func GetStats(ctx context.Context, conn *Conn) (Stats, error) {
	var s Stats

	queries := []struct {
		query string
		dest  *int
	}{
		{"SELECT COUNT(*) FROM posts", &s.PostCount},
		{"SELECT COUNT(*) FROM messages", &s.MessageCount},
		{"SELECT COUNT(*) FROM medias", &s.MediaCount},
		{"SELECT COUNT(*) FROM stories", &s.StoryCount},
		{"SELECT COUNT(*) FROM labels", &s.LabelCount},
		{"SELECT COUNT(*) FROM medias WHERE downloaded = 1", &s.Downloaded},
	}

	for _, q := range queries {
		if err := conn.DB.QueryRowContext(ctx, q.query).Scan(q.dest); err != nil {
			return s, fmt.Errorf("stats query failed: %w", err)
		}
	}

	// Total size of downloaded media.
	conn.DB.QueryRowContext(ctx,
		"SELECT COALESCE(SUM(size), 0) FROM medias WHERE downloaded = 1",
	).Scan(&s.TotalSize)

	return s, nil
}

// ---------------------------------------------------------------------------
// Transaction helper
// ---------------------------------------------------------------------------

// WithTx runs a function within a database transaction. The transaction is
// committed if fn returns nil, rolled back otherwise.
//
// Parameters:
//   - ctx: Context for the transaction.
//   - conn: Database connection.
//   - fn: Function to run within the transaction.
//
// Returns:
//   - Error from fn or from commit/rollback.
func WithTx(ctx context.Context, conn *Conn, fn func(tx *sql.Tx) error) error {
	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ---------------------------------------------------------------------------
// Row types
// ---------------------------------------------------------------------------

// PostRow represents a row in the posts/messages/stories/others tables.
type PostRow struct {
	PostID    int64
	Text      sql.NullString
	Price     float64
	Paid      int
	Archived  int
	CreatedAt sql.NullString
	ModelID   int64
}

// MediaRow represents a row in the medias table.
type MediaRow struct {
	MediaID    int64
	PostID     int64
	Link       sql.NullString
	Directory  sql.NullString
	Filename   sql.NullString
	Size       int64
	APIType    sql.NullString
	MediaType  sql.NullString
	Preview    bool
	Linked     sql.NullString
	Downloaded bool
	CreatedAt  sql.NullString
	PostedAt   sql.NullString
	Hash       sql.NullString
	ModelID    int64
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanMediaRows(rows *sql.Rows) ([]MediaRow, error) {
	var medias []MediaRow
	for rows.Next() {
		var m MediaRow
		var preview, downloaded int
		if err := rows.Scan(
			&m.MediaID, &m.PostID, &m.Link, &m.Directory, &m.Filename,
			&m.Size, &m.APIType, &m.MediaType, &preview, &m.Linked,
			&downloaded, &m.CreatedAt, &m.PostedAt, &m.Hash, &m.ModelID,
		); err != nil {
			return nil, err
		}
		m.Preview = preview == 1
		m.Downloaded = downloaded == 1
		medias = append(medias, m)
	}
	return medias, rows.Err()
}

// ensure time import is used
var _ = time.Now
