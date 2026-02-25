// =============================================================================
// FILE: internal/filter/model_price.go
// PURPOSE: Model price filters. Filters users/models by their subscription
//          pricing (current, regular, promo, renewal). Ports Python
//          filters/models/price.py.
// =============================================================================

package filter

import (
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// Current price filter
// ---------------------------------------------------------------------------

// ByCurrentPrice returns a filter that keeps users whose current subscription
// price is within [minPrice, maxPrice].
//
// Parameters:
//   - minPrice: Minimum price (negative = no lower bound).
//   - maxPrice: Maximum price (negative = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are negative.
func ByCurrentPrice(minPrice, maxPrice float64) ModelFilter {
	if minPrice < 0 && maxPrice < 0 {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			price := u.FinalCurrentPrice()
			minOK := minPrice < 0 || price >= minPrice
			maxOK := maxPrice < 0 || price <= maxPrice
			if minOK && maxOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Regular price filter
// ---------------------------------------------------------------------------

// ByRegularPrice returns a filter that keeps users whose regular subscription
// price is within [minPrice, maxPrice].
//
// Parameters:
//   - minPrice: Minimum price (negative = no lower bound).
//   - maxPrice: Maximum price (negative = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are negative.
func ByRegularPrice(minPrice, maxPrice float64) ModelFilter {
	if minPrice < 0 && maxPrice < 0 {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			price := u.RegularPrice()
			minOK := minPrice < 0 || price >= minPrice
			maxOK := maxPrice < 0 || price <= maxPrice
			if minOK && maxOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Promo price filter
// ---------------------------------------------------------------------------

// ByPromoPrice returns a filter that keeps users whose lowest promo price
// is within [minPrice, maxPrice].
//
// Parameters:
//   - minPrice: Minimum price (negative = no lower bound).
//   - maxPrice: Maximum price (negative = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are negative.
func ByPromoPrice(minPrice, maxPrice float64) ModelFilter {
	if minPrice < 0 && maxPrice < 0 {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			price := u.FinalPromoPrice()
			minOK := minPrice < 0 || price >= minPrice
			maxOK := maxPrice < 0 || price <= maxPrice
			if minOK && maxOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Renewal price filter
// ---------------------------------------------------------------------------

// ByRenewalPrice returns a filter that keeps users whose renewal price
// is within [minPrice, maxPrice].
//
// Parameters:
//   - minPrice: Minimum price (negative = no lower bound).
//   - maxPrice: Maximum price (negative = no upper bound).
//
// Returns:
//   - A ModelFilter, or nil if both bounds are negative.
func ByRenewalPrice(minPrice, maxPrice float64) ModelFilter {
	if minPrice < 0 && maxPrice < 0 {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			price := u.FinalRenewalPrice()
			minOK := minPrice < 0 || price >= minPrice
			maxOK := maxPrice < 0 || price <= maxPrice
			if minOK && maxOK {
				result = append(result, u)
			}
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Free accounts filter
// ---------------------------------------------------------------------------

// ByFreeAccount returns a filter for free (price=0) vs paid accounts.
//
// Parameters:
//   - mode: "free" keeps only free, "paid" keeps only paid, "" = no filter.
//
// Returns:
//   - A ModelFilter, or nil if mode is empty.
func ByFreeAccount(mode string) ModelFilter {
	if mode == "" {
		return nil
	}

	return func(users []model.User) []model.User {
		var result []model.User
		for _, u := range users {
			isFree := u.FinalCurrentPrice() == 0
			switch mode {
			case "free":
				if isFree {
					result = append(result, u)
				}
			case "paid":
				if !isFree {
					result = append(result, u)
				}
			default:
				result = append(result, u)
			}
		}
		return result
	}
}
