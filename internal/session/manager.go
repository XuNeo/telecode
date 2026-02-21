package session

import (
	"sync"
)

// Manager manages session IDs per chat_id
type Manager struct {
	sessions map[int64]string
	mu       sync.RWMutex
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[int64]string),
	}
}

// Get returns the session ID for a chat_id
func (m *Manager) Get(chatID int64) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[chatID]
}

// Set saves the session ID for a chat_id
func (m *Manager) Set(chatID int64, sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[chatID] = sessionID
}

// Delete removes the session for a chat_id
func (m *Manager) Delete(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, chatID)
}

// Exists checks if a session exists
func (m *Manager) Exists(chatID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.sessions[chatID]
	return exists
}
