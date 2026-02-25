// =============================================================================
// FILE: internal/config/menu.go
// PURPOSE: Config menu operations for interactive config editing. Provides
//          functions to read/modify config values via the TUI menu system.
//          Ports Python utils/config/menu.py.
// =============================================================================

package config

// ---------------------------------------------------------------------------
// Menu value types
// ---------------------------------------------------------------------------

// MenuCategory groups related config fields for the interactive editor.
type MenuCategory struct {
	Name   string      // Display name for the category
	Fields []MenuField // Fields within this category
}

// MenuField represents a single editable config field in the menu.
type MenuField struct {
	Key          string      // JSON key path (e.g., "file_options.save_location")
	Label        string      // Human-readable label
	Description  string      // Help text for the field
	Type         string      // Field type: "string", "int", "bool", "choice", "list"
	Choices      []string    // Valid choices (for "choice" type)
	CurrentValue interface{} // Current value from config
}

// ---------------------------------------------------------------------------
// Menu builders
// ---------------------------------------------------------------------------

// BuildMenuCategories constructs the menu structure from the current config.
// Returns all categories with their editable fields populated with current values.
//
// Returns:
//   - Slice of MenuCategory representing the full config editor.
func BuildMenuCategories() []MenuCategory {
	cfg := Get()

	return []MenuCategory{
		{
			Name: "File Options",
			Fields: []MenuField{
				{Key: "file_options.save_location", Label: "Save Location", Type: "string", CurrentValue: cfg.File.SaveLocation},
				{Key: "file_options.dir_format", Label: "Directory Format", Type: "string", CurrentValue: cfg.File.DirFormat},
				{Key: "file_options.file_format", Label: "File Format", Type: "string", CurrentValue: cfg.File.FileFormat},
				{Key: "file_options.textlength", Label: "Text Length", Type: "int", CurrentValue: cfg.File.TextLength},
				{Key: "file_options.space_replacer", Label: "Space Replacer", Type: "string", CurrentValue: cfg.File.SpaceReplacer},
				{Key: "file_options.date", Label: "Date Format", Type: "string", CurrentValue: cfg.File.DateFormat},
				{Key: "file_options.text_type_default", Label: "Text Type", Type: "choice", Choices: TextTypeOptions, CurrentValue: cfg.File.TextType},
				{Key: "file_options.truncation_default", Label: "Truncation", Type: "bool", CurrentValue: cfg.File.Truncation},
			},
		},
		{
			Name: "Download Options",
			Fields: []MenuField{
				{Key: "download_options.filter", Label: "Media Filter", Type: "list", Choices: MediaFilterOptions, CurrentValue: cfg.Download.Filter},
				{Key: "download_options.auto_resume", Label: "Auto Resume", Type: "bool", CurrentValue: cfg.Download.AutoResume},
				{Key: "download_options.system_free_min", Label: "Min Free Space (bytes)", Type: "int", CurrentValue: cfg.Download.SystemFreeMin},
				{Key: "download_options.max_post_count", Label: "Max Post Count", Type: "int", CurrentValue: cfg.Download.MaxPostCount},
			},
		},
		{
			Name: "CDM Options",
			Fields: []MenuField{
				{Key: "cdm_options.key-mode-default", Label: "Key Mode", Type: "choice", Choices: KeyOptions, CurrentValue: cfg.CDM.KeyMode},
				{Key: "cdm_options.private-key", Label: "Private Key Path", Type: "string", CurrentValue: cfg.CDM.PrivateKey},
				{Key: "cdm_options.client-id", Label: "Client ID Path", Type: "string", CurrentValue: cfg.CDM.ClientID},
			},
		},
		{
			Name: "Advanced Options",
			Fields: []MenuField{
				{Key: "advanced_options.dynamic-mode-default", Label: "Dynamic Mode", Type: "choice", Choices: DynamicOptions, CurrentValue: cfg.Advanced.DynamicMode},
				{Key: "advanced_options.cache-mode", Label: "Cache Mode", Type: "choice", Choices: CacheModes, CurrentValue: cfg.Advanced.CacheMode},
				{Key: "advanced_options.rotate_logs", Label: "Rotate Logs", Type: "bool", CurrentValue: cfg.Advanced.RotateLogs},
				{Key: "advanced_options.ssl_verify", Label: "SSL Verify", Type: "bool", CurrentValue: cfg.Advanced.SSLVerify},
				{Key: "advanced_options.sanitize_text", Label: "Sanitize DB Text", Type: "bool", CurrentValue: cfg.Advanced.SanitizeText},
			},
		},
	}
}
