// =============================================================================
// FILE: internal/auth/error.go
// PURPOSE: Auth error types and handling. Defines structured errors for
//          authentication failures with guidance for resolution. Ports Python
//          utils/auth/utils/error.py.
// =============================================================================

package auth

import (
	"fmt"
)

// ---------------------------------------------------------------------------
// Auth errors
// ---------------------------------------------------------------------------

// ErrCode identifies the type of authentication failure.
type ErrCode int

const (
	ErrCodeMissing    ErrCode = iota // Auth file or credentials missing
	ErrCodeInvalid                   // Credentials present but invalid
	ErrCodeExpired                   // Session cookie expired
	ErrCodeForbidden                 // Account restricted/banned
	ErrCodeRateLimit                 // Too many auth attempts
)

// AuthError represents an authentication failure with context.
type AuthError struct {
	Code    ErrCode
	Message string
	Hint    string // User-facing guidance for fixing the issue
}

func (e *AuthError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("auth error [%d]: %s (hint: %s)", e.Code, e.Message, e.Hint)
	}
	return fmt.Sprintf("auth error [%d]: %s", e.Code, e.Message)
}

// NewMissingError creates an error for missing auth credentials.
func NewMissingError(field string) *AuthError {
	return &AuthError{
		Code:    ErrCodeMissing,
		Message: fmt.Sprintf("missing required auth field: %s", field),
		Hint:    "Run the auth setup or manually edit the auth.json file",
	}
}

// NewInvalidError creates an error for invalid auth credentials.
func NewInvalidError(detail string) *AuthError {
	return &AuthError{
		Code:    ErrCodeInvalid,
		Message: detail,
		Hint:    "Check that your cookies and tokens are current",
	}
}

// NewExpiredError creates an error for expired session.
func NewExpiredError() *AuthError {
	return &AuthError{
		Code:    ErrCodeExpired,
		Message: "session cookie has expired",
		Hint:    "Log in to OnlyFans in your browser and re-export cookies",
	}
}

// NewForbiddenError creates an error for account restrictions.
func NewForbiddenError() *AuthError {
	return &AuthError{
		Code:    ErrCodeForbidden,
		Message: "account access is forbidden",
		Hint:    "Your account may be restricted or banned",
	}
}

// IsAuthError checks if an error is an AuthError.
func IsAuthError(err error) bool {
	_, ok := err.(*AuthError)
	return ok
}
