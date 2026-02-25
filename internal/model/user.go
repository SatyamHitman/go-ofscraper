// =============================================================================
// FILE: internal/model/user.go
// PURPOSE: Defines the User (Model/Creator) domain model representing an
//          OnlyFans content creator. Ports Python classes/of/models.py with
//          all subscription pricing, promo, and status fields.
// =============================================================================

package model

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// ---------------------------------------------------------------------------
// Promo represents a promotional pricing offer.
// ---------------------------------------------------------------------------

// Promo holds data for a single promotional pricing entry from the API.
type Promo struct {
	Price    float64 `json:"price"`
	CanClaim bool    `json:"canClaim"`
}

// ---------------------------------------------------------------------------
// SubscribedData holds subscription details when actively subscribed.
// ---------------------------------------------------------------------------

// SubscribedData contains details about the current subscription relationship.
type SubscribedData struct {
	RegularPrice float64 `json:"regularPrice"`
	StartDate    string  `json:"startDate,omitempty"`
	ExpireDate   string  `json:"expireDate,omitempty"`
	RenewedAt    string  `json:"renewedAt,omitempty"`
	Status       string  `json:"status,omitempty"` // e.g. "Set to Expire"
}

// ---------------------------------------------------------------------------
// User struct
// ---------------------------------------------------------------------------

// User represents an OnlyFans content creator (also called "Model" in the API).
// Contains identity, subscription pricing, promo offers, and activity data.
type User struct {
	// --- Identification ---
	ID       int64  `json:"id"`
	Name     string `json:"name"`     // Username
	Avatar   string `json:"avatar"`   // Avatar image URL
	Header   string `json:"header"`   // Header/banner image URL

	// --- Activity ---
	LastSeen string `json:"last_seen,omitempty"` // Last activity timestamp

	// --- Subscription pricing ---
	CurrentSubscribePrice float64         `json:"currentSubscribePrice"`
	SubscribePrice        float64         `json:"subscribePrice"`     // Base/regular price
	SubscribedData        *SubscribedData `json:"subscribedByData,omitempty"`
	SubscribedExpiredDate string          `json:"subscribedByExpireDate,omitempty"`

	// --- Subscription status ---
	SubscribedAt string `json:"subscribedAt,omitempty"` // Subscribe start date
	RenewedAt    string `json:"renewedAt,omitempty"`
	ExpiredAt    string `json:"expiredAt,omitempty"`

	// --- Promos ---
	Promos []Promo `json:"promos,omitempty"`

	// --- Flags ---
	IsRealPerformer bool `json:"isRealPerformer"`
	IsRestricted    bool `json:"isRestricted"`
}

// ---------------------------------------------------------------------------
// Promo computation methods
// ---------------------------------------------------------------------------

// AllClaimablePromos returns promos where CanClaim is true, sorted by price
// ascending (cheapest first).
//
// Returns:
//   - Sorted slice of claimable Promo entries.
func (u *User) AllClaimablePromos() []Promo {
	var claimable []Promo
	for _, p := range u.Promos {
		if p.CanClaim {
			claimable = append(claimable, p)
		}
	}
	sort.Slice(claimable, func(i, j int) bool {
		return claimable[i].Price < claimable[j].Price
	})
	return claimable
}

// LowestClaimablePromo returns the lowest claimable promo price, or -1 if none.
//
// Returns:
//   - The lowest claimable promo price, or -1 if no claimable promos exist.
func (u *User) LowestClaimablePromo() float64 {
	promos := u.AllClaimablePromos()
	if len(promos) == 0 {
		return -1
	}
	return promos[0].Price
}

// AllPromosSorted returns all promos sorted by price ascending.
//
// Returns:
//   - Sorted slice of all Promo entries.
func (u *User) AllPromosSorted() []Promo {
	sorted := make([]Promo, len(u.Promos))
	copy(sorted, u.Promos)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Price < sorted[j].Price
	})
	return sorted
}

// LowestPromoAll returns the lowest promo price across all promos, or -1.
//
// Returns:
//   - The lowest promo price, or -1 if no promos exist.
func (u *User) LowestPromoAll() float64 {
	if len(u.Promos) == 0 {
		return -1
	}
	all := u.AllPromosSorted()
	return all[0].Price
}

// ---------------------------------------------------------------------------
// Price calculation methods
// ---------------------------------------------------------------------------

// RegularPrice returns the regular subscription price. Checks SubscribedData
// first, then falls back to SubscribePrice.
//
// Returns:
//   - The regular price, or 0 if not available.
func (u *User) RegularPrice() float64 {
	if u.SubscribedData != nil && u.SubscribedData.RegularPrice > 0 {
		return u.SubscribedData.RegularPrice
	}
	return u.SubscribePrice
}

// FinalCurrentPrice returns the effective current subscription price.
// Priority: CurrentSubscribePrice > LowestClaimablePromo > RegularPrice > 0.
//
// Returns:
//   - The effective current price.
func (u *User) FinalCurrentPrice() float64 {
	if u.CurrentSubscribePrice > 0 {
		return u.CurrentSubscribePrice
	}
	lcp := u.LowestClaimablePromo()
	if lcp >= 0 {
		return lcp
	}
	rp := u.RegularPrice()
	if rp > 0 {
		return rp
	}
	return 0
}

// FinalRenewalPrice returns the renewal price.
// Priority: LowestClaimablePromo > RegularPrice > 0.
//
// Returns:
//   - The effective renewal price.
func (u *User) FinalRenewalPrice() float64 {
	lcp := u.LowestClaimablePromo()
	if lcp >= 0 {
		return lcp
	}
	rp := u.RegularPrice()
	if rp > 0 {
		return rp
	}
	return 0
}

// FinalPromoPrice returns the lowest promo price across all promos.
// Fallback: RegularPrice > 0.
//
// Returns:
//   - The lowest promo price.
func (u *User) FinalPromoPrice() float64 {
	lpa := u.LowestPromoAll()
	if lpa >= 0 {
		return lpa
	}
	rp := u.RegularPrice()
	if rp > 0 {
		return rp
	}
	return 0
}

// ---------------------------------------------------------------------------
// Subscription status methods
// ---------------------------------------------------------------------------

// IsActive returns true if the user has an active subscription.
// Active means: status is "Set to Expire", has renewedAt, or not yet expired.
//
// Returns:
//   - true if subscription is currently active.
func (u *User) IsActive() bool {
	// Check subscribed data status
	if u.SubscribedData != nil {
		if u.SubscribedData.Status == "Set to Expire" {
			return true
		}
		if u.SubscribedData.RenewedAt != "" {
			return true
		}
	}

	// Check if renewed
	if u.RenewedAt != "" {
		return true
	}

	// Check expiry — active if not yet expired
	expiry := u.ExpiryString()
	if expiry == "" {
		return false
	}

	t, err := parseFlexibleDate(expiry)
	if err != nil {
		return false
	}

	return t.After(time.Now())
}

// ExpiryString returns the expiry date string from the best available source.
// Priority: ExpiredAt > SubscribedExpiredDate > SubscribedData.ExpireDate.
//
// Returns:
//   - The expiry date string, or empty if not available.
func (u *User) ExpiryString() string {
	if u.ExpiredAt != "" {
		return u.ExpiredAt
	}
	if u.SubscribedExpiredDate != "" {
		return u.SubscribedExpiredDate
	}
	if u.SubscribedData != nil && u.SubscribedData.ExpireDate != "" {
		return u.SubscribedData.ExpireDate
	}
	return ""
}

// FinalExpired returns the expiry as a Unix timestamp float, or 0 if not set.
//
// Returns:
//   - Unix timestamp as float64, or 0.
func (u *User) FinalExpired() float64 {
	expiry := u.ExpiryString()
	if expiry == "" {
		return 0
	}

	t, err := parseFlexibleDate(expiry)
	if err != nil {
		return 0
	}

	return float64(t.Unix())
}

// ---------------------------------------------------------------------------
// Date formatting methods
// ---------------------------------------------------------------------------

// LastSeenFormatted returns the last seen date formatted as "2006-01-02".
//
// Returns:
//   - Formatted date string, or empty if not available.
func (u *User) LastSeenFormatted() string {
	return formatDateShort(u.LastSeen)
}

// FinalLastSeen returns the last seen timestamp as float64 (now if missing).
//
// Returns:
//   - Unix timestamp as float64.
func (u *User) FinalLastSeen() float64 {
	if u.LastSeen == "" {
		return float64(time.Now().Unix())
	}
	t, err := parseFlexibleDate(u.LastSeen)
	if err != nil {
		return float64(time.Now().Unix())
	}
	return float64(t.Unix())
}

// SubscribedFormatted returns the subscription start date as "2006-01-02".
//
// Returns:
//   - Formatted date string, or empty.
func (u *User) SubscribedFormatted() string {
	return formatDateShort(u.SubscribedAt)
}

// FinalSubscribed returns the subscribe date as float64 Unix timestamp, or 0.
//
// Returns:
//   - Unix timestamp as float64.
func (u *User) FinalSubscribed() float64 {
	if u.SubscribedAt == "" {
		return 0
	}
	t, err := parseFlexibleDate(u.SubscribedAt)
	if err != nil {
		return 0
	}
	return float64(t.Unix())
}

// RenewedFormatted returns the renewal date as "2006-01-02".
//
// Returns:
//   - Formatted date string, or empty.
func (u *User) RenewedFormatted() string {
	return formatDateShort(u.RenewedAt)
}

// ExpiredFormatted returns the expiry date as "2006-01-02".
//
// Returns:
//   - Formatted date string, or empty.
func (u *User) ExpiredFormatted() string {
	return formatDateShort(u.ExpiryString())
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseFlexibleDate attempts to parse a date string using common API formats.
//
// Parameters:
//   - s: The date string to parse.
//
// Returns:
//   - The parsed time.Time and nil error, or zero time and error on failure.
func parseFlexibleDate(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse date: %s", s)
}

// formatDateShort parses a date string and returns "2006-01-02" format.
//
// Parameters:
//   - s: The date string to format.
//
// Returns:
//   - Formatted short date, or empty if parsing fails.
func formatDateShort(s string) string {
	if s == "" {
		return ""
	}
	t, err := parseFlexibleDate(s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// abs returns the absolute value of an integer.
// Used in price calculations and checksum computation.
//
// Parameters:
//   - x: The input integer.
//
// Returns:
//   - The absolute value.
func abs(x int) int {
	return int(math.Abs(float64(x)))
}

