# GoFScraper

A high-performance OnlyFans content downloader and manager written in Go. Full-feature port of the Python [OF-Scraper](https://github.com/datawhores/OF-Scraper) project, rebuilt from the ground up for speed, reliability, and ease of deployment.

## Features

- **Multi-action pipeline** -- Download media, like/unlike posts, update metadata
- **All content areas** -- Timeline, messages, stories, highlights, archived, pinned, streams, purchased, labels
- **DRM support** -- DASH/MPD manifest parsing with Widevine CDM and CDRM service integration
- **Interactive TUI** -- Terminal UI with filterable tables, progress bars, and live status
- **Daemon mode** -- Scheduled recurring scrapes on a configurable interval
- **Multiple profiles** -- Switch between configurations for different accounts
- **Advanced filtering** -- Filter by media type, date, size, duration, price, text regex, and more
- **Resume support** -- Interrupted downloads resume from where they left off
- **Rate limiting** -- Adaptive request throttling with exponential backoff on 429/403
- **Script hooks** -- Run custom scripts after downloads, for naming, skip logic, and cleanup
- **Discord notifications** -- Webhook integration for scrape completion alerts
- **Cross-platform** -- Native binaries for Linux, macOS, and Windows (amd64 + arm64)
- **Single binary** -- No Python runtime, no pip, no virtual environments

## Quick Start

### Download a Release

Grab the latest binary from the [Releases](https://github.com/your-org/gofscraper/releases) page for your platform.

### Build from Source

**Prerequisites:** Go 1.23+ and FFmpeg (for DRM content)

```bash
git clone https://github.com/your-org/gofscraper.git
cd gofscraper
make build
```

The binary is output to `bin/gofscraper`.

### First Run

```bash
# Show all commands and flags
gofscraper --help

# Run the scraper (interactive mode -- prompts for action and users)
gofscraper scraper

# Download from specific users
gofscraper scraper -u username1,username2 -a download

# Check messages for a user
gofscraper msg_check -u username1
```

### Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or build manually
docker build -t gofscraper .
docker run -v ./config:/root/.config/gofscraper -v ./data:/root/Data/ofscraper gofscraper scraper
```

## Authentication

GoFScraper requires OnlyFans authentication credentials. Create an auth file at `~/.config/gofscraper/auth.json`:

```json
{
  "auth_cookie": "your-sess-cookie-value",
  "auth_user_agent": "your-browser-user-agent",
  "auth_x_bc": "your-x-bc-token",
  "auth_user_id": "your-numeric-user-id",
  "auth_app_token": "33d57ade8c02dbc5a333db99ff9ae26a"
}
```

**How to get these values:**

1. Log into OnlyFans in your browser
2. Open DevTools (F12) > Network tab
3. Navigate to any page and find a request to `onlyfans.com/api2/v2/`
4. From the request headers, copy:
   - `cookie` (the `sess=` value) -> `auth_cookie`
   - `user-agent` -> `auth_user_agent`
   - `x-bc` -> `auth_x_bc`
   - `user-id` -> `auth_user_id`

Or use the interactive auth setup:

```bash
gofscraper scraper
# Select "Auth" from the main menu
```

## Configuration

Config file location: `~/.config/gofscraper/config.json`

A default config is created on first run. See [docs/CONFIGURATION.md](docs/CONFIGURATION.md) for the full reference.

**Minimal example:**

```json
{
  "file_options": {
    "save_location": "/home/user/Data/ofscraper",
    "dir_format": "{model_username}/{responsetype}/{mediatype}/",
    "file_format": "{filename}.{ext}"
  },
  "download_options": {
    "filter": ["Images", "Videos", "Audios"],
    "auto_resume": true
  },
  "performance_options": {
    "download_sems": 6
  }
}
```

## Commands

| Command | Description |
|---------|-------------|
| `scraper` | Main scraper -- download, like, unlike content |
| `manual` | Download from direct URLs |
| `msg_check` | Inspect and list messages |
| `story_check` | Inspect and list stories |
| `paid_check` | Inspect and list purchased content |
| `post_check` | Inspect and list posts |
| `metadata` | Update metadata for downloaded content |
| `db` | Database backup and merge operations |

Run `gofscraper <command> --help` for detailed flag information.

## Documentation

| Document | Description |
|----------|-------------|
| [CONFIGURATION.md](docs/CONFIGURATION.md) | Full config file reference with all options |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | System design, package layout, data flow |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md) | Build setup, code style, contributing guide |
| [CLI_REFERENCE.md](docs/CLI_REFERENCE.md) | Complete CLI flags and usage examples |
| [MIGRATION.md](docs/MIGRATION.md) | Migrating from Python OF-Scraper |

## Project Stats

| Metric | Value |
|--------|-------|
| Go source files | 318 |
| Lines of code | 28,000+ |
| Packages | 44 |
| Binary size | ~3.6 MB |
| External dependencies | 14 direct |
| Supported platforms | 6 (Linux/macOS/Windows x amd64/arm64) |

## License

This project is provided as-is for educational and personal use.
