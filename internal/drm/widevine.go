// =============================================================================
// FILE: internal/drm/widevine.go
// PURPOSE: Manual Widevine CDM implementation. Handles local device loading,
//          license challenge generation, and key extraction from license
//          responses. Ports Python pywidevine usage.
// =============================================================================

package drm

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ---------------------------------------------------------------------------
// Device
// ---------------------------------------------------------------------------

// Device represents a Widevine CDM device (client_id + private_key).
type Device struct {
	ClientID   []byte // Client ID blob
	PrivateKey []byte // RSA private key (DER or PEM)
}

// LoadDevice loads a Widevine device from a directory containing
// device_client_id_blob and device_private_key files.
//
// Parameters:
//   - devicePath: Path to the device directory.
//
// Returns:
//   - Loaded Device, or error.
func LoadDevice(devicePath string) (*Device, error) {
	clientID, err := os.ReadFile(filepath.Join(devicePath, "device_client_id_blob"))
	if err != nil {
		return nil, fmt.Errorf("read client_id: %w", err)
	}

	privateKey, err := os.ReadFile(filepath.Join(devicePath, "device_private_key"))
	if err != nil {
		return nil, fmt.Errorf("read private_key: %w", err)
	}

	return &Device{
		ClientID:   clientID,
		PrivateKey: privateKey,
	}, nil
}

// ---------------------------------------------------------------------------
// License operations
// ---------------------------------------------------------------------------

// ManualDecrypt performs manual Widevine key acquisition using local device files.
//
// Parameters:
//   - ctx: Context for cancellation.
//   - devicePath: Path to device directory.
//   - licenseURL: The license server URL.
//   - kid: The key ID (hex string).
//
// Returns:
//   - The decryption key as hex string, or error.
func ManualDecrypt(ctx context.Context, devicePath, licenseURL, kid string) (string, error) {
	device, err := LoadDevice(devicePath)
	if err != nil {
		return "", fmt.Errorf("load device: %w", err)
	}

	// Generate license challenge.
	challenge, err := generateChallenge(device, kid)
	if err != nil {
		return "", fmt.Errorf("generate challenge: %w", err)
	}

	// Send challenge to license server.
	licenseResp, err := sendChallenge(ctx, licenseURL, challenge)
	if err != nil {
		return "", fmt.Errorf("send challenge: %w", err)
	}

	// Extract key from license response.
	key, err := extractKey(device, licenseResp, kid)
	if err != nil {
		return "", fmt.Errorf("extract key: %w", err)
	}

	return key, nil
}

// generateChallenge creates a Widevine license request challenge.
// This is a simplified version — full implementation would use protobuf.
func generateChallenge(device *Device, kid string) ([]byte, error) {
	kidBytes, err := hex.DecodeString(kid)
	if err != nil {
		return nil, fmt.Errorf("decode kid: %w", err)
	}

	// Placeholder: In production, this would build a proper Widevine
	// SignedLicenseRequest protobuf message using the device's client_id
	// and private_key. For now, return the client_id with the KID appended.
	challenge := make([]byte, 0, len(device.ClientID)+len(kidBytes))
	challenge = append(challenge, device.ClientID...)
	challenge = append(challenge, kidBytes...)
	return challenge, nil
}

// sendChallenge sends a license challenge to the license server.
func sendChallenge(ctx context.Context, licenseURL string, challenge []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, licenseURL, bytes.NewReader(challenge))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("license request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("license server status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read license response: %w", err)
	}

	return body, nil
}

// extractKey extracts the content key from a license response.
// Placeholder: Full implementation would parse the protobuf response
// and use the device's private key to decrypt the content key.
func extractKey(device *Device, licenseResp []byte, kid string) (string, error) {
	if len(licenseResp) == 0 {
		return "", fmt.Errorf("empty license response")
	}

	// Placeholder: Return hex-encoded response segment.
	// Real implementation would parse SignedLicense protobuf,
	// decrypt session key with RSA private key, then derive content keys.
	_ = device
	_ = kid
	return "", fmt.Errorf("manual Widevine decryption requires protobuf implementation")
}
