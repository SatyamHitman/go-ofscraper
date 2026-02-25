// =============================================================================
// FILE: internal/auth/signing.go
// PURPOSE: SHA1-based request signing algorithm. Creates the "sign" and "time"
//          headers required by the OF API using dynamic signature rules fetched
//          from external providers. Exact port of Python
//          utils/auth/request.py create_login_sign.
// =============================================================================

package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Signature parameters (from dynamic rules)
// ---------------------------------------------------------------------------

// SignatureParams holds the dynamic signing parameters fetched from a provider.
type SignatureParams struct {
	StaticParam      string // Static portion of the sign message
	Format           string // Format string for the final sign value
	ChecksumConstant int    // Constant added to checksum
	ChecksumIndexes  []int  // Character indexes in SHA1 hex for checksum
	AppToken         string // App token override (if provided by rules)
	Prefix           string // Sign prefix
	Suffix           string // Sign suffix
}

// ---------------------------------------------------------------------------
// Sign creation
// ---------------------------------------------------------------------------

// CreateSign generates the "sign" and "time" headers for an API request.
// This is the exact algorithm used by the OF web client.
//
// Parameters:
//   - link: The full request URL (path + query used for signing).
//   - headers: The request headers map — "sign" and "time" are set in place.
//   - params: Dynamic signature parameters from a provider.
func CreateSign(link string, headers map[string]string, params SignatureParams) {
	// Current time in milliseconds.
	time2 := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Extract path + query from the URL.
	u, err := url.Parse(link)
	if err != nil {
		return
	}
	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	// Build the message: static_param\ntime\npath\nuser_id
	msg := strings.Join([]string{
		params.StaticParam,
		time2,
		path,
		headers["user-id"],
	}, "\n")

	// SHA1 hash of the message.
	h := sha1.New()
	h.Write([]byte(msg))
	sha1Hex := hex.EncodeToString(h.Sum(nil))

	// Calculate checksum from specific character positions.
	checksum := 0
	for _, idx := range params.ChecksumIndexes {
		if idx < len(sha1Hex) {
			checksum += int(sha1Hex[idx])
		}
	}
	checksum += params.ChecksumConstant

	// Build the sign value using the format template.
	sign := fmt.Sprintf(params.Format, sha1Hex, int(math.Abs(float64(checksum))))

	headers["sign"] = sign
	headers["time"] = time2
}
