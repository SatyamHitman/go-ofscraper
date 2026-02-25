// =============================================================================
// FILE: internal/config/constants.go
// PURPOSE: Central repository of all application constants. Defines key mode
//          options, dynamic rule providers, metadata modes, and other
//          enumerated values used across the application. Ports Python
//          utils/const.py and the various default values scattered across
//          the Python codebase.
// =============================================================================

package config

// ---------------------------------------------------------------------------
// Key mode options (CDM provider selection)
// ---------------------------------------------------------------------------

// KeyOptions defines the supported key/CDM provider modes for DRM decryption.
var KeyOptions = []string{"cdrm", "manual"}

// ---------------------------------------------------------------------------
// Dynamic rule source options
// ---------------------------------------------------------------------------

// DynamicOptions defines the supported dynamic rule providers for auth signing.
var DynamicOptions = []string{
	"digitalcriminals",
	"manual",
	"generic",
	"datawhores",
	"xagler",
	"rafa",
}

// DynamicOptionsAll includes all recognized aliases for dynamic rule providers.
// Used for CLI argument validation.
var DynamicOptionsAll = []string{
	"manual",
	"generic",
	"xagler",
	"rafa",
	"datawhores",
	"digitalcriminals",
	"dv",
	"dev",
	"dc",
	"digital",
	"digitals",
}

// ---------------------------------------------------------------------------
// Metadata operation modes
// ---------------------------------------------------------------------------

// MetadataOptions defines the supported metadata operation modes.
var MetadataOptions = []string{"complete", "update", "check"}

// ---------------------------------------------------------------------------
// Cache modes
// ---------------------------------------------------------------------------

// CacheModes defines the supported caching backends.
var CacheModes = []string{"sqlite", "json", "disabled"}

// ---------------------------------------------------------------------------
// Media type filter options
// ---------------------------------------------------------------------------

// MediaFilterOptions defines the recognized media type filter values.
var MediaFilterOptions = []string{"Images", "Audios", "Videos"}

// ---------------------------------------------------------------------------
// Text truncation modes
// ---------------------------------------------------------------------------

// TextTypeOptions defines supported text truncation modes.
var TextTypeOptions = []string{"letter", "word"}

// ---------------------------------------------------------------------------
// Application-wide string constants
// ---------------------------------------------------------------------------

const (
	// AppToken is the OnlyFans application token used in API headers.
	AppToken = "33d57ade8c02dbc5a333db99ff9ae26a"

	// SiteName is the canonical site name used in path templates.
	SiteName = "Onlyfans"

	// DefaultProfile is the default configuration profile name.
	DefaultProfile = "main_profile"

	// DefaultUserList is the default user list name.
	DefaultUserList = "main"

	// DefaultDateFormat is the default date display format.
	DefaultDateFormat = "MM-DD-YYYY"

	// DeletedModelPlaceholder is the username used for deleted/unavailable models.
	DeletedModelPlaceholder = "modeldeleted"

	// ModelPricePlaceholder is the price string used when price is unknown.
	ModelPricePlaceholder = "Unknown_Price"

	// DefaultDirFormat is the default directory template for media downloads.
	DefaultDirFormat = "{model_username}/{responsetype}/{mediatype}/"

	// DefaultFileFormat is the default filename template for media downloads.
	DefaultFileFormat = "{filename}.{ext}"

	// DefaultMetadataFormat is the default metadata/database path template.
	DefaultMetadataFormat = "{configpath}/{profile}/.data/{model_id}"
)

// ---------------------------------------------------------------------------
// Default response type display names
// ---------------------------------------------------------------------------

// DefaultResponseTypeMap maps API response types to their default display names.
// Users can override these in config.
var DefaultResponseTypeMap = map[string]string{
	"message":    "Messages",
	"timeline":   "Posts",
	"archived":   "Archived",
	"paid":       "Messages",
	"stories":    "Stories",
	"highlights": "Stories",
	"profile":    "Profile",
	"pinned":     "Posts",
	"streams":    "Streams",
}

// ---------------------------------------------------------------------------
// Numeric defaults
// ---------------------------------------------------------------------------

const (
	// DefaultThreads is the default number of download threads.
	DefaultThreads = 2

	// DefaultDownloadSem is the default download semaphore count.
	DefaultDownloadSem = 6

	// DefaultMaxCount is the default maximum post count (0 = unlimited).
	DefaultMaxCount = 0

	// DefaultFileSizeMax is the default max file size filter (0 = no limit).
	DefaultFileSizeMax = 0

	// DefaultFileSizeMin is the default min file size filter (0 = no limit).
	DefaultFileSizeMin = 0

	// DefaultMinLength is the default minimum media duration (0 = no limit).
	DefaultMinLength = 0

	// DefaultMaxLength is the default maximum media duration (0 = no limit).
	DefaultMaxLength = 0

	// DefaultTextLength is the default text truncation length (0 = no truncation).
	DefaultTextLength = 0

	// DefaultDownloadLimit is the default bandwidth limit (0 = unlimited).
	DefaultDownloadLimit = 0

	// DefaultSystemFreeMin is the minimum free disk space required (0 = no check).
	DefaultSystemFreeMin = 0
)

// ---------------------------------------------------------------------------
// Boolean defaults
// ---------------------------------------------------------------------------

const (
	// DefaultTruncation enables path truncation by default.
	DefaultTruncation = true

	// DefaultRotateLogs enables log rotation by default.
	DefaultRotateLogs = true

	// DefaultAvatar enables avatar downloads by default.
	DefaultAvatar = true

	// DefaultResume enables download resume by default.
	DefaultResume = true

	// DefaultSSLValidation enables SSL certificate validation by default.
	DefaultSSLValidation = true

	// DefaultBlockAds disables ad blocking by default.
	DefaultBlockAds = false

	// DefaultSanitizeDB disables database text sanitization by default.
	DefaultSanitizeDB = false

	// DefaultProgress disables progress bars by default.
	DefaultProgress = false

	// DefaultInfiniteLoop disables infinite loop mode by default.
	DefaultInfiniteLoop = false

	// DefaultEnableAutoAfter disables auto-after mode by default.
	DefaultEnableAutoAfter = false

	// DefaultIncludeLabelsAll disables include-all-labels by default.
	DefaultIncludeLabelsAll = false

	// DefaultDiscordThreadOverride disables Discord thread override by default.
	DefaultDiscordThreadOverride = false

	// DefaultDiscordAsync disables async Discord by default.
	DefaultDiscordAsync = false

	// DefaultUseWivCacheKey enables Widevine cache key by default.
	DefaultUseWivCacheKey = true

	// DefaultContinueBool enables auto-continue by default.
	DefaultContinueBool = true

	// DefaultFileCountPlaceholder enables file count placeholder by default.
	DefaultFileCountPlaceholder = true

	// DefaultFilterSelfMedia enables self-media filtering by default.
	DefaultFilterSelfMedia = true

	// DefaultAllowDupeMedia disables duplicate media by default.
	DefaultAllowDupeMedia = false

	// DefaultShowAvatar enables avatar display by default.
	DefaultShowAvatar = true

	// DefaultShowResultsLog enables results logging by default.
	DefaultShowResultsLog = true
)

// ---------------------------------------------------------------------------
// String defaults
// ---------------------------------------------------------------------------

const (
	// DefaultSpaceReplacer is the default character to replace spaces in paths.
	DefaultSpaceReplacer = " "

	// DefaultTextType is the default text truncation mode.
	DefaultTextType = "letter"

	// DefaultKeyMode is the default CDM key mode.
	DefaultKeyMode = "cdrm"

	// DefaultCacheMode is the default cache backend.
	DefaultCacheMode = "sqlite"

	// DefaultDynamicRule is the default dynamic rule provider.
	DefaultDynamicRule = "digital"

	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "DEBUG"

	// DefaultFFmpeg is the default FFmpeg binary path (empty = auto-detect).
	DefaultFFmpeg = ""

	// DefaultDiscord is the default Discord webhook URL (empty = disabled).
	DefaultDiscord = ""

	// SuppressLogLevel is the log level threshold for output suppression.
	SuppressLogLevel = 21
)

// ---------------------------------------------------------------------------
// Regex patterns
// ---------------------------------------------------------------------------

const (
	// NumberRegex matches numeric characters.
	NumberRegex = "[0-9]"

	// UsernameRegex matches valid username characters (not a slash).
	UsernameRegex = "[^/]"
)

// ---------------------------------------------------------------------------
// Output suppression levels
// ---------------------------------------------------------------------------

// SuppressOutputs contains log level names that suppress standard output.
var SuppressOutputs = map[string]bool{
	"CRITICAL": true,
	"ERROR":    true,
	"WARNING":  true,
	"OFF":      true,
	"LOW":      true,
	"PROMPT":   true,
}

// ---------------------------------------------------------------------------
// Refresh rate
// ---------------------------------------------------------------------------

const (
	// RefreshScreen is the TUI refresh interval in milliseconds.
	RefreshScreen = 50
)
