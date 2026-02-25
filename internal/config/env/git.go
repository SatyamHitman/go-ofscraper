// =============================================================================
// FILE: internal/config/env/git.go
// PURPOSE: Git/HTTP request defaults for fetching dynamic rules from GitHub.
//          Ports Python of_env/values/req/git.py.
// =============================================================================

package env

// GitMinWait returns the minimum wait between git/rule fetch requests.
func GitMinWait() float64 {
	return GetFloat64("OF_GIT_MIN_WAIT", 1.0)
}

// GitMaxWait returns the maximum wait between git/rule fetch requests.
func GitMaxWait() float64 {
	return GetFloat64("OF_GIT_MAX_WAIT", 5.0)
}

// GitNumTries returns the retry count for git/rule fetch requests.
func GitNumTries() int {
	return GetInt("OF_GIT_NUM_TRIES", 3)
}
