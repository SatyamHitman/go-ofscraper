// =============================================================================
// FILE: internal/model/label.go
// PURPOSE: Defines the Label domain model representing a user-created label
//          (category/tag) that groups posts together. Ports Python
//          classes/labels.py with all fields and post associations.
// =============================================================================

package model

// Label represents a creator's content label (tag/category) that groups
// related posts together. Each label has a unique ID, name, type, and a
// collection of associated posts.
type Label struct {
	// LabelID is the unique identifier for this label.
	LabelID int64 `json:"id"`

	// Name is the human-readable label name.
	Name string `json:"name"`

	// Type classifies the label (e.g., custom, system).
	Type string `json:"type"`

	// Posts contains all posts associated with this label.
	Posts []*Post `json:"posts,omitempty"`

	// ModelID is the creator's user ID who owns this label.
	ModelID int64 `json:"model_id"`

	// Username is the creator's username who owns this label.
	Username string `json:"username"`
}

// PostCount returns the number of posts associated with this label.
//
// Returns:
//   - The count of associated posts.
func (l *Label) PostCount() int {
	return len(l.Posts)
}

// PostIDs returns a slice of all post IDs associated with this label.
// Useful for batch database operations.
//
// Returns:
//   - Slice of post ID integers.
func (l *Label) PostIDs() []int64 {
	ids := make([]int64, len(l.Posts))
	for i, p := range l.Posts {
		ids[i] = p.ID
	}
	return ids
}
