// =============================================================================
// FILE: internal/utils/encoding.go
// PURPOSE: Encoding utility functions. Provides helpers for JSON, base64, and
//          URL encoding/decoding operations used across the codebase.
//          Ports Python utils/encoding.py.
// =============================================================================

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

// ---------------------------------------------------------------------------
// JSON helpers
// ---------------------------------------------------------------------------

// MustMarshalJSON serialises the value to JSON, panicking on failure. Use
// only for values known to be serialisable (e.g. config structs).
//
// Parameters:
//   - v: The value to serialise.
//
// Returns:
//   - JSON bytes.
func MustMarshalJSON(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("json marshal failed: %v", err))
	}
	return data
}

// PrettyJSON returns an indented JSON string for the given value.
//
// Parameters:
//   - v: The value to serialise.
//
// Returns:
//   - Indented JSON string, or the error message if serialisation fails.
func PrettyJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("<json error: %v>", err)
	}
	return string(data)
}

// ParseJSON deserialises JSON bytes into the target value.
//
// Parameters:
//   - data: JSON bytes.
//   - target: Pointer to the destination value.
//
// Returns:
//   - Error if parsing fails.
func ParseJSON(data []byte, target any) error {
	return json.Unmarshal(data, target)
}

// ---------------------------------------------------------------------------
// Base64 helpers
// ---------------------------------------------------------------------------

// Base64Encode encodes bytes to a standard base64 string.
//
// Parameters:
//   - data: The raw bytes to encode.
//
// Returns:
//   - Base64-encoded string.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode decodes a standard base64 string to bytes.
//
// Parameters:
//   - s: The base64 string.
//
// Returns:
//   - Decoded bytes, and any error.
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// ---------------------------------------------------------------------------
// URL helpers
// ---------------------------------------------------------------------------

// URLEncode encodes a map of key-value pairs as a URL query string.
//
// Parameters:
//   - params: The key-value pairs to encode.
//
// Returns:
//   - URL-encoded query string.
func URLEncode(params map[string]string) string {
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	return v.Encode()
}

// URLJoin concatenates a base URL with path segments safely.
//
// Parameters:
//   - base: The base URL (e.g. "https://api.example.com").
//   - paths: Path segments to append.
//
// Returns:
//   - The joined URL string, and any error.
func URLJoin(base string, paths ...string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("invalid base URL %q: %w", base, err)
	}
	for _, p := range paths {
		ref, err := url.Parse(p)
		if err != nil {
			return "", fmt.Errorf("invalid path segment %q: %w", p, err)
		}
		u = u.ResolveReference(ref)
	}
	return u.String(), nil
}
