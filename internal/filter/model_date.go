// =============================================================================
// FILE: internal/filter/model_date.go
// PURPOSE: Model date filters. Filters users/models by subscription dates
//          (subscribed, expired, last seen). Ports Python filters/models/date.py.
// =============================================================================

package filter

import (
	"time"

	"gofscraper/internal/model"
	"gofscraper/internal/utils"
)

// ---------------------------------------------------------------------------
// Last seen date filter
// ---------------------------------------------------------------------------

// ByLastSeen returns a filter that keeps users whose last seen date is within
// the given range.
//
// Parameters:
//   - after: Only keep users last seen after this time (zero = no lower bound).
//   - before: Only keep users last seen before this time (zero = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are zero.
func ByLastSeen(after, before time.Time) ModelFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if u.LastSeen == "" {
				result = append(result, u)
				continue
			}
			t, err := utils.ParseFlexibleDate(u.LastSeen)
			if err != nil {
				result = append(result, u)
				continue
			}
			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)
			if afterOK && beforeOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Subscribed date filter
// ---------------------------------------------------------------------------

// BySubscribedDate returns a filter that keeps users whose subscription start
// date is within the given range.
//
// Parameters:
//   - after: Only keep users subscribed after this time (zero = no lower bound).
//   - before: Only keep users subscribed before this time (zero = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are zero.
func BySubscribedDate(after, before time.Time) ModelFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if u.SubscribedAt == "" {
				result = append(result, u)
				continue
			}
			t, err := utils.ParseFlexibleDate(u.SubscribedAt)
			if err != nil {
				result = append(result, u)
				continue
			}
			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)
			if afterOK && beforeOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Expired date filter
// ---------------------------------------------------------------------------

// ByExpiredDate returns a filter that keeps users whose expiry date is within
// the given range.
//
// Parameters:
//   - after: Only keep users expiring after this time (zero = no lower bound).
//   - before: Only keep users expiring before this time (zero = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are zero.
func ByExpiredDate(after, before time.Time) ModelFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			expStr := u.ExpiryString()
			if expStr == "" {
				result = append(result, u)
				continue
			}
			t, err := utils.ParseFlexibleDate(expStr)
			if err != nil {
				result = append(result, u)
				continue
			}
			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)
			if afterOK && beforeOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Renewed date filter
// ---------------------------------------------------------------------------

// ByRenewedDate returns a filter that keeps users whose renewal date is within
// the given range.
//
// Parameters:
//   - after: Only keep users renewed after this time (zero = no lower bound).
//   - before: Only keep users renewed before this time (zero = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are zero.
func ByRenewedDate(after, before time.Time) ModelFilter {
	if after.IsZero() && before.IsZero() {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			if u.RenewedAt == "" {
				result = append(result, u)
				continue
			}
			t, err := utils.ParseFlexibleDate(u.RenewedAt)
			if err != nil {
				result = append(result, u)
				continue
			}
			afterOK := after.IsZero() || !t.Before(after)
			beforeOK := before.IsZero() || !t.After(before)
			if afterOK && beforeOK {
				result = append(result, u)
			}
		}
		return result
	}
}
