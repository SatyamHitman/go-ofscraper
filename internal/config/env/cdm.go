// =============================================================================
// FILE: internal/config/env/cdm.go
// PURPOSE: CDM (Content Decryption Module) environment variable defaults.
//          Ports Python of_env/values/req/cdm.py.
// =============================================================================

package env

// CDMPrivateKey returns the path to the Widevine private key file.
func CDMPrivateKey() string {
	return GetString("OF_CDM_PRIVATE_KEY", "")
}

// CDMClientID returns the path to the Widevine client ID file.
func CDMClientID() string {
	return GetString("OF_CDM_CLIENT_ID", "")
}
