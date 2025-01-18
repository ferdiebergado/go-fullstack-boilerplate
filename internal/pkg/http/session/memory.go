package session

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrSessionDoesNotExist = errors.New("session does not exist")
	ErrSessionExists       = errors.New("session already exists")
	ErrSessionExpired      = errors.New("session has expired")
)

type data struct {
	value     string
	createdAt time.Time
}

type sessions map[string]data

type MemorySessionStore struct {
	sessions sessions
	mu       sync.RWMutex
	ttl      time.Duration
	stopChan <-chan struct{}
}

func NewMemorySessionStore(ttl time.Duration, stopChan <-chan struct{}) Manager {
	store := &MemorySessionStore{
		sessions: make(sessions),
		ttl:      ttl,
		stopChan: stopChan,
	}
	store.StartCleanup(10 * time.Minute)
	return store
}

func (s *MemorySessionStore) Save(sessionKey, sessionData string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionKey]; exists {
		return ErrSessionExists
	}

	s.sessions[sessionKey] = data{
		value:     sessionData,
		createdAt: time.Now(),
	}

	return nil
}

func (s *MemorySessionStore) Session(sessionKey string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.sessions[sessionKey]
	if !ok {
		return "", ErrSessionDoesNotExist
	}

	if s.ttl > 0 && time.Since(item.createdAt) > s.ttl {
		return "", ErrSessionExpired
	}

	return item.value, nil
}

func (s *MemorySessionStore) DeleteSession(sessionKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionKey]; !exists {
		return ErrSessionDoesNotExist
	}

	delete(s.sessions, sessionKey)
	return nil
}

func (s *MemorySessionStore) cleanUpExpiredSessions() {
	s.mu.Lock()
	for key, data := range s.sessions {
		if s.ttl > 0 && time.Since(data.createdAt) > s.ttl {
			delete(s.sessions, key)
		}
	}
	s.mu.Unlock()
}

func (s *MemorySessionStore) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cleanUpExpiredSessions()
			case <-s.stopChan:
				return // Exit the goroutine
			}
		}
	}()
}
