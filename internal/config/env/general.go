// =============================================================================
// FILE: internal/config/env/general.go
// PURPOSE: General-purpose environment variable defaults including app token,
//          display settings, placeholders, and regex patterns. Ports Python
//          of_env/values/general.py and of_env/values/main.py.
// =============================================================================

package env

// ---------------------------------------------------------------------------
// Save path
// ---------------------------------------------------------------------------

// DefaultSavePath returns the default base save location for downloads.
// Uses the user's home directory with a Data/ofscraper subfolder.
func DefaultSavePath() string {
	return GetString("OF_SAVE_PATH", "{home}/Data/ofscraper")
}

// ---------------------------------------------------------------------------
// Display settings
// ---------------------------------------------------------------------------

// ShowAvatar returns whether to display creator avatars in output.
func ShowAvatar() bool {
	return GetBool("OF_SHOW_AVATAR", true)
}

// ShowResultsLog returns whether to display the results summary log.
func ShowResultsLog() bool {
	return GetBool("OF_SHOW_RESULTS_LOG", true)
}

// RefreshScreenMs returns the TUI refresh interval in milliseconds.
func RefreshScreenMs() int {
	return GetInt("OF_REFRESH_SCREEN", 50)
}

// ---------------------------------------------------------------------------
// Content filtering
// ---------------------------------------------------------------------------

// FilterSelfMedia returns whether to filter out the user's own media.
func FilterSelfMedia() bool {
	return GetBool("OF_FILTER_SELF_MEDIA", true)
}

// AllowDupeMedia returns whether to allow duplicate media downloads.
func AllowDupeMedia() bool {
	return GetBool("OF_ALLOW_DUPE_MEDIA", false)
}

// ContinueBool returns the default continue/auto-proceed behavior.
func ContinueBool() bool {
	return GetBool("OF_CONTINUE_BOOL", true)
}

// FileCountPlaceholder returns whether to use file count in filenames.
func FileCountPlaceholder() bool {
	return GetBool("OF_FILE_COUNT_PLACEHOLDER", true)
}
