// =============================================================================
// FILE: internal/auth/providers/provider.go
// PURPOSE: SignatureProvider interface. Defines the contract for dynamic
//          signing rule sources. Each provider fetches signature parameters
//          (static_param, checksum indexes, etc.) from a different source.
//          Ports Python auth/request.py dynamic rule fetching.
// =============================================================================

package providers

import (
	"context"

	"gofscraper/internal/auth"
)

// ---------------------------------------------------------------------------
// Provider interface
// ---------------------------------------------------------------------------

// SignatureProvider fetches dynamic signing rules from an external source.
type SignatureProvider interface {
	// Name returns the human-readable provider identifier.
	Name() string

	// FetchRules retrieves signing parameters from this provider's source.
	// Returns error if the source is unavailable or returns invalid data.
	FetchRules(ctx context.Context) (*auth.SignatureParams, error)

	// Priority returns the provider's priority (lower = preferred).
	Priority() int
}

// ---------------------------------------------------------------------------
// Provider registry
// ---------------------------------------------------------------------------

// registry holds all registered providers sorted by priority.
var registry []SignatureProvider

// Register adds a provider to the registry.
//
// Parameters:
//   - p: The provider to register.
func Register(p SignatureProvider) {
	registry = append(registry, p)
}

// All returns all registered providers.
//
// Returns:
//   - Slice of registered SignatureProviders.
func All() []SignatureProvider {
	return registry
}

// ByName returns a provider by its name.
//
// Parameters:
//   - name: The provider name.
//
// Returns:
//   - The provider, or nil if not found.
func ByName(name string) SignatureProvider {
	for _, p := range registry {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// FetchFirst tries each provider in order and returns the first successful result.
//
// Parameters:
//   - ctx: Context for cancellation.
//
// Returns:
//   - The signing parameters from the first successful provider, and any error.
func FetchFirst(ctx context.Context) (*auth.SignatureParams, error) {
	var lastErr error
	for _, p := range registry {
		params, err := p.FetchRules(ctx)
		if err == nil {
			return params, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, auth.NewMissingError("no signature providers available")
}
