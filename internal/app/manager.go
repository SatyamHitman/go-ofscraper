// =============================================================================
// FILE: internal/app/manager.go
// PURPOSE: Manager singleton that manages shared state across the application
//          lifetime. Holds references to active sessions, configs, and
//          processing state. Ports Python managers/manager.py.
// =============================================================================

package app

import (
	"sync"

	"gofscraper/internal/config"
	gohttp "gofscraper/internal/http"
)

// ---------------------------------------------------------------------------
// Manager singleton
// ---------------------------------------------------------------------------

var (
	globalManager *Manager
	managerOnce   sync.Once
)

// Manager holds shared references and state across the application lifetime.
type Manager struct {
	mu sync.RWMutex

	// session is the active HTTP session manager.
	session *gohttp.SessionManager

	// cfg is the active configuration.
	cfg *config.AppConfig

	// activeUsers tracks usernames currently being processed.
	activeUsers map[string]bool
}

// GetManager returns the global Manager singleton, creating it if needed.
//
// Returns:
//   - The global Manager instance.
func GetManager() *Manager {
	managerOnce.Do(func() {
		globalManager = &Manager{
			activeUsers: make(map[string]bool),
		}
	})
	return globalManager
}

// ---------------------------------------------------------------------------
// Session management
// ---------------------------------------------------------------------------

// SetSession updates the active HTTP session.
//
// Parameters:
//   - session: The new session manager.
func (m *Manager) SetSession(session *gohttp.SessionManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.session = session
}

// Session returns the active HTTP session.
//
// Returns:
//   - The current SessionManager, or nil if not set.
func (m *Manager) Session() *gohttp.SessionManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session
}

// ---------------------------------------------------------------------------
// Config management
// ---------------------------------------------------------------------------

// SetConfig updates the active configuration reference.
//
// Parameters:
//   - cfg: The new configuration.
func (m *Manager) SetConfig(cfg *config.AppConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg = cfg
}

// Config returns the active configuration.
//
// Returns:
//   - The current AppConfig, or nil if not set.
func (m *Manager) Config() *config.AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg
}

// ---------------------------------------------------------------------------
// Active user tracking
// ---------------------------------------------------------------------------

// MarkUserActive records that a user is currently being processed.
//
// Parameters:
//   - username: The username to mark as active.
//
// Returns:
//   - true if the user was not already active.
func (m *Manager) MarkUserActive(username string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeUsers[username] {
		return false
	}
	m.activeUsers[username] = true
	return true
}

// MarkUserDone records that a user is no longer being processed.
//
// Parameters:
//   - username: The username to mark as done.
func (m *Manager) MarkUserDone(username string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeUsers, username)
}

// IsUserActive checks if a user is currently being processed.
//
// Parameters:
//   - username: The username to check.
//
// Returns:
//   - true if the user is currently active.
func (m *Manager) IsUserActive(username string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeUsers[username]
}

// ActiveUserCount returns the number of users currently being processed.
//
// Returns:
//   - The count of active users.
func (m *Manager) ActiveUserCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.activeUsers)
}
