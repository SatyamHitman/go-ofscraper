# Architecture

GoFScraper is organized as a modular Go application with clear package boundaries, dependency injection, and a unidirectional data flow from CLI input to disk output.

---

## High-Level Data Flow

```
CLI Input (cobra)
    |
    v
App Orchestrator (internal/app)
    |
    +---> Config (internal/config)
    +---> Auth (internal/auth)
    +---> HTTP Session (internal/http)
    |
    v
Command Dispatch (internal/commands)
    |
    v
API Client (internal/api)  <---->  OF API (onlyfans.com/api2/v2)
    |
    v
Filters (internal/filter)
    |
    v
Download System (internal/download)
    +---> DRM Decrypt (internal/drm)
    +---> FFmpeg (external)
    |
    v
Database (internal/db)  <---->  SQLite
    |
    v
TUI / Progress (internal/tui)
    |
    v
Disk (media files)
```

---

## Package Dependency Graph

```
cmd/gofscraper/main.go
  └── internal/cli           CLI framework (cobra commands)
       ├── flags/            Flag definitions
       ├── bundles/          Flag grouping
       ├── accessors/        Flag value readers
       ├── callbacks/        Validation callbacks
       ├── mutators/         Pre-run transformations
       └── types/            Custom pflag.Value types

internal/app                 Application lifecycle
  ├── internal/config        Configuration (viper-like JSON)
  │    └── env/              Environment variable loading
  ├── internal/auth          Authentication
  │    └── providers/        Auth signing rule providers
  ├── internal/http          HTTP session management
  ├── internal/logging       Structured logging (slog)
  └── internal/db            SQLite database

internal/commands            Command implementations
  ├── scraper/               Main scraper pipeline
  ├── metadata/              Metadata update flow
  └── utils/                 Shared command utilities

internal/api                 OF API client
  └── internal/model         Domain models (Post, Media, User)

internal/filter              Content filtering engine
internal/download            Download orchestration
  ├── progress/              Progress tracking
  └── internal/drm           DRM decryption

internal/tui                 Terminal UI
  ├── fields/                Filter input widgets
  ├── inputs/                Raw input components
  ├── sections/              Layout sections
  ├── live/                  Live display management
  ├── live_classes/          Display data classes
  └── utils/                 TUI utilities

internal/prompts             Interactive prompts
internal/scripts             External script hooks
internal/worker              Generic worker pool
internal/hash                XXHash file deduplication
internal/paths               Path resolution and sanitization
internal/cache               Caching layer
internal/utils               General utilities
  └── system/                OS-level utilities

pkg/version                  Version info (public)
```

---

## Package Descriptions

### `cmd/gofscraper`

Entry point. Calls `cli.Execute()`. Nothing else.

### `internal/app`

Application lifecycle orchestrator. Key types:

| Type | Responsibility |
|------|---------------|
| `App` | Holds context, config, session, logger. Init/Shutdown lifecycle. |
| `Manager` | Singleton managing shared state across operations |
| `ModelManager` | Processes a single user: fetch areas, filter, dispatch actions |
| `PostCollection` | Aggregates posts from multiple content areas |
| `Stats` | Atomic counters for tracking all metrics |
| `State` | Current processing phase and progress |
| `DaemonConfig` | Scheduler configuration for daemon mode |

### `internal/api`

OF API client. Built around a `ContentFetcher` pattern:

- **Endpoint methods**: `GetTimeline`, `GetMessages`, `GetStories`, `GetHighlights`, `GetPinned`, `GetArchived`, `GetStreams`, `GetLabels`, `GetPurchased`, `GetSubscriptions`, `GetProfile`, `GetMe`, `PostFavorite`
- **Pagination**: Cursor-based with `afterPublishTime` / `tailMarker`
- **Response parsing**: `parsePost()` and `parseMedia()` in `common.go`

### `internal/auth`

Authentication subsystem:

```
auth.Load(path) -> ReadAuthFile -> Validate -> Set globally
                                                    |
                                                    v
                                      signing.CreateSign(url, headers, params)
                                                    |
                                                    v
                                        SHA1(static_param + time + path + user_id)
                                                    |
                                                    v
                                              headers["sign"] = result
```

**Providers** supply signing parameters (`static_param`, `checksum_indexes`, `checksum_constant`, `format`) from different sources.

### `internal/http`

HTTP session management:

| Type | Responsibility |
|------|---------------|
| `SessionManager` | Wraps `http.Client` with auth, rate limiting, retry |
| `Request` | Request builder with action flags |
| `Response` | Response wrapper with JSON decode helpers |
| `RateLimiter` | Token bucket rate limiter |
| `AdaptiveSleeper` | Exponential backoff on 429/403 with time-based decay |

### `internal/download`

Download pipeline:

```
Orchestrator.Run(ctx, media)
    |
    +---> per media item via worker channel
    |
    +---> downloadOne(media)
           |
           +---> Normal: HTTP GET with Range header, .part file, SpeedLimitReader
           |
           +---> Protected: DRM pipeline (DASH parse -> key fetch -> FFmpeg decrypt)
```

| Type | Responsibility |
|------|---------------|
| `Orchestrator` | Dispatches download workers via channel |
| `SpeedLimitReader` | io.Reader wrapper enforcing bytes/sec limit |
| `RetryPolicy` | Configurable retry with exponential backoff |

### `internal/drm`

DRM decryption:

```
Manager.Decrypt(ctx, mpdURL, outputPath)
    |
    +---> ParseMPD(url) -> Manifest
    |         |
    |         +---> KeyID(), PSSH()
    |
    +---> GetKey (mode-dependent)
    |         |
    |         +---> "cdrm": CDRMClient.GetKey(pssh, licenseURL)
    |         +---> "manual": ManualDecrypt(device, licenseURL, pssh)
    |
    +---> FFmpegDecrypt(input, output, kid:key)
```

### `internal/filter`

Composable filter chains using function types:

```go
type MediaFilter func([]model.Media) []model.Media
type PostFilter  func([]model.Post) []model.Post
type ModelFilter func([]model.User) []model.User

// Usage:
filtered := filter.ChainMedia(
    filter.ByMediaType("videos"),
    filter.ByMediaDate(after, before),
    filter.ByURLPresence(),
    filter.ByMediaDupe(),
)(allMedia)
```

Each filter returns `nil` to skip (no-op in chain).

### `internal/model`

Domain types:

| Type | Description |
|------|-------------|
| `Post` | Content post with text, media list, dates, metadata |
| `Media` | Single downloadable item (image/video/audio) with DRM fields |
| `User` | Creator profile with subscription info, pricing, flags |
| `Label` | Content label/category |

`Media` contains a `sync.Mutex` for thread-safe status updates during concurrent downloads.

### `internal/db`

SQLite database layer:

- **10 tables**: posts, messages, medias, stories, labels, others, products, profiles, models, schema_flags
- **Connection pool**: One connection per username, cached in map
- **WAL mode**: Write-Ahead Logging for concurrent read access
- **Schema migration**: `transition.go` handles upgrades via `schema_flags`

### `internal/tui`

Terminal UI built on Bubbletea:

```
App (bubbletea.Model)
  |
  +---> View: Menu | Table | Progress
  |
  +---> Sections
  |       +---> Sidebar (filter fields)
  |       +---> TableSection (data + pagination + sorting)
  |       +---> ConsoleSection (plain-text fallback)
  |
  +---> Live Display
          +---> ProgressBar
          +---> TaskDisplay
          +---> ScreenManager
```

### `internal/worker`

Generic worker pool with Go type parameters:

```go
pool := worker.NewPool[model.Media, DownloadResult](8, downloadFunc)
results := pool.Run(ctx, mediaItems)
```

---

## Concurrency Model

| Pattern | Go Implementation |
|---------|-------------------|
| Concurrent downloads | N goroutines reading from `chan model.Media` |
| API pagination | `errgroup` for parallel page fetching |
| Rate limiting | `x/time/rate` token bucket |
| Adaptive backoff | `AdaptiveSleeper` with time-decaying multiplier |
| Shared state | `sync.Mutex` / `sync.RWMutex` on all mutable globals |
| Cancellation | `context.Context` propagated through all operations |
| Signal handling | `os/signal.Notify` -> `context.CancelFunc` |

---

## Database Schema

```sql
-- Core content tables
posts       (id, post_id, text, price, paid, created_at, model_id)
messages    (id, post_id, text, price, paid, created_at, model_id)
medias      (id, media_id, post_id, link, directory, filename, size, media_type,
             duration, downloaded, created_at, model_id, hash)
stories     (id, post_id, text, price, paid, created_at, model_id)
labels      (id, label_id, name, type, post_id, model_id)

-- Supporting tables
others      (id, post_id, text, price, paid, created_at, model_id)
products    (id, post_id, text, price, paid, created_at, model_id)
profiles    (id, user_id, username, model_id)
models      (id, model_id, username)
schema_flags (id, flag_name, flag_value)
```

Each user gets their own SQLite database file at the configured metadata path.

---

## Error Handling

Structured error types with codes:

```go
type AppError struct {
    Code    ErrorCode
    Message string
    Err     error
}
```

| Code | HTTP Status | Meaning |
|------|-------------|---------|
| `ErrAuth` | 401, 400 | Authentication failure |
| `ErrRateLimit` | 429, 504 | Rate limited / gateway timeout |
| `ErrForbidden` | 403 | Access denied |
| `ErrNotFound` | 404 | Resource not found |
| `ErrNetwork` | - | Connection failure |
| `ErrDatabase` | - | SQLite operation failure |
| `ErrFFmpeg` | - | FFmpeg execution failure |
| `ErrDRM` | - | DRM decryption failure |
| `ErrConfig` | - | Configuration error |

Individual download failures do not abort the batch. Errors are collected and reported in the final summary.
