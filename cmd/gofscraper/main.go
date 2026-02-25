// =============================================================================
// FILE: cmd/gofscraper/main.go
// PURPOSE: Entry point for the gofscraper application. Initializes the root
//          CLI command and launches the application lifecycle. This is the
//          equivalent of Python OF-Scraper's __main__.py.
// =============================================================================

package main

import (
	"gofscraper/internal/cli"
)

// main is the application entry point. It initializes the CLI framework,
// parses arguments, and delegates to the appropriate command handler.
//
// Exit codes:
//   - 0: Successful execution
//   - 1: General error (CLI parse failure, runtime error)
func main() {
	cli.Execute()
}
