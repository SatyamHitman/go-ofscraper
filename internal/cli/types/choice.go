// =============================================================================
// FILE: internal/cli/types/choice.go
// PURPOSE: ChoiceValue implements pflag.Value for restricted choice inputs.
//          Only values from a predefined set are accepted.
// =============================================================================

package types

import (
	"fmt"
	"strings"
)

// ChoiceValue is a custom pflag.Value that restricts input to a set of allowed
// values.
type ChoiceValue struct {
	Allowed []string
	Value   string
}

// NewChoiceValue returns a ChoiceValue with the given allowed values and a
// default selection.
func NewChoiceValue(allowed []string, defaultVal string) *ChoiceValue {
	return &ChoiceValue{
		Allowed: allowed,
		Value:   defaultVal,
	}
}

// String returns the current value.
func (c *ChoiceValue) String() string {
	return c.Value
}

// Set validates and stores the given string.
func (c *ChoiceValue) Set(s string) error {
	lower := strings.ToLower(strings.TrimSpace(s))
	for _, a := range c.Allowed {
		if strings.ToLower(a) == lower {
			c.Value = a
			return nil
		}
	}
	return fmt.Errorf("invalid value %q: must be one of [%s]", s, strings.Join(c.Allowed, ", "))
}

// Type returns the type name for help output.
func (c *ChoiceValue) Type() string {
	return "choice"
}

// Get returns the current value as a string.
func (c *ChoiceValue) Get() string {
	return c.Value
}
