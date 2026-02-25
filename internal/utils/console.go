// =============================================================================
// FILE: internal/utils/console.go
// PURPOSE: Console output utilities. Provides formatted terminal output
//          functions for status messages, headers, and styled text that
//          respect the --quiet and --no-color flags. Ports Python
//          utils/console.py.
// =============================================================================

package utils

import (
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// Console output
// ---------------------------------------------------------------------------

// quiet suppresses non-essential console output when set to true.
var quiet bool

// SetQuiet enables or disables quiet mode for console output.
//
// Parameters:
//   - q: true to suppress output.
func SetQuiet(q bool) {
	quiet = q
}

// PrintHeader writes a formatted section header to stdout.
//
// Parameters:
//   - title: The header text to display.
func PrintHeader(title string) {
	if quiet {
		return
	}
	line := strings.Repeat("=", 60)
	fmt.Fprintln(os.Stdout, line)
	fmt.Fprintf(os.Stdout, "  %s\n", title)
	fmt.Fprintln(os.Stdout, line)
}

// PrintStatus writes a status line with a label and value.
//
// Parameters:
//   - label: The status label (left side).
//   - value: The status value (right side).
func PrintStatus(label, value string) {
	if quiet {
		return
	}
	fmt.Fprintf(os.Stdout, "  %-25s %s\n", label+":", value)
}

// PrintBullet writes a bullet-pointed line to stdout.
//
// Parameters:
//   - text: The text to display after the bullet.
func PrintBullet(text string) {
	if quiet {
		return
	}
	fmt.Fprintf(os.Stdout, "  - %s\n", text)
}

// PrintSeparator writes a horizontal line separator to stdout.
func PrintSeparator() {
	if quiet {
		return
	}
	fmt.Fprintln(os.Stdout, strings.Repeat("-", 60))
}

// PrintError writes an error message to stderr.
//
// Parameters:
//   - msg: The error message.
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
}

// PrintWarning writes a warning message to stderr.
//
// Parameters:
//   - msg: The warning message.
func PrintWarning(msg string) {
	if quiet {
		return
	}
	fmt.Fprintf(os.Stderr, "WARNING: %s\n", msg)
}
