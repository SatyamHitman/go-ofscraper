// =============================================================================
// FILE: internal/model/placeholder.go
// PURPOSE: Implements path template expansion for generating download file
//          paths and directory structures. Supports template variables like
//          {model_username}, {post_id}, {media_type}, etc. Ports the Python
//          classes/placeholder.py with all placeholder classes.
// =============================================================================

package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Template variable names (constants for all supported placeholders)
// ---------------------------------------------------------------------------

const (
	VarConfigPath      = "config_path"
	VarProfile         = "profile"
	VarSiteName        = "site_name"
	VarSaveLocation    = "save_location"
	VarMyID            = "my_id"
	VarMyUsername       = "my_username"
	VarRoot            = "root"
	VarUsername         = "user_name"
	VarModelUsername    = "model_username"
	VarModelID         = "model_id"
	VarPostID          = "post_id"
	VarMediaID         = "media_id"
	VarFirstLetter     = "first_letter"
	VarMediaType       = "media_type"
	VarValue           = "value"
	VarDate            = "date"
	VarResponseType    = "response_type"
	VarLabel           = "label"
	VarDownloadType    = "download_type"
	VarQuality         = "quality"
	VarFilename        = "file_name"
	VarOriginalFilename = "original_filename"
	VarOnlyFilename    = "only_file_name"
	VarText            = "text"
	VarExt             = "ext"
	VarCurrentPrice    = "current_price"
	VarRegularPrice    = "regular_price"
	VarPromoPrice      = "promo_price"
	VarRenewalPrice    = "renewal_price"
)

// ---------------------------------------------------------------------------
// PlaceholderContext holds all resolved template variables.
// ---------------------------------------------------------------------------

// PlaceholderContext stores resolved variable values used for path template
// expansion. Variables are populated from media, post, user, and config data.
type PlaceholderContext struct {
	// Variables maps placeholder names (e.g., "model_username") to their values.
	Variables map[string]string
}

// NewPlaceholderContext creates an empty placeholder context.
//
// Returns:
//   - A new PlaceholderContext with an initialized (empty) variable map.
func NewPlaceholderContext() *PlaceholderContext {
	return &PlaceholderContext{
		Variables: make(map[string]string),
	}
}

// Set adds or updates a variable in the placeholder context.
//
// Parameters:
//   - key: The variable name (e.g., "model_username").
//   - value: The resolved value for this variable.
func (pc *PlaceholderContext) Set(key, value string) {
	pc.Variables[key] = value
}

// Get retrieves a variable value by name. Returns empty string if not found.
//
// Parameters:
//   - key: The variable name to look up.
//
// Returns:
//   - The variable value, or empty string if the key does not exist.
func (pc *PlaceholderContext) Get(key string) string {
	return pc.Variables[key]
}

// ---------------------------------------------------------------------------
// Template expansion
// ---------------------------------------------------------------------------

// ExpandTemplate replaces all {variable_name} placeholders in a format string
// with their resolved values from the PlaceholderContext.
//
// Parameters:
//   - format: The template string containing {placeholder} tokens.
//   - ctx: The PlaceholderContext with resolved variable values.
//
// Returns:
//   - The expanded string with all placeholders replaced.
func ExpandTemplate(format string, ctx *PlaceholderContext) string {
	result := format
	for key, value := range ctx.Variables {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ---------------------------------------------------------------------------
// MediaPlaceholder handles path generation for media downloads.
// ---------------------------------------------------------------------------

// MediaPlaceholder generates file paths and directory structures for a single
// media item download. It combines config format strings with resolved
// variable values to produce the final filesystem path.
type MediaPlaceholder struct {
	// Media is the media item being downloaded.
	Media *Media

	// Extension is the file extension (e.g., "mp4", "jpg").
	Extension string

	// Context holds all resolved template variables.
	Context *PlaceholderContext

	// resolved paths (set after Init)
	filename string
	mediaDir string
	filePath string
}

// NewMediaPlaceholder creates a new MediaPlaceholder for the given media and extension.
//
// Parameters:
//   - media: The media item to generate paths for.
//   - ext: The file extension to use.
//
// Returns:
//   - A new MediaPlaceholder instance.
func NewMediaPlaceholder(media *Media, ext string) *MediaPlaceholder {
	return &MediaPlaceholder{
		Media:     media,
		Extension: ext,
		Context:   NewPlaceholderContext(),
	}
}

// SetBaseVariables populates the base-level template variables common to all
// placeholder types (config path, profile, site name, save location, etc.).
//
// Parameters:
//   - configPath: The configuration directory path.
//   - profile: The active profile name.
//   - saveLocation: The base download directory.
//   - myID: The current authenticated user's ID.
//   - myUsername: The current authenticated user's username.
func (mp *MediaPlaceholder) SetBaseVariables(configPath, profile, saveLocation string, myID int64, myUsername string) {
	mp.Context.Set(VarConfigPath, configPath)
	mp.Context.Set(VarProfile, profile)
	mp.Context.Set(VarSiteName, "Onlyfans")
	mp.Context.Set(VarSaveLocation, saveLocation)
	mp.Context.Set(VarMyID, fmt.Sprintf("%d", myID))
	mp.Context.Set(VarMyUsername, myUsername)
	mp.Context.Set(VarRoot, saveLocation)
}

// SetMediaVariables populates media-specific template variables from the
// Media and its parent Post.
//
// Parameters:
//   - username: The creator's username.
//   - modelID: The creator's user ID.
//   - quality: The selected quality string (e.g., "source", "720").
func (mp *MediaPlaceholder) SetMediaVariables(username string, modelID int64, quality string) {
	m := mp.Media

	mp.Context.Set(VarUsername, username)
	mp.Context.Set(VarModelUsername, username)
	mp.Context.Set(VarModelID, fmt.Sprintf("%d", modelID))
	mp.Context.Set(VarPostID, fmt.Sprintf("%d", m.PostID))
	mp.Context.Set(VarMediaID, fmt.Sprintf("%d", m.ID))
	mp.Context.Set(VarExt, mp.Extension)

	// First letter of username (capitalized)
	if len(username) > 0 {
		mp.Context.Set(VarFirstLetter, strings.ToUpper(username[:1]))
	}

	// Media type (capitalized)
	mp.Context.Set(VarMediaType, strings.Title(string(m.MediaType())))

	// Value: "Free" or "Paid"
	if m.Value == "paid" {
		mp.Context.Set(VarValue, "Paid")
	} else {
		mp.Context.Set(VarValue, "Free")
	}

	// Response type
	mp.Context.Set(VarResponseType, m.ResponseType)

	// Label
	mp.Context.Set(VarLabel, m.Label)

	// Download type
	mp.Context.Set(VarDownloadType, string(m.DownloadKind()))

	// Quality
	mp.Context.Set(VarQuality, quality)

	// Filename from URL
	origFilename := m.Filename()
	mp.Context.Set(VarOriginalFilename, origFilename)
	mp.Context.Set(VarOnlyFilename, origFilename)

	// Filename with quality suffix
	if quality != "" && quality != "source" {
		mp.Context.Set(VarFilename, origFilename+"_"+quality)
	} else {
		mp.Context.Set(VarFilename, origFilename)
	}

	// Text (from parent post, cleaned for filenames)
	if m.Post != nil {
		mp.Context.Set(VarText, m.Post.FileText())
	}

	// Date
	mp.Context.Set(VarDate, m.FormattedDate())
}

// SetPriceVariables populates price-related template variables from the user model.
//
// Parameters:
//   - user: The creator's User model with pricing data.
func (mp *MediaPlaceholder) SetPriceVariables(user *User) {
	if user == nil {
		mp.Context.Set(VarCurrentPrice, "Unknown_Price")
		mp.Context.Set(VarRegularPrice, "Unknown_Price")
		mp.Context.Set(VarPromoPrice, "Unknown_Price")
		mp.Context.Set(VarRenewalPrice, "Unknown_Price")
		return
	}

	mp.Context.Set(VarCurrentPrice, fmt.Sprintf("%.2f", user.FinalCurrentPrice()))
	mp.Context.Set(VarRegularPrice, fmt.Sprintf("%.2f", user.RegularPrice()))
	mp.Context.Set(VarPromoPrice, fmt.Sprintf("%.2f", user.FinalPromoPrice()))
	mp.Context.Set(VarRenewalPrice, fmt.Sprintf("%.2f", user.FinalRenewalPrice()))
}

// GenerateDir expands the directory format template and resolves the path.
//
// Parameters:
//   - dirFormat: The directory template string (e.g., "{model_username}/{responsetype}/{mediatype}/").
//   - rootDir: The root download directory path.
//   - createDirs: Whether to create the directory if it doesn't exist.
//
// Returns:
//   - The resolved directory path, and any error from directory creation.
func (mp *MediaPlaceholder) GenerateDir(dirFormat, rootDir string, createDirs bool) (string, error) {
	expanded := ExpandTemplate(dirFormat, mp.Context)
	dir := filepath.Join(rootDir, expanded)
	mp.mediaDir = dir

	if createDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create media directory %s: %w", dir, err)
		}
	}

	return dir, nil
}

// GenerateFilename expands the filename format template.
//
// Parameters:
//   - fileFormat: The filename template string (e.g., "{filename}.{ext}").
//
// Returns:
//   - The resolved filename string.
func (mp *MediaPlaceholder) GenerateFilename(fileFormat string) string {
	expanded := ExpandTemplate(fileFormat, mp.Context)
	mp.filename = FileCleanup(expanded)
	return mp.filename
}

// FilePath returns the full resolved file path (directory + filename).
//
// Returns:
//   - The complete file path string.
func (mp *MediaPlaceholder) FilePath() string {
	if mp.filePath != "" {
		return mp.filePath
	}
	mp.filePath = filepath.Join(mp.mediaDir, mp.filename)
	return mp.filePath
}

// Dir returns the resolved media directory path.
//
// Returns:
//   - The directory path string.
func (mp *MediaPlaceholder) Dir() string {
	return mp.mediaDir
}

// ---------------------------------------------------------------------------
// DatabasePlaceholder handles database file path generation.
// ---------------------------------------------------------------------------

// DatabasePlaceholder generates paths for per-model SQLite database files.
type DatabasePlaceholder struct {
	Context *PlaceholderContext
}

// NewDatabasePlaceholder creates a new DatabasePlaceholder.
//
// Returns:
//   - A new DatabasePlaceholder instance.
func NewDatabasePlaceholder() *DatabasePlaceholder {
	return &DatabasePlaceholder{
		Context: NewPlaceholderContext(),
	}
}

// DatabasePath generates the database file path for a specific model.
//
// Parameters:
//   - metadataFormat: The metadata path template (e.g., "{configpath}/{profile}/.data/{model_id}").
//   - configPath: The configuration directory.
//   - profile: The active profile name.
//   - modelID: The creator's user ID.
//   - modelUsername: The creator's username.
//
// Returns:
//   - The resolved database file path string.
func (dp *DatabasePlaceholder) DatabasePath(metadataFormat, configPath, profile string, modelID int64, modelUsername string) string {
	dp.Context.Set(VarConfigPath, configPath)
	dp.Context.Set(VarProfile, profile)
	dp.Context.Set(VarModelID, fmt.Sprintf("%d", modelID))
	dp.Context.Set(VarModelUsername, modelUsername)
	dp.Context.Set(VarUsername, modelUsername)

	if len(modelUsername) > 0 {
		dp.Context.Set(VarFirstLetter, strings.ToUpper(modelUsername[:1]))
	}

	return ExpandTemplate(metadataFormat, dp.Context)
}
