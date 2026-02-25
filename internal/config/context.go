// =============================================================================
// FILE: internal/config/context.go
// PURPOSE: Provides context-based configuration helpers. Allows storing and
//          retrieving config values from Go context.Context for use in
//          request-scoped operations. Ports Python utils/config/utils/context.py.
// =============================================================================

package config

import "context"

// ---------------------------------------------------------------------------
// Context keys
// ---------------------------------------------------------------------------

// contextKey is an unexported type for config context keys to avoid collisions.
type contextKey string

const (
	// configCtxKey stores the AppConfig in context.
	configCtxKey contextKey = "gofscraper_config"

	// settingsCtxKey stores the resolved Settings in context.
	settingsCtxKey contextKey = "gofscraper_settings"
)

// ---------------------------------------------------------------------------
// Context accessors
// ---------------------------------------------------------------------------

// WithConfig returns a new context with the AppConfig attached.
//
// Parameters:
//   - ctx: The parent context.
//   - cfg: The config to attach.
//
// Returns:
//   - A derived context containing the config.
func WithConfig(ctx context.Context, cfg *AppConfig) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}

// FromContext retrieves the AppConfig from a context. Returns the global
// config if none is attached to the context.
//
// Parameters:
//   - ctx: The context to read from.
//
// Returns:
//   - The AppConfig from context, or the global config as fallback.
func FromContext(ctx context.Context) *AppConfig {
	if cfg, ok := ctx.Value(configCtxKey).(*AppConfig); ok {
		return cfg
	}
	return Get()
}

// WithSettings returns a new context with the resolved Settings attached.
//
// Parameters:
//   - ctx: The parent context.
//   - s: The settings to attach.
//
// Returns:
//   - A derived context containing the settings.
func WithSettings(ctx context.Context, s *Settings) context.Context {
	return context.WithValue(ctx, settingsCtxKey, s)
}

// SettingsFromContext retrieves the resolved Settings from a context.
// Returns nil if no settings are attached.
//
// Parameters:
//   - ctx: The context to read from.
//
// Returns:
//   - The Settings from context, or nil.
func SettingsFromContext(ctx context.Context) *Settings {
	if s, ok := ctx.Value(settingsCtxKey).(*Settings); ok {
		return s
	}
	return nil
}
