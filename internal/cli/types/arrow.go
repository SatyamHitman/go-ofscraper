// =============================================================================
// FILE: internal/cli/types/arrow.go
// PURPOSE: DateValue implements pflag.Value for date inputs. Accepts date
//          strings in multiple formats and stores them as time.Time.
// =============================================================================

package types

import (
	"time"

	"gofscraper/internal/cli/flags"
)

// DateValue is a custom pflag.Value that parses date strings.
type DateValue struct {
	Value time.Time
	Raw   string
}

// NewDateValue returns a DateValue initialised with the zero time.
func NewDateValue() *DateValue {
	return &DateValue{}
}

// String returns the raw string representation.
func (d *DateValue) String() string {
	return d.Raw
}

// Set parses the given string into a time.Time using the shared date parser.
func (d *DateValue) Set(s string) error {
	t, err := flags.ParseDate(s)
	if err != nil {
		return err
	}
	d.Value = t
	d.Raw = s
	return nil
}

// Type returns the type name for help output.
func (d *DateValue) Type() string {
	return "date"
}

// Time returns the parsed time value.
func (d *DateValue) Time() time.Time {
	return d.Value
}

// IsZero returns true if no date has been set.
func (d *DateValue) IsZero() bool {
	return d.Value.IsZero()
}
