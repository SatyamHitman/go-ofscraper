// =============================================================================
// FILE: internal/prompts/strings.go
// PURPOSE: Prompt display strings and labels. Centralizes all user-facing text
//          for interactive prompts. Ports Python prompts/prompt_strings.py.
// =============================================================================

package prompts

// ---------------------------------------------------------------------------
// Menu labels
// ---------------------------------------------------------------------------

const (
	// Main menu
	LabelMainMenu   = "Main Menu"
	LabelAction     = "Select Action"
	LabelAreas      = "Select Content Areas"
	LabelUsers      = "Select Users"
	LabelProfile    = "Select Profile"
	LabelConfig     = "Edit Configuration"

	// Actions
	ActionDownload   = "Download"
	ActionLike       = "Like"
	ActionUnlike     = "Unlike"
	ActionMetadata   = "Metadata"

	// Areas
	AreaTimeline     = "Timeline"
	AreaMessages     = "Messages"
	AreaArchived     = "Archived"
	AreaHighlights   = "Highlights"
	AreaStories      = "Stories"
	AreaPinned       = "Pinned"
	AreaStreams       = "Streams"
	AreaLabels       = "Labels"
	AreaPurchased    = "Purchased"

	// Misc
	LabelAll         = "All"
	LabelNone        = "None"
	LabelBack        = "Back"
	LabelExit        = "Exit"
	LabelContinue    = "Continue"
	LabelCancel      = "Cancel"
	LabelConfirm     = "Confirm"
	LabelYes         = "Yes"
	LabelNo          = "No"
)

// ---------------------------------------------------------------------------
// Action descriptions
// ---------------------------------------------------------------------------

// ActionDescriptions maps action names to their descriptions.
var ActionDescriptions = map[string]string{
	"download": "Download media files from selected creators",
	"like":     "Like/favorite posts from selected creators",
	"unlike":   "Remove likes from posts by selected creators",
	"metadata": "Update metadata for downloaded media",
}

// ---------------------------------------------------------------------------
// Area descriptions
// ---------------------------------------------------------------------------

// AreaDescriptions maps area names to their descriptions.
var AreaDescriptions = map[string]string{
	"timeline":   "Regular timeline posts",
	"messages":   "Direct messages",
	"archived":   "Archived posts",
	"highlights": "Profile highlights",
	"stories":    "Stories (24-hour posts)",
	"pinned":     "Pinned posts",
	"streams":    "Live stream recordings",
	"labels":     "Labelled/tagged posts",
	"purchased":  "Purchased content",
}

// AllAreas returns all available content area names.
func AllAreas() []string {
	return []string{
		"timeline", "messages", "archived", "highlights",
		"stories", "pinned", "streams", "labels", "purchased",
	}
}

// AllActions returns all available action names.
func AllActions() []string {
	return []string{"download", "like", "unlike", "metadata"}
}
