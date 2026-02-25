// =============================================================================
// FILE: internal/tui/utils/names.go
// PURPOSE: Column name constants and display name mappings. Maps internal
//          column keys to human-readable display names for the TUI table.
//          Ports Python classes/table/utils/names.py.
// =============================================================================

package tuiutils

// ---------------------------------------------------------------------------
// Column key constants
// ---------------------------------------------------------------------------

const (
	ColIndex        = "index"
	ColMediaID      = "media_id"
	ColPostID       = "post_id"
	ColUsername      = "username"
	ColDate         = "date"
	ColMediaType    = "media_type"
	ColFilename     = "filename"
	ColSize         = "size"
	ColDuration     = "duration"
	ColPrice        = "price"
	ColPaid         = "paid"
	ColText         = "text"
	ColResponseType = "response_type"
	ColDownloaded   = "downloaded"
	ColLabel        = "label"
	ColURL          = "url"
	ColProtected    = "protected"
	ColLinked       = "linked"
	ColStatus       = "status"
	ColProgress     = "progress"
	ColSpeed        = "speed"
	ColETA          = "eta"
	ColTask         = "task"
	ColElapsed      = "elapsed"
)

// ---------------------------------------------------------------------------
// Display name mapping
// ---------------------------------------------------------------------------

// DisplayNames maps internal column keys to human-readable display names.
var DisplayNames = map[string]string{
	ColIndex:        "#",
	ColMediaID:      "Media ID",
	ColPostID:       "Post ID",
	ColUsername:      "Username",
	ColDate:         "Date",
	ColMediaType:    "Type",
	ColFilename:     "Filename",
	ColSize:         "Size",
	ColDuration:     "Duration",
	ColPrice:        "Price",
	ColPaid:         "Paid",
	ColText:         "Text",
	ColResponseType: "Source",
	ColDownloaded:   "Downloaded",
	ColLabel:        "Label",
	ColURL:          "URL",
	ColProtected:    "Protected",
	ColLinked:       "Linked",
	ColStatus:       "Status",
	ColProgress:     "Progress",
	ColSpeed:        "Speed",
	ColETA:          "ETA",
	ColTask:         "Task",
	ColElapsed:      "Elapsed",
}

// DisplayName returns the human-readable name for the given column key.
// Falls back to the key itself if no mapping exists.
func DisplayName(key string) string {
	if name, ok := DisplayNames[key]; ok {
		return name
	}
	return key
}

// ---------------------------------------------------------------------------
// Default column sets
// ---------------------------------------------------------------------------

// MediaTableColumns returns the default column keys for the media table view.
func MediaTableColumns() []string {
	return []string{
		ColIndex,
		ColMediaID,
		ColPostID,
		ColUsername,
		ColDate,
		ColMediaType,
		ColFilename,
		ColSize,
		ColPrice,
		ColResponseType,
		ColDownloaded,
	}
}

// ProgressColumns returns the column keys for the progress display.
func ProgressColumns() []string {
	return []string{
		ColIndex,
		ColFilename,
		ColStatus,
		ColProgress,
		ColSpeed,
		ColETA,
	}
}

// TaskColumns returns the column keys for the task display.
func TaskColumns() []string {
	return []string{
		ColIndex,
		ColTask,
		ColStatus,
		ColElapsed,
	}
}
