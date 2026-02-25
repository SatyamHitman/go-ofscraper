// =============================================================================
// FILE: internal/config/settings.go
// PURPOSE: Settings resolution layer that merges CLI arguments, config file
//          values, and environment defaults into a single resolved Settings
//          struct. Implements the priority chain: CLI args > config > env > defaults.
//          Ports Python utils/settings.py merged_settings.
// =============================================================================

package config

// ---------------------------------------------------------------------------
// Settings holds the fully resolved runtime settings.
// ---------------------------------------------------------------------------

// Settings represents the fully resolved configuration after merging CLI args,
// config file, and environment defaults. This is the single source of truth
// used by all runtime operations.
type Settings struct {
	// Key/CDM mode
	KeyMode string

	// Cache settings
	CacheDisabled    bool
	APICacheDisabled bool

	// Dynamic rules
	DynamicRules string

	// Display
	DownloadBars bool
	DiscordLevel string
	LogLevel     string

	// Text processing
	Truncate      bool
	TextType      string
	SpaceReplacer string
	TextLength    int

	// Lists
	UserList  []string
	BlackList []string

	// Size filters
	SizeMax int64
	SizeMin int64

	// Performance
	DownloadSems   int
	SystemFreeMin  int64
	MaxPostCount   int
	DownloadLimit  int64

	// Media filters
	MediaTypes []string

	// CDM
	PrivateKey string
	ClientID   string

	// Duration filters
	LengthMax int
	LengthMin int

	// Hash dedup
	HashEnabled bool

	// Scripts
	PostScript          string
	AfterActionScript   string
	NamingScript        string
	SkipDownloadScript  string
	AfterDownloadScript string

	// Logging
	RotateLogs     bool
	LogsExpireTime int

	// Download
	AutoResume bool
	AutoAfter  bool
	SSLVerify  bool

	// Env files
	EnvFiles []string

	// Text download
	TextEnabled bool
	TextOnly    bool

	// Anonymous mode
	Anon bool
}

// ---------------------------------------------------------------------------
// ResolveSettings merges all configuration sources into a Settings struct.
// ---------------------------------------------------------------------------

// ResolveSettings creates a fully resolved Settings by merging CLI arguments
// over config file values over environment defaults. Parameters set to their
// zero value in cliOverrides are treated as "not set" and fall through to
// config defaults.
//
// Parameters:
//   - cliOverrides: Settings from CLI argument parsing. Zero-valued fields
//     indicate "use config default".
//
// Returns:
//   - A fully resolved Settings struct.
func ResolveSettings(cliOverrides *Settings) *Settings {
	cfg := Get()
	s := &Settings{}

	// Key mode: CLI > config
	s.KeyMode = coalesceStr(cliOverrides.KeyMode, GetKeyMode())

	// Cache
	s.CacheDisabled = cliOverrides.CacheDisabled
	s.APICacheDisabled = cliOverrides.APICacheDisabled

	// Dynamic rules
	s.DynamicRules = coalesceStr(cliOverrides.DynamicRules, GetDynamicMode())

	// Display
	s.DownloadBars = cliOverrides.DownloadBars || GetDownloadBars()
	s.DiscordLevel = coalesceStr(cliOverrides.DiscordLevel, "")
	s.LogLevel = coalesceStr(cliOverrides.LogLevel, DefaultLogLevel)

	// Text processing
	s.Truncate = GetTruncation()
	s.TextType = coalesceStr(cliOverrides.TextType, GetTextType())
	s.SpaceReplacer = coalesceStr(cliOverrides.SpaceReplacer, GetSpaceReplacer())
	s.TextLength = coalesceInt(cliOverrides.TextLength, GetTextLength())

	// Lists
	s.UserList = coalesceStrSlice(cliOverrides.UserList, GetDefaultUserList())
	s.BlackList = coalesceStrSlice(cliOverrides.BlackList, GetDefaultBlackList())

	// Size filters
	s.SizeMax = coalesceInt64(cliOverrides.SizeMax, GetFileSizeMax())
	s.SizeMin = coalesceInt64(cliOverrides.SizeMin, GetFileSizeMin())

	// Performance
	s.DownloadSems = coalesceInt(cliOverrides.DownloadSems, GetDownloadSemaphores())
	s.SystemFreeMin = coalesceInt64(cliOverrides.SystemFreeMin, GetSystemFreeSize())
	s.MaxPostCount = coalesceInt(cliOverrides.MaxPostCount, GetMaxPostCount())
	s.DownloadLimit = coalesceInt64(cliOverrides.DownloadLimit, GetDownloadLimit())

	// Media types
	s.MediaTypes = coalesceStrSlice(cliOverrides.MediaTypes, GetFilter())

	// CDM
	s.PrivateKey = coalesceStr(cliOverrides.PrivateKey, GetPrivateKey())
	s.ClientID = coalesceStr(cliOverrides.ClientID, GetClientID())

	// Duration
	s.LengthMax = coalesceInt(cliOverrides.LengthMax, cfg.Content.LengthMax)
	s.LengthMin = coalesceInt(cliOverrides.LengthMin, cfg.Content.LengthMin)

	// Hash
	s.HashEnabled = cliOverrides.HashEnabled || GetHashEnabled()

	// Scripts
	s.PostScript = coalesceStr(cliOverrides.PostScript, GetPostScript())
	s.AfterActionScript = coalesceStr(cliOverrides.AfterActionScript, GetAfterActionScript())
	s.NamingScript = coalesceStr(cliOverrides.NamingScript, GetNamingScript())
	s.SkipDownloadScript = coalesceStr(cliOverrides.SkipDownloadScript, GetSkipDownloadScript())
	s.AfterDownloadScript = coalesceStr(cliOverrides.AfterDownloadScript, GetAfterDownloadScript())

	// Logging
	s.RotateLogs = GetRotateLogs()
	s.LogsExpireTime = cfg.Advanced.LogsExpireTime

	// Download
	s.AutoResume = GetAutoResume()
	s.AutoAfter = GetEnableAutoAfter()
	s.SSLVerify = GetSSLVerify()

	// Env files
	s.EnvFiles = coalesceStrSlice(cliOverrides.EnvFiles, cfg.Advanced.EnvFiles)

	// Text download
	s.TextEnabled = cliOverrides.TextEnabled
	s.TextOnly = cliOverrides.TextOnly

	// Anon
	s.Anon = cliOverrides.Anon

	return s
}

// ---------------------------------------------------------------------------
// Coalesce helpers (prefer first non-zero value)
// ---------------------------------------------------------------------------

// coalesceStr returns the first non-empty string.
func coalesceStr(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// coalesceInt returns the first non-zero integer.
func coalesceInt(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

// coalesceInt64 returns the first non-zero int64.
func coalesceInt64(values ...int64) int64 {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

// coalesceStrSlice returns the first non-empty string slice.
func coalesceStrSlice(values ...[]string) []string {
	for _, v := range values {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}
