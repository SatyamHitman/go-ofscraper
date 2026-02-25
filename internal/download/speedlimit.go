// =============================================================================
// FILE: internal/download/speedlimit.go
// PURPOSE: Bandwidth throttle for downloads. Implements an io.Reader wrapper
//          that limits read throughput to a configured bytes-per-second rate.
//          Ports Python utils/leaky.py.
// =============================================================================

package download

import (
	"io"
	"time"
)

// ---------------------------------------------------------------------------
// Speed limit reader
// ---------------------------------------------------------------------------

// SpeedLimitReader wraps an io.Reader to limit read throughput.
type SpeedLimitReader struct {
	reader    io.Reader
	bytesPS   int64     // bytes per second limit
	bytesRead int64     // bytes read in current window
	windowStart time.Time // current tracking window start
}

// NewSpeedLimitReader creates a new speed-limited reader.
//
// Parameters:
//   - reader: The underlying reader.
//   - bytesPerSecond: Maximum bytes per second (0 = unlimited).
//
// Returns:
//   - A speed-limited reader.
func NewSpeedLimitReader(reader io.Reader, bytesPerSecond int64) *SpeedLimitReader {
	return &SpeedLimitReader{
		reader:      reader,
		bytesPS:     bytesPerSecond,
		windowStart: time.Now(),
	}
}

// Read implements io.Reader with speed limiting.
func (s *SpeedLimitReader) Read(p []byte) (int, error) {
	if s.bytesPS <= 0 {
		return s.reader.Read(p)
	}

	// Check if we've exceeded the rate for the current second.
	elapsed := time.Since(s.windowStart)
	if elapsed >= time.Second {
		// Reset window.
		s.bytesRead = 0
		s.windowStart = time.Now()
	} else if s.bytesRead >= s.bytesPS {
		// We've hit the limit — sleep until the next window.
		sleepDuration := time.Second - elapsed
		time.Sleep(sleepDuration)
		s.bytesRead = 0
		s.windowStart = time.Now()
	}

	// Limit read size to remaining budget.
	remaining := s.bytesPS - s.bytesRead
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err := s.reader.Read(p)
	s.bytesRead += int64(n)
	return n, err
}
