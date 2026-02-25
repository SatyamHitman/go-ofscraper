// =============================================================================
// FILE: internal/auth/providers/generic.go
// PURPOSE: Generic URL-based signature provider. Fetches signing rules from
//          any URL that returns the standard JSON format. Used as a fallback
//          and as a base for URL-specific providers.
// =============================================================================

package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gofscraper/internal/auth"
)

// ---------------------------------------------------------------------------
// Generic provider
// ---------------------------------------------------------------------------

// GenericProvider fetches signing rules from a configurable URL.
type GenericProvider struct {
	name     string
	url      string
	priority int
}

// NewGenericProvider creates a generic URL-based provider.
//
// Parameters:
//   - name: Provider identifier.
//   - url: URL to fetch rules from.
//   - priority: Priority (lower = preferred).
//
// Returns:
//   - A new GenericProvider.
func NewGenericProvider(name, url string, priority int) *GenericProvider {
	return &GenericProvider{name: name, url: url, priority: priority}
}

// Name returns the provider name.
func (g *GenericProvider) Name() string { return g.name }

// Priority returns the provider priority.
func (g *GenericProvider) Priority() int { return g.priority }

// FetchRules downloads and parses signing rules from the configured URL.
func (g *GenericProvider) FetchRules(ctx context.Context) (*auth.SignatureParams, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.url, nil)
	if err != nil {
		return nil, fmt.Errorf("provider %s: failed to create request: %w", g.name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider %s: request failed: %w", g.name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider %s: HTTP %d", g.name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("provider %s: failed to read body: %w", g.name, err)
	}

	return parseRulesJSON(body, g.name)
}

// ---------------------------------------------------------------------------
// JSON parsing
// ---------------------------------------------------------------------------

// rulesResponse is the expected JSON format from rule providers.
type rulesResponse struct {
	StaticParam      string `json:"static_param"`
	Format           string `json:"format"`
	ChecksumConstant int    `json:"checksum_constant"`
	ChecksumIndexes  []int  `json:"checksum_indexes"`
	AppToken         string `json:"app_token"`
	Prefix           string `json:"prefix"`
	Suffix           string `json:"suffix"`
}

// parseRulesJSON parses the standard rules JSON format.
func parseRulesJSON(data []byte, providerName string) (*auth.SignatureParams, error) {
	var rules rulesResponse
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("provider %s: invalid JSON: %w", providerName, err)
	}

	if rules.StaticParam == "" {
		return nil, fmt.Errorf("provider %s: missing static_param", providerName)
	}

	return &auth.SignatureParams{
		StaticParam:      rules.StaticParam,
		Format:           rules.Format,
		ChecksumConstant: rules.ChecksumConstant,
		ChecksumIndexes:  rules.ChecksumIndexes,
		AppToken:         rules.AppToken,
		Prefix:           rules.Prefix,
		Suffix:           rules.Suffix,
	}, nil
}
