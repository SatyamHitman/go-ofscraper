// =============================================================================
// FILE: internal/http/request.go
// PURPOSE: Request builder for HTTP requests. Provides a fluent API for
//          constructing requests with headers, query params, and body content.
// =============================================================================

package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ---------------------------------------------------------------------------
// Request
// ---------------------------------------------------------------------------

// Request represents an HTTP request to be executed by SessionManager.
type Request struct {
	Method  string            // HTTP method (GET, POST, etc.)
	URL     string            // Full URL
	Headers map[string]string // Request headers
	Params  url.Values        // Query parameters
	Body    io.Reader         // Request body (for POST/PUT)
}

// NewRequest creates a GET request for the given URL.
//
// Parameters:
//   - url: The request URL.
//
// Returns:
//   - A new Request configured for GET.
func NewRequest(requestURL string) *Request {
	return &Request{
		Method:  http.MethodGet,
		URL:     requestURL,
		Headers: make(map[string]string),
		Params:  make(url.Values),
	}
}

// NewPostRequest creates a POST request with a JSON body.
//
// Parameters:
//   - url: The request URL.
//   - body: The request body bytes.
//
// Returns:
//   - A new Request configured for POST.
func NewPostRequest(requestURL string, body []byte) *Request {
	return &Request{
		Method:  http.MethodPost,
		URL:     requestURL,
		Headers: map[string]string{"Content-Type": "application/json"},
		Params:  make(url.Values),
		Body:    bytes.NewReader(body),
	}
}

// WithHeader adds a header to the request.
//
// Parameters:
//   - key: Header name.
//   - value: Header value.
//
// Returns:
//   - The Request (for chaining).
func (r *Request) WithHeader(key, value string) *Request {
	r.Headers[key] = value
	return r
}

// WithParam adds a query parameter.
//
// Parameters:
//   - key: Parameter name.
//   - value: Parameter value.
//
// Returns:
//   - The Request (for chaining).
func (r *Request) WithParam(key, value string) *Request {
	r.Params.Set(key, value)
	return r
}

// Build constructs the stdlib http.Request from the builder fields.
//
// Parameters:
//   - ctx: Context for the request.
//
// Returns:
//   - A configured *http.Request, and any error.
func (r *Request) Build(ctx context.Context) (*http.Request, error) {
	// Append query params to URL.
	reqURL := r.URL
	if len(r.Params) > 0 {
		u, err := url.Parse(r.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL %q: %w", r.URL, err)
		}
		q := u.Query()
		for k, vals := range r.Params {
			for _, v := range vals {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	var body io.Reader
	if r.Body != nil {
		body = r.Body
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return req, nil
}
