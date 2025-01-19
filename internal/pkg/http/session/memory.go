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

type session struct {
	value     string
	createdAt time.Time
}

type sessions map[string]session

type InMemorySession struct {
	sessions sessions
	ttl      time.Duration
	mu       sync.RWMutex
}

func NewInMemorySession(ttl time.Duration) Manager {
	return &InMemorySession{
		sessions: make(sessions),
		ttl:      ttl,
	}
}

func (s *InMemorySession) Save(sessionKey, sessionData string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionKey]; exists {
		return ErrSessionExists
	}

	s.sessions[sessionKey] = session{
		value:     sessionData,
		createdAt: time.Now(),
	}

	return nil
}

func (s *InMemorySession) Session(sessionKey string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionKey]
	if !ok {
		return "", ErrSessionDoesNotExist
	}

	if s.ttl > 0 && time.Since(session.createdAt) > s.ttl {
		return "", ErrSessionExpired
	}

	return session.value, nil
}

func (s *InMemorySession) Destroy(sessionKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionKey]; !exists {
		return ErrSessionDoesNotExist
	}

	delete(s.sessions, sessionKey)
	return nil
}

func (s *InMemorySession) cleanUpExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, session := range s.sessions {
		if s.ttl > 0 && time.Since(session.createdAt) > s.ttl {
			delete(s.sessions, key)
		}
	}
}

func (s *InMemorySession) StartCleanup(wg *sync.WaitGroup, stopChan <-chan struct{}, interval time.Duration) {
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cleanUpExpiredSessions()
			case <-stopChan:
				return
			}
		}
	}()
}
