// =============================================================================
// FILE: internal/cli/flags/advanced_user_filter.go
// PURPOSE: Advanced user filter flag definitions: min-price, max-price,
//          last-seen-after, last-seen-before, sub-after, sub-before,
//          expired-after, expired-before, renewal-on, renewal-off,
//          promo-only, free-only, all-promo-only, regular-only.
// =============================================================================

package flags

import (
	"github.com/spf13/cobra"
)

// RegisterAdvancedUserFilterFlags adds advanced user filtering flags to the
// given command.
func RegisterAdvancedUserFilterFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.Float64("min-price", 0, "Minimum subscription price filter")
	f.Float64("max-price", 0, "Maximum subscription price filter (0 = no limit)")
	f.String("last-seen-after", "", "Only include users last seen after this date (YYYY-MM-DD)")
	f.String("last-seen-before", "", "Only include users last seen before this date (YYYY-MM-DD)")
	f.String("sub-after", "", "Only include users subscribed after this date")
	f.String("sub-before", "", "Only include users subscribed before this date")
	f.String("expired-after", "", "Only include users expired after this date")
	f.String("expired-before", "", "Only include users expired before this date")
	f.Bool("renewal-on", false, "Only include users with renewal enabled")
	f.Bool("renewal-off", false, "Only include users with renewal disabled")
	f.Bool("promo-only", false, "Only include users currently on promotion")
	f.Bool("free-only", false, "Only include free subscriptions")
	f.Bool("all-promo-only", false, "Only include users with any promo (current or past)")
	f.Bool("regular-only", false, "Only include regular (non-promo) subscriptions")
}
