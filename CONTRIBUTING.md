# Contributing to GoFScraper

Thank you for your interest in contributing. This document provides guidelines for contributing to the project.

---

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/gofscraper.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `make test`
6. Run linter: `make lint`
7. Commit: `git commit -m "feat: add my feature"`
8. Push: `git push origin feature/my-feature`
9. Open a Pull Request

---

## Development Setup

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for full setup instructions.

Quick start:

```bash
go mod download
make build
make test
```

---

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new media filter for resolution
fix: handle nil pointer in pagination
docs: update CLI reference for new flags
refactor: simplify download retry logic
test: add unit tests for auth signing
chore: update dependencies
```

**Types:** `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `ci`

---

## Code Guidelines

### General

- Run `go vet ./...` and `golangci-lint run ./...` before submitting
- All exported types and functions must have godoc comments
- Keep functions focused -- one function, one responsibility
- Wrap errors with context: `fmt.Errorf("download media %d: %w", id, err)`
- Use `context.Context` for all operations that could block or be cancelled

### File Organization

- One concern per file (e.g., `media_type.go` for `ByMediaType` filter)
- File names match their primary type or function
- Each file starts with a header comment block

### Testing

- Tests live in `*_test.go` files alongside the code they test
- Use table-driven tests for multiple cases
- Use `testing.T.Helper()` in test helpers
- Aim for meaningful tests, not coverage metrics

### Dependencies

- Prefer stdlib over external packages where practical
- Discuss new dependencies in the PR before adding them
- Keep `go.mod` tidy (`go mod tidy`)

---

## Pull Request Process

1. **Description** -- Clearly describe what the PR does and why
2. **Tests** -- Add or update tests for your changes
3. **Docs** -- Update documentation if adding features or changing behavior
4. **Single concern** -- One PR per feature or fix. Split large changes.
5. **Clean history** -- Squash fixup commits before requesting review

---

## Reporting Issues

When reporting bugs, please include:

- GoFScraper version (`gofscraper --version`)
- Operating system and architecture
- Steps to reproduce
- Expected vs actual behavior
- Relevant log output (use `--log-level debug`)

---

## Feature Requests

Open an issue with the `enhancement` label. Include:

- What problem the feature solves
- Proposed solution (if you have one)
- Whether you'd be willing to implement it
