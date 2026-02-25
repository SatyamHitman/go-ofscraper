# Migration Guide: Python OF-Scraper to GoFScraper

This guide helps existing OF-Scraper (Python) users transition to GoFScraper (Go).

---

## Why Migrate?

| Aspect | Python OF-Scraper | GoFScraper |
|--------|-------------------|------------|
| **Installation** | Python 3.10+, pip, venv, 200+ pip packages | Single binary, zero dependencies |
| **Startup time** | 2-5 seconds (interpreter + imports) | <50ms |
| **Memory usage** | 100-300 MB typical | 20-50 MB typical |
| **Concurrency** | asyncio (single-thread cooperative) | goroutines (true parallelism) |
| **Deployment** | Complex (pip freeze, venv, Python version) | Copy binary, done |
| **Cross-platform** | Requires Python per platform | Native binary per platform |
| **Updates** | pip install, dependency conflicts possible | Replace binary |

---

## Compatibility

### What Stays the Same

- **Config file format** -- Same JSON structure, same field names. Your existing `config.json` works as-is.
- **Auth file format** -- Same JSON structure (`auth_cookie`, `auth_user_agent`, `auth_x_bc`, `auth_user_id`).
- **Database schema** -- Same 10 SQLite tables, same column names. Existing databases are compatible.
- **Directory structure** -- Same path templates, same default directory layout.
- **Filename patterns** -- Same placeholders (`{model_username}`, `{filename}`, `{ext}`, etc.).
- **API behavior** -- Same endpoints, same pagination logic, same auth signing.

### What Changed

| Feature | Python | Go |
|---------|--------|----|
| CLI framework | Click/Cloup | Cobra |
| Flag syntax | `--flag value` | `--flag value` or `--flag=value` |
| Config library | Custom JSON loader | Custom JSON loader (Viper-compatible) |
| TUI framework | Rich + Textual | Bubbletea + Lipgloss |
| Logging | Python logging + Rich | slog + tint (colored) |
| Async model | asyncio | goroutines + channels |
| Package manager | pip | Go modules (vendored in binary) |

---

## Step-by-Step Migration

### 1. Install GoFScraper

Download the binary for your platform from the [Releases](https://github.com/your-org/gofscraper/releases) page, or build from source:

```bash
git clone https://github.com/your-org/gofscraper.git
cd gofscraper
make build
sudo cp bin/gofscraper /usr/local/bin/
```

### 2. Keep Your Config

Your existing config file works without changes:

```bash
# GoFScraper looks in the same default location
ls ~/.config/gofscraper/config.json

# Or specify explicitly
gofscraper scraper --config ~/.config/ofscraper/config.json
```

### 3. Keep Your Auth

Same auth file format:

```bash
ls ~/.config/gofscraper/auth.json
```

### 4. Keep Your Databases

SQLite databases are binary-compatible. GoFScraper reads them as-is:

```bash
# Default location (from metadata path in config)
ls ~/.config/gofscraper/main_profile/.data/*/
```

### 5. Keep Your Downloaded Files

No changes to existing downloads. GoFScraper uses the same path patterns and will recognize already-downloaded files for deduplication.

### 6. Update Any Wrapper Scripts

If you have shell scripts calling the Python version:

```bash
# Before (Python)
ofscraper --action download --posts timeline --usernames user1

# After (Go) -- same flags, slightly different syntax
gofscraper scraper --action download --posts timeline --users user1
```

---

## CLI Flag Mapping

Most flags are identical. Key differences:

| Python OF-Scraper | GoFScraper | Notes |
|-------------------|------------|-------|
| `ofscraper` | `gofscraper scraper` | Main command is now a subcommand |
| `--action` | `--action` / `-a` | Same |
| `--posts` | `--posts` / `-o` | Same |
| `--usernames` | `--users` / `-u` | Renamed for brevity |
| `--excluded-username` | `--excluded-users` | Renamed |
| `--daemon` | `--daemon` / `-d` | Same |
| `--config-path` | `--config` / `-c` | Renamed |
| `--profile` | `--profile` / `-p` | Same |
| `--log-level` | `--log-level` / `-l` | Same values |
| `--no-prompt` | `--no-interactive` / `-n` | Renamed |
| `ofscraper msg_check` | `gofscraper msg_check` | Same subcommand |
| `ofscraper story_check` | `gofscraper story_check` | Same subcommand |
| `ofscraper paid_check` | `gofscraper paid_check` | Same subcommand |
| `ofscraper post_check` | `gofscraper post_check` | Same subcommand |
| `ofscraper metadata` | `gofscraper metadata` | Same subcommand |
| `ofscraper db` | `gofscraper db` | Same subcommand |
| `ofscraper manual` | `gofscraper manual` | Same subcommand |

---

## Config Field Mapping

All config fields are identical between Python and Go versions. No changes needed.

The one difference is that GoFScraper validates the config more strictly. If you had typos or unused fields in your Python config, GoFScraper will ignore them silently (no errors), but they won't have any effect.

---

## Environment Variable Mapping

| Python | Go | Notes |
|--------|-----|-------|
| `OF_AUTH_FILE` | `OF_AUTH_FILE` | Same |
| `OF_CONFIG_FILE` | `OF_CONFIG_FILE` | Same |
| `OF_CONFIG_DIR` | `OF_CONFIG_DIR` | Same |
| `OF_LOG_LEVEL` | `OF_LOG_LEVEL` | Same |
| `OF_DISCORD_WEBHOOK` | `OF_DISCORD_WEBHOOK` | Same |

---

## Docker Migration

### Before (Python)

```dockerfile
FROM python:3.11
RUN pip install ofscraper
ENTRYPOINT ["ofscraper"]
```

### After (Go)

```dockerfile
FROM alpine:3.19
COPY gofscraper /usr/local/bin/
ENTRYPOINT ["gofscraper"]
```

Or use the provided `docker-compose.yml`:

```bash
docker-compose up -d
```

---

## Script Hook Compatibility

Script hooks receive the same environment variables:

| Variable | Python | Go | Same? |
|----------|--------|-----|-------|
| `OF_USERNAME` | Yes | Yes | Yes |
| `OF_ACTION` | Yes | Yes | Yes |
| `OF_MEDIA_COUNT` | Yes | Yes | Yes |
| `OF_POST_COUNT` | Yes | Yes | Yes |
| `OF_FILE_PATH` | Yes | Yes | Yes |

Your existing hook scripts will work without modification.

---

## Troubleshooting

### "Config file not found"

GoFScraper looks for config at `~/.config/gofscraper/config.json` by default. If your Python config was at a different path:

```bash
gofscraper scraper --config /path/to/your/config.json
```

### "Auth validation failed"

Check that all required fields are present in your auth.json:

```bash
cat ~/.config/gofscraper/auth.json | python -m json.tool
```

Required fields: `auth_cookie`, `auth_user_agent`, `auth_x_bc`, `auth_user_id`.

### "FFmpeg not found"

For DRM content, install FFmpeg and either:
- Add it to your PATH, or
- Set the path in config: `"binary_options": { "ffmpeg": "/path/to/ffmpeg" }`

### Database Locked

GoFScraper uses WAL mode for SQLite. If you get "database is locked" errors, make sure no other process (including the Python version) has the database open.

### Performance Differences

GoFScraper may be significantly faster due to true parallelism. If this causes rate limiting issues, reduce concurrency:

```json
{
  "performance_options": {
    "download_sems": 3
  }
}
```

---

## Running Both Versions

You can safely run both versions side-by-side as long as they don't access the same database files simultaneously. They share the same config and auth format, so you can use the same config directory.

```bash
# Python version
ofscraper --action download --usernames user1

# Go version (same config, same output)
gofscraper scraper -a download -u user1
```

---

## Uninstalling Python OF-Scraper

Once you've verified GoFScraper works correctly:

```bash
pip uninstall ofscraper

# Or if using venv
rm -rf ~/.venvs/ofscraper
```

Your config, auth, databases, and downloaded files are not affected by uninstalling.
