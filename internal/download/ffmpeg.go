// =============================================================================
// FILE: internal/download/ffmpeg.go
// PURPOSE: FFmpeg path detection and operations. Locates the FFmpeg binary
//          and provides wrappers for common FFmpeg operations.
//          Ports Python utils/ffmpeg.py.
// =============================================================================

package download

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ---------------------------------------------------------------------------
// FFmpeg detection
// ---------------------------------------------------------------------------

// FindFFmpeg locates the FFmpeg binary. Checks the provided path first,
// then falls back to PATH lookup.
//
// Parameters:
//   - configPath: Configured FFmpeg path (may be empty).
//
// Returns:
//   - The absolute path to ffmpeg, or error if not found.
func FindFFmpeg(configPath string) (string, error) {
	if configPath != "" {
		if _, err := exec.LookPath(configPath); err == nil {
			return configPath, nil
		}
	}

	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}
	return path, nil
}

// FFmpegVersion returns the FFmpeg version string.
//
// Parameters:
//   - ffmpegPath: Path to the FFmpeg binary.
//
// Returns:
//   - Version string, or error.
func FFmpegVersion(ffmpegPath string) (string, error) {
	cmd := exec.Command(ffmpegPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg -version: %w", err)
	}
	// First line contains version.
	lines := strings.SplitN(string(output), "\n", 2)
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

// ConcatFiles concatenates multiple media files using FFmpeg.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - ffmpegPath: Path to FFmpeg.
//   - inputs: Input file paths.
//   - output: Output file path.
//
// Returns:
//   - Error if the operation fails.
func ConcatFiles(ctx context.Context, ffmpegPath string, inputs []string, output string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("no input files")
	}

	// Build concat filter input.
	var args []string
	args = append(args, "-y")
	for _, input := range inputs {
		args = append(args, "-i", input)
	}

	if len(inputs) > 1 {
		filter := fmt.Sprintf("concat=n=%d:v=1:a=1", len(inputs))
		args = append(args, "-filter_complex", filter)
	}

	args = append(args, "-c", "copy", output)

	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat: %w\noutput: %s", err, string(out))
	}
	return nil
}
