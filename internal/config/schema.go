// =============================================================================
// FILE: internal/config/schema.go
// PURPOSE: Defines the configuration schema with all default values. Provides
//          the canonical config structure used for JSON serialization and
//          schema validation. Ports Python utils/config/schema.py.
// =============================================================================

package config

// ---------------------------------------------------------------------------
// AppConfig is the top-level configuration structure.
// ---------------------------------------------------------------------------

// AppConfig represents the complete application configuration.
// All fields have JSON tags matching the config file keys.
type AppConfig struct {
	MainProfile string            `json:"main_profile"`
	Metadata    string            `json:"metadata"`
	Discord     string            `json:"discord"`
	File        FileOptions       `json:"file_options"`
	Download    DownloadOptions   `json:"download_options"`
	Binary      BinaryOptions     `json:"binary_options"`
	CDM         CDMOptions        `json:"cdm_options"`
	Performance PerformanceOpts   `json:"performance_options"`
	Content     ContentFilterOpts `json:"content_filter_options"`
	Advanced    AdvancedOptions   `json:"advanced_options"`
	Scripts     ScriptOptions     `json:"script_options"`
	Response    ResponseTypeMap   `json:"responsetype"`
}

// FileOptions controls file/path formatting and text processing.
type FileOptions struct {
	SaveLocation string `json:"save_location"`
	DirFormat    string `json:"dir_format"`
	FileFormat   string `json:"file_format"`
	TextLength   int    `json:"textlength"`
	SpaceReplacer string `json:"space_replacer"`
	DateFormat   string `json:"date"`
	TextType     string `json:"text_type_default"`
	Truncation   bool   `json:"truncation_default"`
}

// DownloadOptions controls download behavior and limits.
type DownloadOptions struct {
	Filter       []string `json:"filter"`
	AutoResume   bool     `json:"auto_resume"`
	SystemFreeMin int64   `json:"system_free_min"`
	MaxPostCount int      `json:"max_post_count"`
}

// BinaryOptions specifies paths to external binaries.
type BinaryOptions struct {
	FFmpeg string `json:"ffmpeg"`
}

// CDMOptions configures the Content Decryption Module for DRM.
type CDMOptions struct {
	PrivateKey string `json:"private-key"`
	ClientID   string `json:"client-id"`
	KeyMode    string `json:"key-mode-default"`
}

// PerformanceOpts controls concurrency and bandwidth limits.
type PerformanceOpts struct {
	DownloadSems  int   `json:"download_sems"`
	DownloadLimit int64 `json:"download_limit"`
}

// ContentFilterOpts controls media content filtering by size and duration.
type ContentFilterOpts struct {
	BlockAds    bool  `json:"block_ads"`
	FileSizeMax int64 `json:"file_size_max"`
	FileSizeMin int64 `json:"file_size_min"`
	LengthMax   int   `json:"length_max"`
	LengthMin   int   `json:"length_min"`
}

// AdvancedOptions holds advanced/power-user configuration.
type AdvancedOptions struct {
	DynamicMode       string   `json:"dynamic-mode-default"`
	DownloadBars      bool     `json:"downloadbars"`
	CacheMode         string   `json:"cache-mode"`
	RotateLogs        bool     `json:"rotate_logs"`
	SanitizeText      bool     `json:"sanitize_text"`
	TempDir           string   `json:"temp_dir"`
	RemoveHashMatch   bool     `json:"remove_hash_match"`
	InfiniteLoopMode  string   `json:"infinite_loop_action_mode"`
	EnableAutoAfter   bool     `json:"enable_auto_after"`
	DefaultUserList   []string `json:"default_user_list"`
	DefaultBlackList  []string `json:"default_black_list"`
	LogsExpireTime    int      `json:"logs_expire_time"`
	SSLVerify         bool     `json:"ssl_verify"`
	EnvFiles          []string `json:"env_files"`
}

// ScriptOptions specifies paths to user-defined hook scripts.
type ScriptOptions struct {
	AfterActionScript   string `json:"after_action_script"`
	PostScript          string `json:"post_script"`
	NamingScript        string `json:"naming_script"`
	AfterDownloadScript string `json:"after_download_script"`
	SkipDownloadScript  string `json:"skip_download_script"`
}

// ResponseTypeMap maps API response types to display names.
type ResponseTypeMap struct {
	Timeline   string `json:"timeline"`
	Message    string `json:"message"`
	Archived   string `json:"archived"`
	Paid       string `json:"paid"`
	Stories    string `json:"stories"`
	Highlights string `json:"highlights"`
	Profile    string `json:"profile"`
	Pinned     string `json:"pinned"`
	Streams    string `json:"streams"`
}

// ---------------------------------------------------------------------------
// DefaultConfig returns a new AppConfig populated with all default values.
// ---------------------------------------------------------------------------

// DefaultConfig creates a new configuration with all fields set to their
// default values. This serves as the base config that gets merged with
// user overrides from the config file and CLI flags.
//
// Returns:
//   - A fully-populated AppConfig with default values.
func DefaultConfig() AppConfig {
	return AppConfig{
		MainProfile: DefaultProfile,
		Metadata:    DefaultMetadataFormat,
		Discord:     DefaultDiscord,
		File: FileOptions{
			SaveLocation:  "{home}/Data/ofscraper",
			DirFormat:     DefaultDirFormat,
			FileFormat:    DefaultFileFormat,
			TextLength:    DefaultTextLength,
			SpaceReplacer: DefaultSpaceReplacer,
			DateFormat:    DefaultDateFormat,
			TextType:      DefaultTextType,
			Truncation:    DefaultTruncation,
		},
		Download: DownloadOptions{
			Filter:        []string{"Images", "Audios", "Videos"},
			AutoResume:    DefaultResume,
			SystemFreeMin: DefaultSystemFreeMin,
			MaxPostCount:  DefaultMaxCount,
		},
		Binary: BinaryOptions{
			FFmpeg: DefaultFFmpeg,
		},
		CDM: CDMOptions{
			PrivateKey: "",
			ClientID:   "",
			KeyMode:    DefaultKeyMode,
		},
		Performance: PerformanceOpts{
			DownloadSems:  DefaultDownloadSem,
			DownloadLimit: DefaultDownloadLimit,
		},
		Content: ContentFilterOpts{
			BlockAds:    DefaultBlockAds,
			FileSizeMax: DefaultFileSizeMax,
			FileSizeMin: DefaultFileSizeMin,
			LengthMax:   DefaultMaxLength,
			LengthMin:   DefaultMinLength,
		},
		Advanced: AdvancedOptions{
			DynamicMode:      DefaultDynamicRule,
			DownloadBars:     DefaultProgress,
			CacheMode:        DefaultCacheMode,
			RotateLogs:       DefaultRotateLogs,
			SanitizeText:     DefaultSanitizeDB,
			TempDir:          "",
			RemoveHashMatch:  false,
			InfiniteLoopMode: "",
			EnableAutoAfter:  DefaultEnableAutoAfter,
			DefaultUserList:  []string{DefaultUserList},
			DefaultBlackList: []string{},
			LogsExpireTime:   0,
			SSLVerify:        DefaultSSLValidation,
			EnvFiles:         []string{},
		},
		Scripts: ScriptOptions{},
		Response: ResponseTypeMap{
			Timeline:   "Posts",
			Message:    "Messages",
			Archived:   "Archived",
			Paid:       "Messages",
			Stories:    "Stories",
			Highlights: "Stories",
			Profile:    "Profile",
			Pinned:     "Posts",
			Streams:    "Streams",
		},
	}
}
