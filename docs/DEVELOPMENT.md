# Development Guide

This document covers how to set up the development environment, build the project, run tests, and contribute code.

---

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.23+ | Compiler and toolchain |
| Git | 2.x+ | Version control |
| FFmpeg | 4.x+ | DRM content decryption (runtime only) |
| Make | any | Build automation |
| golangci-lint | latest | Linting (optional) |
| sqlc | 1.25+ | SQL code generation (optional, only if modifying queries) |
| GoReleaser | 2.x | Release builds (optional) |

### Installing Go

```bash
# Linux/macOS
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Windows (use installer)
# https://go.dev/dl/go1.23.6.windows-amd64.msi

# Verify
go version
```

### Installing FFmpeg

```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# macOS
brew install ffmpeg

# Windows (download from https://ffmpeg.org/download.html)
# Add to PATH
```

---

## Getting Started

```bash
# Clone the repository
git clone https://github.com/your-org/gofscraper.git
cd gofscraper

# Download dependencies
go mod download

# Build
make build

# Run tests
make test

# Run the binary
./bin/gofscraper --help
```

---

## Build Commands

```bash
# Standard build (with version info)
make build

# Quick build (no version ldflags)
go build -o bin/gofscraper ./cmd/gofscraper

# Build and run
make run

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o bin/gofscraper-linux ./cmd/gofscraper
GOOS=darwin GOARCH=arm64 go build -o bin/gofscraper-macos ./cmd/gofscraper
GOOS=windows GOARCH=amd64 go build -o bin/gofscraper.exe ./cmd/gofscraper

# Release build (all platforms)
goreleaser release --snapshot --clean
```

---

## Project Layout

```
gofscraper/
  cmd/gofscraper/          Entry point (main.go only)
  internal/                Private application code
    app/                   Application lifecycle and orchestration
    api/                   OnlyFans API client
    auth/                  Authentication and signing
      providers/           Auth rule providers
    cache/                 Caching layer (sqlite, json, memory)
    cli/                   CLI framework (cobra commands)
      accessors/           Flag value readers
      bundles/             Flag grouping
      callbacks/           Validation functions
      flags/               Flag definitions
      mutators/            Pre-run transformations
      types/               Custom pflag types
    commands/              Command implementations
      scraper/             Main scraper pipeline
      metadata/            Metadata updates
      utils/               Shared command utilities
    config/                Configuration management
      env/                 Environment variable loading
    db/                    SQLite database operations
      queries/             SQL query files (for sqlc)
      sqlc/                Auto-generated query code
    download/              Download system
      progress/            Progress tracking
    drm/                   DRM decryption
    filter/                Content filtering
    hash/                  File hashing (xxh3)
    http/                  HTTP session management
    logging/               Structured logging
    model/                 Domain models
    paths/                 Path resolution
    posts/                 Post processing
    prompts/               Interactive prompts
      groups/              Prompt group definitions
      utils/               Prompt helpers
    scripts/               External script hooks
    tui/                   Terminal UI
      fields/              Filter input widgets
      inputs/              Raw input components
      live/                Live display management
      live_classes/        Display data classes
      sections/            Layout sections
      utils/               TUI utilities
    utils/                 General utilities
      system/              OS-level utilities
    worker/                Generic worker pool
  pkg/version/             Public version info
  docs/                    Documentation
```

### Key Design Principles

1. **`internal/` for everything** -- All application code is in `internal/` to prevent external imports. Only `pkg/version` is public.

2. **One concern per file** -- Each file handles a single responsibility. File names match their contents (`media_type.go` contains `ByMediaType`).

3. **Interfaces at boundaries** -- Package boundaries use interfaces (`ContentFetcher`, `Cache`, `ScriptRunner`) for testability.

4. **No global mutable state without sync** -- All package-level variables are protected by `sync.Mutex` or `sync.RWMutex`.

5. **Context everywhere** -- Every long-running operation accepts `context.Context` for cancellation.

---

## Adding a New Feature

### Adding a New CLI Flag

1. Define the flag in `internal/cli/flags/<category>.go`:
   ```go
   func RegisterMyFlags(cmd *cobra.Command) {
       cmd.Flags().String("my-flag", "default", "Description")
   }
   ```

2. Add an accessor in `internal/cli/accessors/`:
   ```go
   func GetMyFlag(cmd *cobra.Command) string {
       val, _ := cmd.Flags().GetString("my-flag")
       return val
   }
   ```

3. Register in the appropriate bundle (`internal/cli/bundles/`):
   ```go
   func RegisterMyBundle(cmd *cobra.Command) {
       flags.RegisterMyFlags(cmd)
   }
   ```

4. Call the bundle in the command's `init()`.

### Adding a New Filter

1. Create `internal/filter/my_filter.go`:
   ```go
   package filter

   import "gofscraper/internal/model"

   func ByMyCondition(param string) MediaFilter {
       if param == "" {
           return nil // nil = skip in chain
       }
       return func(items []model.Media) []model.Media {
           var out []model.Media
           for _, m := range items {
               if matchesCondition(m, param) {
                   out = append(out, m)
               }
           }
           return out
       }
   }
   ```

2. Add to the filter chain in the scraper pipeline.

### Adding a New API Endpoint

1. Add the URL template in `internal/api/endpoints.go`
2. Create a new file `internal/api/myendpoint.go`
3. Implement the fetcher method on the `Client` struct
4. Add to the `ContentFetcher` interface if applicable

### Adding a New Content Area

1. Add the area name to `internal/prompts/strings.go`
2. Add the API fetcher method
3. Add the response type mapping in config schema
4. Update the scraper to fetch from the new area

---

## Database Development

### Modifying the Schema

1. Edit `internal/db/queries/schema.sql`
2. Add migration logic in `internal/db/transition.go`
3. Update or create query files in `internal/db/queries/`
4. Run `make generate` (or `sqlc generate`) to regenerate Go code
5. Update `internal/db/operations.go` to use the new queries

### Query Files

SQL queries follow the [sqlc](https://sqlc.dev/) convention:

```sql
-- name: GetMediaByPostID :many
SELECT * FROM medias WHERE post_id = ? ORDER BY media_id;

-- name: InsertMedia :exec
INSERT INTO medias (media_id, post_id, link, media_type, model_id)
VALUES (?, ?, ?, ?, ?);
```

---

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./internal/filter/... -v

# Run a specific test
go test ./internal/filter/ -run TestByMediaType -v

# Run with race detector
go test ./... -race

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Writing Tests

Tests live alongside the code they test:

```
internal/filter/
  media_type.go
  media_type_test.go     <-- tests for media_type.go
```

Or in the `test/` directory for integration tests:

```
test/
  db/
  general/
  post/
  calls/
```

---

## Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint

# Run go vet only
go vet ./...
```

---

## Code Style

- **Formatting**: `gofmt` / `goimports` (enforced)
- **Naming**: Standard Go conventions (`CamelCase` exports, `camelCase` unexported)
- **Comments**: Godoc-style on all exported types and functions
- **Error handling**: Always wrap errors with context: `fmt.Errorf("operation: %w", err)`
- **Imports**: Grouped as stdlib, external, internal (separated by blank lines)
- **File headers**: Each file has a header comment block with FILE, PURPOSE

---

## Releasing

### Manual Release

```bash
# Tag a version
git tag v1.0.0
git push origin v1.0.0

# Build release artifacts
goreleaser release --clean
```

### CI Release

Push a tag matching `v*` to trigger the release workflow. GoReleaser builds binaries for all 6 platform targets and creates a GitHub release with checksums.

### Docker Release

```bash
docker build -t gofscraper:latest .
docker tag gofscraper:latest ghcr.io/your-org/gofscraper:latest
docker push ghcr.io/your-org/gofscraper:latest
```

---

## Debugging

### Verbose Logging

```bash
# Set log level to debug
gofscraper scraper --log-level debug

# Or trace for maximum detail
gofscraper scraper --log-level trace
```

### Environment Variable Overrides

```bash
# Override specific settings without touching config
OF_LOG_LEVEL=DEBUG OF_DOWNLOAD_DIR=/tmp/test gofscraper scraper
```

### SQLite Inspection

```bash
# Database files are at the metadata path
sqlite3 ~/.config/gofscraper/main_profile/.data/<model_id>/db.sqlite

# Useful queries
.tables
SELECT COUNT(*) FROM medias;
SELECT * FROM medias WHERE downloaded = 0 LIMIT 10;
```

---

## Dependency Management

```bash
# Add a new dependency
go get github.com/some/package@latest

# Update all dependencies
go get -u ./...

# Tidy (remove unused, add missing)
go mod tidy

# Verify checksums
go mod verify
```

### Key Dependencies

| Package | Purpose | Update Frequency |
|---------|---------|------------------|
| `spf13/cobra` | CLI framework | Stable, rare updates |
| `modernc.org/sqlite` | Pure-Go SQLite | Monthly patches |
| `charmbracelet/bubbletea` | TUI framework | Active development |
| `charmbracelet/lipgloss` | TUI styling | Active development |
| `zeebo/xxh3` | Fast hashing | Stable |
| `dustin/go-humanize` | Size formatting | Stable, rare updates |
| `lmittmann/tint` | Colored slog | Stable |
| `google/uuid` | UUID generation | Stable |
