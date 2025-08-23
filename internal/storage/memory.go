// Package storage internal/storage/memory.go
package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/mohar9h/goauth/internal/entity"
)

type memoryDriver struct {
	tokensByHash map[string]*entity.PersonalAccessToken // key is hashed token string
	tokensByID   map[int64]*entity.PersonalAccessToken  // key is token ID for O(1) lookups
	mu           sync.RWMutex
	nextID       int64 // Auto-incrementing ID
}

var _ Driver = (*memoryDriver)(nil)

func NewMemoryDriver() Driver {
	return &memoryDriver{
		tokensByHash: make(map[string]*entity.PersonalAccessToken),
		tokensByID:   make(map[int64]*entity.PersonalAccessToken),
		nextID:       1,
	}
}

// StoreToken stores the token using its hashed value as key
func (m *memoryDriver) StoreToken(t *entity.PersonalAccessToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Assign ID if not set
	if t.ID == 0 {
		t.ID = m.nextID
		m.nextID++
	}

	m.tokensByHash[t.Token] = t
	m.tokensByID[t.ID] = t
	return nil
}

// FindByID looks up token by its internal ID (numeric) - O(1) lookup
func (m *memoryDriver) FindByID(id int64) (*entity.PersonalAccessToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tok, ok := m.tokensByID[id]
	if !ok {
		return nil, errors.New("token not found")
	}

	if tok.ExpiresAt != nil && time.Now().After(*tok.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return tok, nil
}

// FindByHash looks up token by its hashed token string - O(1) lookup
func (m *memoryDriver) FindByHash(hash string) (*entity.PersonalAccessToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tok, ok := m.tokensByHash[hash]
	if !ok {
		return nil, errors.New("token not found")
	}

	if tok.ExpiresAt != nil && time.Now().After(*tok.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return tok, nil
}

// RevokeToken removes a token by its hashed value
func (m *memoryDriver) RevokeToken(hash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tok, ok := m.tokensByHash[hash]
	if !ok {
		return errors.New("token not found")
	}

	delete(m.tokensByHash, hash)
	delete(m.tokensByID, tok.ID)
	return nil
}

// TouchLastUsed updates the last used time for analytics or session freshness
func (m *memoryDriver) TouchLastUsed(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tok, ok := m.tokensByID[id]
	if !ok {
		return errors.New("token not found")
	}

	now := time.Now()
	tok.LastUsedAt = &now
	return nil
}
