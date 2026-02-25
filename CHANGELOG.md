# Changelog

All notable changes to GoFScraper will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/).

---

## [0.1.0] - 2026-02-25

### Added

- Complete Go rewrite of Python OF-Scraper with full feature parity
- **CLI**: 8 subcommands (`scraper`, `manual`, `msg_check`, `story_check`, `paid_check`, `post_check`, `metadata`, `db`) with 100+ flags via Cobra
- **API client**: Full OnlyFans API v2 integration with pagination, auth signing, and rate limiting
- **Authentication**: SHA1-based request signing with multiple dynamic rule providers (digitalcriminals, datawhores, xagler, rafa, manual, generic)
- **Download system**: Concurrent downloads with resume support, bandwidth throttling, and .part file handling
- **DRM support**: DASH/MPD manifest parsing, Widevine CDM integration (manual + CDRM service modes), FFmpeg decryption
- **Filtering engine**: 29 composable filters for media (type, date, size, duration, URL, viewable, dedup), posts (date, text regex, mass message, ads, timed), and models (price, date, flags, subscription type)
- **Database**: SQLite with 10 tables, WAL mode, schema migration, backup, and merge operations
- **TUI**: Interactive terminal UI with Bubbletea (filterable tables, progress bars, live task display, sidebar filters)
- **Configuration**: JSON config with environment variable overrides, multiple profiles, 28 env var categories
- **Logging**: Structured slog logging with colored terminal output (tint), file handler with rotation, Discord webhook handler, sensitive data redaction
- **Script hooks**: 6 hook points (after-action, after-download, after-like, naming, skip-check, final) with environment variable context
- **Daemon mode**: Scheduled recurring scrapes with configurable interval
- **Interactive menus**: Main menu with scraper, config, profile, auth, and exit options
- **Worker pool**: Generic type-safe concurrent processing with Go generics
- **Path management**: Template-based path resolution with 15+ variables, sanitization, and truncation
- **Caching**: SQLite, JSON, and in-memory cache backends
- **File hashing**: XXHash128-based deduplication
- **Discord integration**: Webhook notifications for scrape completion
- **Docker support**: Multi-stage Dockerfile, docker-compose.yml with volume mounts
- **Release automation**: GoReleaser config for 6 platform targets (Linux/macOS/Windows x amd64/arm64)
- **Full documentation**: README, configuration reference, architecture guide, development guide, CLI reference, migration guide from Python, contributing guidelines

### Technical Details

- 318 Go source files across 44 packages
- 28,000+ lines of Go code
- ~3.6 MB compiled binary (zero runtime dependencies)
- Pure-Go SQLite via modernc.org/sqlite (no CGO required)
- Clean `go build` and `go vet` on all packages
