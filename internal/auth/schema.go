// =============================================================================
// FILE: internal/auth/schema.go
// PURPOSE: Auth schema validation. Defines the expected structure and field
//          requirements for auth data, including field types and constraints.
//          Ports Python utils/auth/schema.py.
// =============================================================================

package auth

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Schema validation
// ---------------------------------------------------------------------------

// FieldRequirement describes a required auth field.
type FieldRequirement struct {
	Name     string // JSON key name
	Label    string // Human-readable label
	Required bool   // Whether the field is mandatory
}

// Schema returns the list of auth field requirements.
//
// Returns:
//   - Slice of FieldRequirement describing all auth fields.
func Schema() []FieldRequirement {
	return []FieldRequirement{
		{Name: "auth_cookie", Label: "Session Cookie (sess)", Required: true},
		{Name: "auth_user_agent", Label: "User Agent", Required: true},
		{Name: "auth_x_bc", Label: "X-BC Token", Required: true},
		{Name: "auth_user_id", Label: "User ID", Required: true},
		{Name: "auth_app_token", Label: "App Token", Required: false},
	}
}

// ValidateSchema checks auth data against the schema and returns all
// validation errors (not just the first one).
//
// Parameters:
//   - data: The auth data to validate.
//
// Returns:
//   - Slice of error messages, empty if valid.
func ValidateSchema(data *Data) []string {
	var errors []string

	if data == nil {
		return []string{"auth data is nil"}
	}

	fields := map[string]string{
		"auth_cookie":     data.Cookie,
		"auth_user_agent": data.UserAgent,
		"auth_x_bc":       data.XBC,
		"auth_user_id":    data.UserID,
	}

	for _, req := range Schema() {
		if !req.Required {
			continue
		}
		val, ok := fields[req.Name]
		if !ok || strings.TrimSpace(val) == "" {
			errors = append(errors, fmt.Sprintf("%s (%s) is required", req.Label, req.Name))
		}
	}

	return errors
}
