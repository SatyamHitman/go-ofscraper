// =============================================================================
// FILE: internal/drm/cdrm.go
// PURPOSE: CDRM service client. Communicates with a remote CDRM (Content
//          Decryption Reference Module) service to obtain Widevine content
//          decryption keys. Ports Python CDRM service integration.
// =============================================================================

package drm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ---------------------------------------------------------------------------
// CDRM client
// ---------------------------------------------------------------------------

// CDRMClient communicates with a remote CDRM service for key acquisition.
type CDRMClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCDRMClient creates a new CDRM client.
//
// Parameters:
//   - baseURL: The CDRM service base URL (e.g. "https://cdrm.example.com").
//
// Returns:
//   - Configured CDRMClient.
func NewCDRMClient(baseURL string) *CDRMClient {
	return &CDRMClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// ---------------------------------------------------------------------------
// Key request/response
// ---------------------------------------------------------------------------

// cdrmRequest is the JSON request body sent to the CDRM service.
type cdrmRequest struct {
	LicenseURL string `json:"license_url"`
	PSSH       string `json:"pssh,omitempty"`
	KID        string `json:"kid,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// cdrmResponse is the JSON response from the CDRM service.
type cdrmResponse struct {
	Keys    []cdrmKey `json:"keys"`
	Message string    `json:"message,omitempty"`
}

// cdrmKey represents a single decryption key from the response.
type cdrmKey struct {
	KID  string `json:"kid"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

// ---------------------------------------------------------------------------
// Key acquisition
// ---------------------------------------------------------------------------

// GetKey obtains a content decryption key from the CDRM service.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - licenseURL: The DRM license server URL.
//   - kid: The content key ID (hex).
//
// Returns:
//   - The decryption key as hex string, or error.
func (c *CDRMClient) GetKey(ctx context.Context, licenseURL, kid string) (string, error) {
	return c.GetKeyWithHeaders(ctx, licenseURL, kid, nil)
}

// GetKeyWithHeaders obtains a key, passing additional headers to the license server.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - licenseURL: The DRM license server URL.
//   - kid: The content key ID (hex).
//   - headers: Additional headers for the license request.
//
// Returns:
//   - The decryption key as hex string, or error.
func (c *CDRMClient) GetKeyWithHeaders(ctx context.Context, licenseURL, kid string, headers map[string]string) (string, error) {
	reqBody := cdrmRequest{
		LicenseURL: licenseURL,
		KID:        kid,
		Headers:    headers,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/decrypt", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("CDRM request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CDRM status %d: %s", resp.StatusCode, string(body))
	}

	var cdrmResp cdrmResponse
	if err := json.Unmarshal(body, &cdrmResp); err != nil {
		return "", fmt.Errorf("parse CDRM response: %w", err)
	}

	// Find the CONTENT key matching the KID.
	for _, k := range cdrmResp.Keys {
		if k.Type == "CONTENT" {
			return k.Key, nil
		}
	}

	// Fallback: return first key if no CONTENT type.
	if len(cdrmResp.Keys) > 0 {
		return cdrmResp.Keys[0].Key, nil
	}

	return "", fmt.Errorf("no keys returned from CDRM service")
}
