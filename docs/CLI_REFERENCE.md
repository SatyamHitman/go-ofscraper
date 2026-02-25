# CLI Reference

Complete reference for all GoFScraper commands and flags.

---

## Global Flags

These flags are available on all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `""` | Path to config file |
| `--profile` | `-p` | `"default"` | Configuration profile name |
| `--log-level` | `-l` | `"info"` | Log level: `trace`, `debug`, `info`, `warn`, `error` |
| `--no-interactive` | `-n` | `false` | Disable interactive prompts |
| `--verbose` | `-v` | `false` | Enable verbose output |
| `--version` | | | Show version and exit |
| `--help` | `-h` | | Show help |

---

## scraper

Main scraper command. Downloads media, likes/unlikes posts, and updates metadata.

```bash
gofscraper scraper [flags]
```

### Scraper Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--action` | `-a` | `download` | Actions: `download`, `like`, `unlike` (comma-separated) |
| `--posts` | `-o` | all | Content areas to process (see below) |
| `--users` | `-u` | `""` | Usernames (comma-separated) |
| `--excluded-users` | | `""` | Usernames to exclude (comma-separated) |
| `--daemon` | `-d` | `false` | Run in daemon mode with scheduled repeats |

### Content Areas

Available values for `--posts`:

| Area | Description |
|------|-------------|
| `timeline` | Timeline/feed posts |
| `messages` | Direct messages |
| `archived` | Archived posts |
| `stories` | Stories |
| `highlights` | Highlights |
| `pinned` | Pinned posts |
| `streams` | Streams |
| `purchased` | Purchased content |
| `labels` | Labelled posts |

Use `all` to select all areas, or comma-separate specific ones:

```bash
gofscraper scraper -o timeline,messages,stories
```

### Download Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--arrow` | `-ar` | `""` | Download direction |
| `--database` | `-db` | `""` | Database path |
| `--save-dir` | `-sd` | `""` | Override save directory |
| `--download-limit` | `-dl` | `0` | Bandwidth limit (bytes/sec) |

### Media Filter Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--quality` | `-q` | `""` | Quality preference |
| `--media-type` | `-mt` | `""` | Filter by type: `images`, `videos`, `audios` |
| `--size-max` | `-sx` | `0` | Max file size (bytes) |
| `--size-min` | `-sm` | `0` | Min file size (bytes) |
| `--media-id` | `-mid` | `""` | Filter by media ID range (`min-max`) |
| `--length-max` | `-lx` | `0` | Max duration (seconds) |
| `--length-min` | `-lm` | `0` | Min duration (seconds) |
| `--max-count` | `-mxc` | `0` | Max items to process (0 = unlimited) |
| `--media-sort` | `-mst` | `""` | Sort by: `date`, `id`, `type`, `size` |
| `--media-desc` | `-mdc` | `false` | Sort descending |
| `--excluded` | `-e` | `""` | Excluded media types |
| `--excluded-quality` | `-eq` | `""` | Excluded quality levels |
| `--before-epoch` | `-be` | `""` | Filter media before date |
| `--after` | `-af` | `""` | Filter media after date |
| `--remove-duplicates` | `-rd` | `false` | Remove duplicate media |

### Post Filter Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--only` | `-o` | `""` | Filter mode |
| `--date-after` | `-da` | `""` | Posts after date (YYYY-MM-DD) |
| `--label-after` | `-la` | `""` | Posts with label after date |
| `--filter-text` | `-ft` | `""` | Keep posts matching regex |
| `--neg-filter` | `-nf` | `""` | Remove posts matching regex |
| `--before-after` | `-ba` | `""` | Posts before date |
| `--timed` | `-t` | `""` | Timed post filter: `only`, `exclude` |
| `--skip-pinned` | `-sp` | `false` | Skip pinned posts |
| `--mass-only` | | `false` | Only mass messages |
| `--mass-exclude` | | `false` | Exclude mass messages |

### User Selection Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--usernames` | `-u` | `""` | Usernames to process |
| `--excluded-users` | `-eu` | `""` | Usernames to exclude |
| `--user-list` | `-ul` | `""` | User list name(s) |
| `--blacklist` | `-bl` | `""` | Blacklist name(s) |
| `--sort-type` | `-st` | `""` | Sort users by: `name`, `subscribed`, `expired`, `price` |
| `--descending-sort` | `-ds` | `false` | Sort descending |

### Advanced User Filters

| Flag | Default | Description |
|------|---------|-------------|
| `--min-price` | `0` | Minimum subscription price |
| `--max-price` | `0` | Maximum subscription price |
| `--last-seen-after` | `""` | User last seen after date |
| `--last-seen-before` | `""` | User last seen before date |
| `--sub-after` | `""` | Subscribed after date |
| `--sub-before` | `""` | Subscribed before date |
| `--expired-after` | `""` | Expired after date |
| `--expired-before` | `""` | Expired before date |
| `--renewal-on` | `false` | Only users with renewal on |
| `--renewal-off` | `false` | Only users with renewal off |
| `--promo-only` | `false` | Only users with active promo |
| `--free-only` | `false` | Only free accounts |
| `--regular-only` | `false` | Only paid accounts |

### File Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file-format` | `-g` | `""` | Filename template override |
| `--text-type` | `-tt` | `""` | Text truncation: `letter` or `word` |
| `--space-replacer` | `-sr` | `""` | Space replacement character |
| `--text-length` | `-tl` | `0` | Max text length in filenames |

### Script Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--after-action-script` | `""` | Script after action completion |
| `--post-script` | `""` | Script after each post |
| `--naming-script` | `""` | Custom filename script |
| `--after-dl-script` | `""` | Script after each download |
| `--skip-dl-script` | `""` | Script to check skip condition |

### Advanced Program Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--no-cache` | `false` | Disable all caching |
| `--no-cache-api` | `false` | Disable API response caching |
| `--key-mode` | `""` | CDM key mode: `cdrm`, `manual` |
| `--private-key` | `""` | Widevine private key path |
| `--check-interval` | `0` | Daemon check interval (minutes) |
| `--auto-like` | `false` | Automatically like downloaded posts |
| `--dynamic-rules` | `""` | Auth rule provider override |
| `--update-profile` | `false` | Update profile data |

---

## manual

Download content from direct URLs.

```bash
gofscraper manual <url1> [url2] [url3] ... [flags]
```

Supports the common, download, file, and media flags.

---

## msg_check

Inspect and list messages for a user.

```bash
gofscraper msg_check -u username [flags]
```

### Check-Specific Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--check-area` | `-ca` | `""` | Content area to check |
| `--user` | `-u` | `""` | Username to check |
| `--file` | `-f` | `""` | Output file path |
| `--force` | `-fo` | `false` | Force refresh (ignore cache) |
| `--table-progress` | `-tp` | `false` | Show table progress |
| `--table-name` | `-tn` | `""` | Custom table name |

---

## story_check

Inspect and list stories for a user.

```bash
gofscraper story_check -u username [flags]
```

Same flags as `msg_check`.

---

## paid_check

Inspect and list purchased content.

```bash
gofscraper paid_check -u username [flags]
```

Same flags as `msg_check` plus media and post filter flags.

---

## post_check

Inspect and list posts for a user.

```bash
gofscraper post_check -u username [flags]
```

Same flags as `msg_check` plus media, post filter, and download flags.

---

## metadata

Update metadata for downloaded content.

```bash
gofscraper metadata -u username [flags]
```

### Metadata Filter Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--update-cache-list` | `false` | Update cache list |
| `--download-list` | `false` | Update download list |
| `--protected-filter` | `false` | Filter protected content |
| `--protected-bypass` | `false` | Bypass protected check |
| `--cache-filter` | `false` | Filter cached content |
| `--cache-bypass` | `false` | Bypass cache check |
| `--preview` | `false` | Preview mode (no writes) |

---

## db

Database management operations.

```bash
gofscraper db [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--backup` | `false` | Create database backup |
| `--merge` | `false` | Merge two databases |

### Examples

```bash
# Backup all user databases
gofscraper db --backup

# Merge databases
gofscraper db --merge -u source_user,dest_user
```

---

## Usage Examples

### Basic Download

```bash
# Download all content from a user
gofscraper scraper -u janedoe -a download

# Download only videos from timeline
gofscraper scraper -u janedoe -o timeline --media-type videos

# Download from multiple users
gofscraper scraper -u user1,user2,user3

# Download with bandwidth limit (1 MB/s)
gofscraper scraper -u janedoe --download-limit 1048576
```

### Filtering

```bash
# Only videos longer than 5 minutes
gofscraper scraper -u janedoe --media-type videos --length-min 300

# Posts from the last 30 days
gofscraper scraper -u janedoe --date-after 2024-12-01

# Skip mass messages
gofscraper scraper -u janedoe --mass-exclude

# Only posts containing "exclusive"
gofscraper scraper -u janedoe --filter-text "exclusive"

# Files larger than 10 MB
gofscraper scraper -u janedoe --size-min 10485760
```

### Daemon Mode

```bash
# Run every 6 hours
gofscraper scraper -d --check-interval 360

# Daemon with specific users
gofscraper scraper -d -u user1,user2 --check-interval 120
```

### Non-Interactive Mode

```bash
# CI/server usage (no prompts)
gofscraper scraper -n -u janedoe -a download -o timeline,messages
```

### Like/Unlike

```bash
# Like all posts from a user
gofscraper scraper -u janedoe -a like

# Unlike all posts
gofscraper scraper -u janedoe -a unlike
```

### User Filtering

```bash
# Only free accounts
gofscraper scraper --free-only

# Users with active promo under $10
gofscraper scraper --promo-only --max-price 10

# Users last seen in the past week
gofscraper scraper --last-seen-after 2024-12-18
```
