// =============================================================================
// FILE: internal/http/response.go
// PURPOSE: Response wrapper for HTTP responses. Provides convenience methods
//          for reading JSON bodies, checking status codes, and extracting
//          headers. Used by the API client for all response handling.
// =============================================================================

package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ---------------------------------------------------------------------------
// Response
// ---------------------------------------------------------------------------

// Response wraps an HTTP response with convenience methods.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser
}

// IsOK reports whether the response has a 2xx status code.
//
// Returns:
//   - true if 200 <= status < 300.
func (r *Response) IsOK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsRateLimit reports whether the response is a rate limit (429 or 504).
func (r *Response) IsRateLimit() bool {
	return r.StatusCode == 429 || r.StatusCode == 504
}

// IsForbidden reports whether the response is a 403.
func (r *Response) IsForbidden() bool {
	return r.StatusCode == 403
}

// IsAuthError reports whether the response is a 401 or 400.
func (r *Response) IsAuthError() bool {
	return r.StatusCode == 401 || r.StatusCode == 400
}

// IsNotFound reports whether the response is a 404.
func (r *Response) IsNotFound() bool {
	return r.StatusCode == 404
}

// JSON reads and decodes the response body into the given target.
//
// Parameters:
//   - target: Pointer to the decode destination.
//
// Returns:
//   - Error if reading or decoding fails.
func (r *Response) JSON(target any) error {
	if r.Body == nil {
		return fmt.Errorf("response body is nil")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode JSON (status %d): %w", r.StatusCode, err)
	}

	return nil
}

// ReadBody reads and returns the full response body as bytes.
//
// Returns:
//   - The body bytes, and any error.
func (r *Response) ReadBody() ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// Close releases the response body.
func (r *Response) Close() {
	if r.Body != nil {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}
}
