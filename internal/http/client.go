// =============================================================================
// FILE: internal/http/client.go
// PURPOSE: SessionManager HTTP client. Manages connection pooling, auth header
//          injection, rate limiting, and adaptive sleeping for OF API requests.
//          Ports Python managers/sessionmanager/sessionmanager.py.
// =============================================================================

package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"gofscraper/internal/auth"
	"gofscraper/internal/config/env"
)

// ---------------------------------------------------------------------------
// SessionManager
// ---------------------------------------------------------------------------

// SessionManager wraps an http.Client with auth, rate limiting, and retry.
type SessionManager struct {
	client    *http.Client
	authData  *auth.Data
	rateLimit *RateLimiter
	sleeper   *AdaptiveSleeper
	mu        sync.RWMutex

	// Request counters for stats.
	totalRequests int64
	failedRequests int64
}

// New creates a SessionManager with the given auth credentials.
//
// Parameters:
//   - authData: Authentication credentials for request signing.
//
// Returns:
//   - A configured SessionManager.
func New(authData *auth.Data) *SessionManager {
	transport := &http.Transport{
		MaxIdleConns:        env.MaxConnections(),
		MaxIdleConnsPerHost: env.MaxConnections(),
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(env.ConnectTimeout()) * time.Second,
	}

	return &SessionManager{
		client:    client,
		authData:  authData,
		rateLimit: NewRateLimiter(env.RateLimitEnabled()),
		sleeper:   NewAdaptiveSleeper(),
	}
}

// Do executes an HTTP request with auth headers, rate limiting, and retry.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - req: The request to execute.
//
// Returns:
//   - The Response wrapper, and any error.
func (sm *SessionManager) Do(ctx context.Context, req *Request) (*Response, error) {
	// Apply rate limiting.
	if err := sm.rateLimit.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	// Apply adaptive sleep if needed.
	sm.sleeper.Sleep(ctx)

	// Build the http.Request.
	httpReq, err := req.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Inject auth headers.
	sm.mu.RLock()
	if sm.authData != nil {
		auth.AddAuthHeaders(req.Headers, sm.authData)
	}
	sm.mu.RUnlock()

	// Apply headers to the http.Request.
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Set cookies.
	if sm.authData != nil {
		httpReq.Header.Set("Cookie", auth.MakeCookies(sm.authData))
	}

	sm.totalRequests++

	// Execute request.
	resp, err := sm.client.Do(httpReq)
	if err != nil {
		sm.failedRequests++
		return nil, fmt.Errorf("request failed: %w", err)
	}

	wrapped := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       resp.Body,
	}

	// Handle rate limit responses.
	if resp.StatusCode == 429 || resp.StatusCode == 504 {
		sm.sleeper.OnRateLimit()
		slog.Debug("rate limited", "status", resp.StatusCode, "url", req.URL)
	}

	// Handle forbidden responses.
	if resp.StatusCode == 403 {
		sm.sleeper.OnForbidden()
		slog.Warn("forbidden response", "status", resp.StatusCode, "url", req.URL)
	}

	return wrapped, nil
}

// DoStream executes a request and returns the response body for streaming.
//
// Parameters:
//   - ctx: Context.
//   - req: The request.
//
// Returns:
//   - A StreamResponse for reading the body in chunks, and any error.
func (sm *SessionManager) DoStream(ctx context.Context, req *Request) (*StreamResponse, error) {
	resp, err := sm.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	return &StreamResponse{
		Response: resp,
		Reader:   resp.Body,
	}, nil
}

// Close releases resources held by the session manager.
//
// Returns:
//   - Always nil.
func (sm *SessionManager) Close() error {
	sm.client.CloseIdleConnections()
	return nil
}

// UpdateAuth replaces the auth credentials used for requests.
//
// Parameters:
//   - data: The new auth credentials.
func (sm *SessionManager) UpdateAuth(data *auth.Data) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.authData = data
}

// Stats returns request statistics.
func (sm *SessionManager) Stats() (total, failed int64) {
	return sm.totalRequests, sm.failedRequests
}

// ---------------------------------------------------------------------------
// StreamResponse
// ---------------------------------------------------------------------------

// StreamResponse wraps a Response with a streaming reader.
type StreamResponse struct {
	*Response
	Reader io.ReadCloser
}
