// =============================================================================
// FILE: internal/utils/system/subprocess.go
// PURPOSE: Subprocess execution utilities. Provides helpers for running
//          external processes (FFmpeg, scripts) with timeouts, output capture,
//          and error handling. Ports Python utils/system/subprocess.py.
// =============================================================================

package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Subprocess execution
// ---------------------------------------------------------------------------

// RunResult holds the outcome of a subprocess execution.
type RunResult struct {
	Stdout   string        // Standard output
	Stderr   string        // Standard error
	ExitCode int           // Process exit code (0 = success)
	Duration time.Duration // How long the process ran
	Err      error         // Error if the process could not be started or timed out
}

// Success reports whether the process exited with code 0.
func (r RunResult) Success() bool {
	return r.ExitCode == 0 && r.Err == nil
}

// Run executes a command with the given arguments and a timeout.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - timeout: Maximum runtime. 0 means no timeout.
//   - name: The executable name or path.
//   - args: Command arguments.
//
// Returns:
//   - A RunResult with output and exit info.
func Run(ctx context.Context, timeout time.Duration, name string, args ...string) RunResult {
	start := time.Now()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	elapsed := time.Since(start)

	result := RunResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: elapsed,
	}

	if err != nil {
		result.Err = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// RunShell executes a shell command string. Uses "sh -c" on Unix and
// "cmd /C" on Windows.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - timeout: Maximum runtime.
//   - command: The shell command string.
//
// Returns:
//   - A RunResult with output and exit info.
func RunShell(ctx context.Context, timeout time.Duration, command string) RunResult {
	shell, flag := "sh", "-c"
	if IsWindows() {
		shell, flag = "cmd", "/C"
	}
	return Run(ctx, timeout, shell, flag, command)
}

// Which searches for an executable in PATH and returns its absolute path.
//
// Parameters:
//   - name: The executable name.
//
// Returns:
//   - The absolute path, or empty string if not found.
func Which(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

// FFmpegPath returns the path to the FFmpeg executable, checking common
// locations if it's not in PATH.
//
// Returns:
//   - The FFmpeg path, or empty string if not found.
func FFmpegPath() string {
	// Check PATH first.
	if p := Which("ffmpeg"); p != "" {
		return p
	}

	// Check common locations on different platforms.
	commonPaths := []string{
		"/usr/bin/ffmpeg",
		"/usr/local/bin/ffmpeg",
		"/opt/homebrew/bin/ffmpeg",
	}
	if IsWindows() {
		commonPaths = []string{
			`C:\ffmpeg\bin\ffmpeg.exe`,
			`C:\Program Files\ffmpeg\bin\ffmpeg.exe`,
		}
	}

	for _, p := range commonPaths {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return ""
}

// FFmpegVersion returns the FFmpeg version string, or an error if FFmpeg
// is not available.
//
// Returns:
//   - The version string (e.g. "6.0"), and any error.
func FFmpegVersion() (string, error) {
	path := FFmpegPath()
	if path == "" {
		return "", fmt.Errorf("ffmpeg not found in PATH")
	}

	result := Run(context.Background(), 10*time.Second, path, "-version")
	if result.Err != nil {
		return "", fmt.Errorf("failed to get ffmpeg version: %w", result.Err)
	}

	// First line is typically: ffmpeg version 6.0 Copyright...
	lines := strings.SplitN(result.Stdout, "\n", 2)
	if len(lines) == 0 {
		return "", fmt.Errorf("empty ffmpeg output")
	}

	parts := strings.Fields(lines[0])
	if len(parts) >= 3 {
		return parts[2], nil
	}
	return lines[0], nil
}
