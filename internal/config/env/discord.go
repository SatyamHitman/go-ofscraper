// =============================================================================
// FILE: internal/config/env/discord.go
// PURPOSE: Discord webhook environment variable defaults.
//          Ports Python of_env/values/req/discord.py.
// =============================================================================

package env

// DiscordWebhookURL returns the Discord webhook URL for notifications.
func DiscordWebhookURL() string {
	return GetString("OF_DISCORD_WEBHOOK", "")
}

// DiscordThreadOverride returns whether to override Discord thread settings.
func DiscordThreadOverride() bool {
	return GetBool("OF_DISCORD_THREAD_OVERRIDE", false)
}

// DiscordAsync returns whether Discord messages should be sent asynchronously.
func DiscordAsync() bool {
	return GetBool("OF_DISCORD_ASYNC", false)
}
