// =============================================================================
// FILE: internal/commands/scraper/prepare.go
// PURPOSE: Data preparation for scrape runs. Resolves usernames to user IDs,
//          fetches subscription lists, and applies user filters (blacklist,
//          price, activity). Ports Python data/models/selector.py and
//          utils/args/accessors/areas.py data loading.
// =============================================================================

package scraper

import (
	"context"
	"log/slog"
	"strings"

	"gofscraper/internal/app"
	"gofscraper/internal/model"
)

// ---------------------------------------------------------------------------
// PrepareData
// ---------------------------------------------------------------------------

// PrepareData resolves the list of users to scrape by fetching subscriptions
// and applying configured filters (user list, blacklist, price, activity).
//
// Parameters:
//   - ctx: Context for cancellation.
//   - a: The application instance providing config and session.
//   - logger: Structured logger.
//
// Returns:
//   - A filtered slice of User pointers ready for processing, and any error.
func PrepareData(ctx context.Context, a *app.App, logger *slog.Logger) ([]*model.User, error) {
	cfg := a.Config()

	// Step 1: Determine which usernames to fetch.
	userList := cfg.Advanced.DefaultUserList
	blackList := cfg.Advanced.DefaultBlackList

	logger.Info("preparing user data",
		"user_list", userList,
		"blacklist", blackList,
	)

	// Step 2: Fetch subscriptions.
	// TODO: Wire to API subscription fetcher based on retrieval mode.
	var allUsers []*model.User

	// Step 3: Apply filters.
	filtered := filterUsers(allUsers, userList, blackList)

	logger.Info("user preparation complete",
		"total_fetched", len(allUsers),
		"after_filter", len(filtered),
	)

	return filtered, nil
}

// ---------------------------------------------------------------------------
// User filtering
// ---------------------------------------------------------------------------

// filterUsers applies user list and blacklist filters to a set of users.
//
// Parameters:
//   - users: The full list of fetched users.
//   - userList: Usernames to include (empty = include all).
//   - blackList: Usernames to exclude.
//
// Returns:
//   - The filtered slice of users.
func filterUsers(users []*model.User, userList, blackList []string) []*model.User {
	blackSet := make(map[string]bool, len(blackList))
	for _, b := range blackList {
		blackSet[strings.ToLower(b)] = true
	}

	// If user list contains "ALL" or is empty, include all non-blacklisted.
	includeAll := len(userList) == 0
	if !includeAll {
		for _, u := range userList {
			if strings.ToUpper(u) == "ALL" {
				includeAll = true
				break
			}
		}
	}

	userSet := make(map[string]bool, len(userList))
	if !includeAll {
		for _, u := range userList {
			userSet[strings.ToLower(u)] = true
		}
	}

	var result []*model.User
	for _, user := range users {
		lower := strings.ToLower(user.Name)

		// Skip blacklisted users.
		if blackSet[lower] {
			continue
		}

		// If not including all, check the user list.
		if !includeAll && !userSet[lower] {
			continue
		}

		result = append(result, user)
	}

	return result
}
