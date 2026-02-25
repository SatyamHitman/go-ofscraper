package version

import "fmt"

// Version info set via ldflags at build time.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("gofscraper %s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
