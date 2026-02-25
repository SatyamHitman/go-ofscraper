// =============================================================================
// FILE: internal/model/retriever.go
// PURPOSE: Defines the subscription retrieval strategy interface and helpers.
//          Ports the data/models/utils/retriver.py retrieval logic that
//          determines how to fetch the list of subscribed models (individual
//          lookup vs. full list fetch). The actual API calls live in the api
//          package; this file provides the strategy selection logic.
// =============================================================================

package model

// ---------------------------------------------------------------------------
// RetrievalMode enumerates the strategies for fetching subscriptions.
// ---------------------------------------------------------------------------

// RetrievalMode specifies how subscription data should be fetched from the API.
type RetrievalMode string

const (
	// RetrievalModeList fetches all active + expired subscriptions in bulk.
	RetrievalModeList RetrievalMode = "list"

	// RetrievalModeIndividual fetches each specified username one at a time.
	RetrievalModeIndividual RetrievalMode = "individual"

	// RetrievalModeMainList fetches the full subscription list without filtering.
	RetrievalModeMainList RetrievalMode = "main_list"
)

// ---------------------------------------------------------------------------
// RetrievalStrategy determines the optimal fetch mode.
// ---------------------------------------------------------------------------

// DetermineRetrievalMode selects the optimal subscription retrieval strategy
// based on the configured usernames, subscription counts, and user preferences.
//
// Decision logic (mirrors Python retriver.py):
//   1. If anonymous mode → Individual
//   2. If allMainModels flag → MainList
//   3. If no usernames specified → List
//   4. If "ALL" in usernames → List
//   5. If username_search == "individual" → Individual
//   6. If username_search == "list" → List
//   7. If (activeCount + expiredCount) / 12 >= len(usernames) → Individual
//   8. Otherwise → List
//
// Parameters:
//   - usernames: The list of specified usernames to scrape.
//   - usernameSearch: The configured search strategy ("individual", "list", or "").
//   - activeCount: The number of active subscriptions.
//   - expiredCount: The number of expired subscriptions.
//   - isAnon: Whether running in anonymous mode.
//   - allMainModels: Whether the allMainModels flag is set.
//
// Returns:
//   - The selected RetrievalMode.
func DetermineRetrievalMode(
	usernames []string,
	usernameSearch string,
	activeCount int,
	expiredCount int,
	isAnon bool,
	allMainModels bool,
) RetrievalMode {
	// Anonymous mode always uses individual lookups
	if isAnon {
		return RetrievalModeIndividual
	}

	// All main models flag forces full list retrieval
	if allMainModels {
		return RetrievalModeMainList
	}

	// No usernames specified — fetch everything
	if len(usernames) == 0 {
		return RetrievalModeList
	}

	// "ALL" keyword present — fetch full list
	for _, u := range usernames {
		if u == "ALL" {
			return RetrievalModeList
		}
	}

	// Explicit search mode preference
	switch usernameSearch {
	case "individual":
		return RetrievalModeIndividual
	case "list":
		return RetrievalModeList
	}

	// Heuristic: if total subs / 12 >= number of usernames, individual is cheaper
	totalCount := activeCount + expiredCount
	if totalCount > 0 && len(usernames) > 0 {
		pagesNeeded := totalCount / 12
		if pagesNeeded >= len(usernames) {
			return RetrievalModeIndividual
		}
	}

	// Default to list mode
	return RetrievalModeList
}

// ---------------------------------------------------------------------------
// Subscription count helper
// ---------------------------------------------------------------------------

// SubCounts holds the active and expired subscription counts for a user.
type SubCounts struct {
	Active  int `json:"active"`
	Expired int `json:"expired"`
}

// Total returns the sum of active and expired subscriptions.
//
// Returns:
//   - The total subscription count.
func (sc SubCounts) Total() int {
	return sc.Active + sc.Expired
}
