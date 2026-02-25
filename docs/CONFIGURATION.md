# Configuration Reference

GoFScraper uses a JSON configuration file. Default location: `~/.config/gofscraper/config.json`.

Override with `--config` or the `OF_CONFIG_FILE` environment variable.

---

## File Structure

```json
{
  "main_profile": "main_profile",
  "metadata": "{configpath}/{profile}/.data/{model_id}",
  "discord": "",
  "file_options": { ... },
  "download_options": { ... },
  "binary_options": { ... },
  "cdm_options": { ... },
  "performance_options": { ... },
  "content_filter_options": { ... },
  "advanced_options": { ... },
  "script_options": { ... },
  "responsetype": { ... }
}
```

---

## Top-Level Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `main_profile` | string | `"main_profile"` | Active configuration profile name |
| `metadata` | string | `"{configpath}/{profile}/.data/{model_id}"` | Database/metadata path template |
| `discord` | string | `""` | Discord webhook URL for notifications (empty = disabled) |

---

## file_options

Controls file naming, paths, and text formatting.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `save_location` | string | `"{home}/Data/ofscraper"` | Root directory for downloaded content |
| `dir_format` | string | `"{model_username}/{responsetype}/{mediatype}/"` | Directory structure template |
| `file_format` | string | `"{filename}.{ext}"` | Filename template |
| `textlength` | int | `0` | Max text length in filenames (0 = unlimited) |
| `space_replacer` | string | `" "` | Character to replace spaces in paths |
| `date` | string | `"MM-DD-YYYY"` | Date format for display |
| `text_type_default` | string | `"letter"` | Truncation mode: `"letter"` or `"word"` |
| `truncation_default` | bool | `true` | Enable path length truncation |

### Path Template Variables

Available in `save_location`, `dir_format`, `file_format`, and `metadata`:

| Variable | Description | Example |
|----------|-------------|---------|
| `{model_username}` | Creator's username | `janedoe` |
| `{model_id}` | Creator's numeric ID | `123456` |
| `{responsetype}` | Content area display name | `Posts`, `Messages` |
| `{mediatype}` | Media type directory | `images`, `videos`, `audios` |
| `{filename}` | Original filename from URL | `photo_2024` |
| `{ext}` | File extension | `jpg`, `mp4` |
| `{media_id}` | Media item ID | `789012` |
| `{post_id}` | Parent post ID | `345678` |
| `{date}` | Post date (formatted per `date` setting) | `01-15-2024` |
| `{text}` | Post text (truncated per `textlength`) | `Check out...` |
| `{label}` | Label name (if applicable) | `favorites` |
| `{count}` | Media index within post | `1`, `2`, `3` |
| `{home}` | User home directory | `/home/user` |
| `{configpath}` | Config directory path | `~/.config/gofscraper` |
| `{profile}` | Active profile name | `main_profile` |
| `{sitename}` | Site name constant | `Onlyfans` |

---

## download_options

Controls what gets downloaded and how.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `filter` | []string | `["Images", "Audios", "Videos"]` | Media types to download |
| `auto_resume` | bool | `true` | Resume interrupted downloads |
| `system_free_min` | int | `0` | Minimum free disk space in bytes (0 = no check) |
| `max_post_count` | int | `0` | Max posts to process per user (0 = unlimited) |

**Filter values:** `"Images"`, `"Audios"`, `"Videos"`

---

## binary_options

Paths to external tools.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `ffmpeg` | string | `""` | Path to FFmpeg binary (empty = auto-detect on PATH) |

FFmpeg is required for DRM-protected content decryption and media concatenation.

---

## cdm_options

Content Decryption Module configuration for DRM-protected media.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `private-key` | string | `""` | Path to Widevine private key file |
| `client-id` | string | `""` | Path to Widevine client ID file |
| `key-mode-default` | string | `"cdrm"` | CDM mode: `"cdrm"` (service) or `"manual"` (local keys) |

**Modes:**
- `cdrm` -- Uses an external CDRM decryption service (no local keys needed)
- `manual` -- Uses local Widevine device files (private key + client ID)

---

## performance_options

Concurrency and bandwidth controls.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `download_sems` | int | `6` | Max concurrent downloads |
| `download_limit` | int | `0` | Bandwidth limit in bytes/sec (0 = unlimited) |

---

## content_filter_options

Filter content by size and duration before downloading.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `block_ads` | bool | `false` | Filter out advertisement posts |
| `file_size_max` | int | `0` | Max file size in bytes (0 = no limit) |
| `file_size_min` | int | `0` | Min file size in bytes (0 = no limit) |
| `length_max` | int | `0` | Max media duration in seconds (0 = no limit) |
| `length_min` | int | `0` | Min media duration in seconds (0 = no limit) |

---

## advanced_options

Power-user settings.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `dynamic-mode-default` | string | `"digital"` | Auth signing rule provider |
| `downloadbars` | bool | `false` | Show per-file download progress bars |
| `cache-mode` | string | `"sqlite"` | Cache backend: `"sqlite"`, `"json"`, `"disabled"` |
| `rotate_logs` | bool | `true` | Enable log file rotation |
| `sanitize_text` | bool | `false` | Sanitize text before storing in database |
| `temp_dir` | string | `""` | Temporary directory for downloads (empty = system default) |
| `remove_hash_match` | bool | `false` | Remove files that match hash of existing downloads |
| `infinite_loop_action_mode` | string | `""` | Action on infinite loop detection |
| `enable_auto_after` | bool | `false` | Auto-set "after" timestamp from last scrape |
| `default_user_list` | []string | `["main"]` | Default user list names |
| `default_black_list` | []string | `[]` | Default blacklisted usernames |
| `logs_expire_time` | int | `0` | Log file expiry in days (0 = never) |
| `ssl_verify` | bool | `true` | Verify SSL certificates |
| `env_files` | []string | `[]` | Additional `.env` files to load |

**Dynamic rule providers:** `"digitalcriminals"`, `"manual"`, `"generic"`, `"datawhores"`, `"xagler"`, `"rafa"`

---

## script_options

Hook scripts executed at various stages. Each value is a path to a shell script.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `after_action_script` | string | `""` | Run after each download/like action completes |
| `post_script` | string | `""` | Run after each post is processed |
| `naming_script` | string | `""` | Custom filename generation script |
| `after_download_script` | string | `""` | Run after each file download completes |
| `skip_download_script` | string | `""` | Script to decide whether to skip a download |

### Script Environment Variables

Scripts receive context via environment variables:

| Variable | Available In | Description |
|----------|-------------|-------------|
| `OF_USERNAME` | after_action, after_like | Creator's username |
| `OF_ACTION` | after_action | Action name (download/like/unlike) |
| `OF_MEDIA_COUNT` | after_action | Number of media items processed |
| `OF_POST_COUNT` | after_like | Number of posts processed |
| `OF_FILE_PATH` | after_download | Full path to downloaded file |
| `OF_TOTAL_SIZE` | skip_download | Total expected file size |

---

## responsetype

Maps API content areas to display directory names.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `timeline` | string | `"Posts"` | Timeline posts directory name |
| `message` | string | `"Messages"` | Messages directory name |
| `archived` | string | `"Archived"` | Archived posts directory name |
| `paid` | string | `"Messages"` | Purchased content directory name |
| `stories` | string | `"Stories"` | Stories directory name |
| `highlights` | string | `"Stories"` | Highlights directory name |
| `profile` | string | `"Profile"` | Profile media directory name |
| `pinned` | string | `"Posts"` | Pinned posts directory name |
| `streams` | string | `"Streams"` | Streams directory name |

---

## Environment Variables

All configuration values can be overridden with environment variables. The naming pattern is:

```
OF_<SECTION>_<FIELD>
```

Common overrides:

| Variable | Overrides | Example |
|----------|-----------|---------|
| `OF_CONFIG_FILE` | Config file path | `/path/to/config.json` |
| `OF_AUTH_FILE` | Auth file path | `/path/to/auth.json` |
| `OF_CONFIG_DIR` | Config directory | `~/.config/gofscraper` |
| `OF_LOG_LEVEL` | Log level | `DEBUG`, `INFO`, `WARN`, `ERROR` |
| `OF_DISCORD_WEBHOOK` | Discord webhook URL | `https://discord.com/api/webhooks/...` |
| `OF_MAX_CONNECTIONS` | Max HTTP connections | `10` |
| `OF_REQUEST_TIMEOUT` | HTTP timeout (seconds) | `30` |
| `OF_RATE_LIMIT_RPS` | Rate limit (requests/sec) | `2.0` |
| `OF_DOWNLOAD_DIR` | Override save location | `/data/downloads` |
| `OF_CDRM_SERVER` | CDRM service URL | `http://localhost:8080` |

Boolean env vars accept: `true`, `1`, `yes`, `on` (case-insensitive).

---

## Profiles

Profiles allow multiple independent configurations. Each profile has its own config, auth, and database.

```bash
# Use a specific profile
gofscraper scraper --profile myprofile

# Profiles are stored at:
# ~/.config/gofscraper/<profile_name>/config.json
# ~/.config/gofscraper/<profile_name>/auth.json
```

---

## Example: Full Configuration

```json
{
  "main_profile": "main_profile",
  "metadata": "{configpath}/{profile}/.data/{model_id}",
  "discord": "https://discord.com/api/webhooks/123/abc",
  "file_options": {
    "save_location": "/data/ofscraper",
    "dir_format": "{model_username}/{responsetype}/{mediatype}/",
    "file_format": "{filename}.{ext}",
    "textlength": 50,
    "space_replacer": "_",
    "date": "YYYY-MM-DD",
    "text_type_default": "word",
    "truncation_default": true
  },
  "download_options": {
    "filter": ["Images", "Videos", "Audios"],
    "auto_resume": true,
    "system_free_min": 1073741824,
    "max_post_count": 0
  },
  "binary_options": {
    "ffmpeg": "/usr/bin/ffmpeg"
  },
  "cdm_options": {
    "private-key": "",
    "client-id": "",
    "key-mode-default": "cdrm"
  },
  "performance_options": {
    "download_sems": 8,
    "download_limit": 0
  },
  "content_filter_options": {
    "block_ads": true,
    "file_size_max": 0,
    "file_size_min": 0,
    "length_max": 0,
    "length_min": 0
  },
  "advanced_options": {
    "dynamic-mode-default": "digital",
    "downloadbars": true,
    "cache-mode": "sqlite",
    "rotate_logs": true,
    "sanitize_text": false,
    "temp_dir": "",
    "remove_hash_match": false,
    "enable_auto_after": true,
    "default_user_list": ["main"],
    "default_black_list": [],
    "ssl_verify": true,
    "env_files": []
  },
  "script_options": {
    "after_action_script": "",
    "post_script": "",
    "naming_script": "",
    "after_download_script": "",
    "skip_download_script": ""
  },
  "responsetype": {
    "timeline": "Posts",
    "message": "Messages",
    "archived": "Archived",
    "paid": "Messages",
    "stories": "Stories",
    "highlights": "Stories",
    "profile": "Profile",
    "pinned": "Posts",
    "streams": "Streams"
  }
}
```
