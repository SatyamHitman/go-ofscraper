// =============================================================================
// FILE: internal/utils/checkers.go
// PURPOSE: Various checker utilities used across the application. Includes
//          file existence checks, network availability tests, and
//          pre-condition validators. Ports Python utils/checkers.py.
// =============================================================================

package utils

import (
	"net"
	"os"
	"time"
)

// ---------------------------------------------------------------------------
// File checks
// ---------------------------------------------------------------------------

// FileExists reports whether a file exists at the given path.
//
// Parameters:
//   - path: The file path to check.
//
// Returns:
//   - true if the path exists and is not a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists reports whether a directory exists at the given path.
//
// Parameters:
//   - path: The directory path to check.
//
// Returns:
//   - true if the path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileSize returns the size of a file in bytes, or 0 if the file doesn't exist.
//
// Parameters:
//   - path: The file path.
//
// Returns:
//   - File size in bytes, or 0 on error.
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// ---------------------------------------------------------------------------
// Network checks
// ---------------------------------------------------------------------------

// IsOnline checks basic network connectivity by attempting a TCP connection
// to a well-known host.
//
// Returns:
//   - true if a TCP connection can be established.
func IsOnline() bool {
	conn, err := net.DialTimeout("tcp", "1.1.1.1:443", 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ---------------------------------------------------------------------------
// Value checks
// ---------------------------------------------------------------------------

// IsEmpty reports whether a string is empty or only whitespace.
//
// Parameters:
//   - s: The string to check.
//
// Returns:
//   - true if s is empty or blank.
func IsEmpty(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}

// CoalesceString returns the first non-empty string from the arguments.
//
// Parameters:
//   - values: Candidate strings in priority order.
//
// Returns:
//   - The first non-empty string, or "" if all are empty.
func CoalesceString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// CoalesceInt returns the first non-zero int from the arguments.
//
// Parameters:
//   - values: Candidate ints in priority order.
//
// Returns:
//   - The first non-zero value, or 0 if all are zero.
func CoalesceInt(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}
