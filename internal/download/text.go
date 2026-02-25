// =============================================================================
// FILE: internal/download/text.go
// PURPOSE: Text file downloads. Saves post text/captions as .txt files
//          alongside media files. Ports Python utils/text.py.
// =============================================================================

package download

import (
	"fmt"
	"os"
	"path/filepath"

	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Text download
// ---------------------------------------------------------------------------

// SavePostText writes the post text to a .txt file.
//
// Parameters:
//   - post: The post whose text to save.
//   - outputDir: Directory to save the text file in.
//   - filename: The text filename (without extension).
//
// Returns:
//   - Error if the write fails.
func SavePostText(post *model.Post, outputDir, filename string) error {
	if post.RawText == "" {
		return nil // Nothing to save.
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create text dir: %w", err)
	}

	path := filepath.Join(outputDir, filename+".txt")

	return os.WriteFile(path, []byte(post.RawText), 0o644)
}
