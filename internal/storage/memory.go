package storage

import (
	"errors"
	"github.com/mohar9h/goauth/internal/entity"
	"sync"
	"time"
)

type memoryDriver struct {
	tokens map[string]*entity.Token // key is hashed token string
	mu     sync.RWMutex
}

var _ Driver = (*memoryDriver)(nil)

func NewMemoryDriver() Driver {
	return &memoryDriver{
		tokens: make(map[string]*entity.Token),
	}
}

// StoreToken stores the token using its hashed value as key
func (m *memoryDriver) StoreToken(t *entity.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[t.Token] = t
	return nil
}

// FindByID looks up token by its internal ID (numeric)
func (m *memoryDriver) FindByID(id int64) (*entity.Token, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.tokens {
		if t.ID == id {
			if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
				return nil, errors.New("t expired")
			}
			return t, nil
		}
	}
	return nil, errors.New("t not found")
}

// FindByHash looks up token by its hashed token string
func (m *memoryDriver) FindByHash(hash string) (*entity.Token, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tok, ok := m.tokens[hash]
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
	delete(m.tokens, hash)
	return nil
}

// TouchLastUsed updates the last used time for analytics or session freshness
func (m *memoryDriver) TouchLastUsed(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, tok := range m.tokens {
		if tok.ID == id {
			now := time.Now()
			tok.LastUsedAt = &now
			return nil
		}
	}
	return errors.New("token not found")
}
