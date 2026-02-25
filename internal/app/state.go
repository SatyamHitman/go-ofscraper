// =============================================================================
// FILE: internal/app/state.go
// PURPOSE: State tracks the current processing state including which user is
//          being processed, which content area is active, and overall progress.
//          Ports Python data/state.py processing state tracking.
// =============================================================================

package app

import (
	"sync"
)

// ---------------------------------------------------------------------------
// ProcessingPhase enumerates the phases of a scrape run.
// ---------------------------------------------------------------------------

// ProcessingPhase identifies the current stage of processing.
type ProcessingPhase string

const (
	PhaseIdle       ProcessingPhase = "idle"
	PhaseInit       ProcessingPhase = "initializing"
	PhasePrepare    ProcessingPhase = "preparing"
	PhaseFetching   ProcessingPhase = "fetching"
	PhaseFiltering  ProcessingPhase = "filtering"
	PhaseDownload   ProcessingPhase = "downloading"
	PhaseLiking     ProcessingPhase = "liking"
	PhaseMetadata   ProcessingPhase = "metadata"
	PhaseCleanup    ProcessingPhase = "cleanup"
	PhaseDone       ProcessingPhase = "done"
)

// ---------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------

// State tracks the current processing state of the application. All fields
// are protected by a mutex for concurrent access.
type State struct {
	mu sync.RWMutex

	// Current phase of processing.
	phase ProcessingPhase

	// Current user being processed.
	currentUser string

	// Current content area being processed.
	currentArea string

	// Progress tracking.
	currentIndex int
	totalCount   int

	// Error from the most recent operation.
	lastError error
}

// NewState creates a State in the idle phase.
//
// Returns:
//   - An initialized State.
func NewState() *State {
	return &State{
		phase: PhaseIdle,
	}
}

// Phase returns the current processing phase.
func (s *State) Phase() ProcessingPhase {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.phase
}

// SetPhase updates the current processing phase.
//
// Parameters:
//   - phase: The new phase.
func (s *State) SetPhase(phase ProcessingPhase) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.phase = phase
}

// CurrentUser returns the username currently being processed.
func (s *State) CurrentUser() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentUser
}

// SetCurrentUser updates the current user.
//
// Parameters:
//   - username: The username now being processed.
func (s *State) SetCurrentUser(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentUser = username
}

// CurrentArea returns the content area currently being processed.
func (s *State) CurrentArea() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentArea
}

// SetCurrentArea updates the current content area.
//
// Parameters:
//   - area: The area now being processed.
func (s *State) SetCurrentArea(area string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentArea = area
}

// Progress returns the current progress index and total count.
//
// Returns:
//   - current: The current item index (0-based).
//   - total: The total number of items.
func (s *State) Progress() (current, total int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentIndex, s.totalCount
}

// SetProgress updates the progress counters.
//
// Parameters:
//   - current: The current item index.
//   - total: The total number of items.
func (s *State) SetProgress(current, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentIndex = current
	s.totalCount = total
}

// LastError returns the most recent error, if any.
func (s *State) LastError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastError
}

// SetLastError records an error.
//
// Parameters:
//   - err: The error to record.
func (s *State) SetLastError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastError = err
}

// Reset returns the state to idle with all fields cleared.
func (s *State) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.phase = PhaseIdle
	s.currentUser = ""
	s.currentArea = ""
	s.currentIndex = 0
	s.totalCount = 0
	s.lastError = nil
}
