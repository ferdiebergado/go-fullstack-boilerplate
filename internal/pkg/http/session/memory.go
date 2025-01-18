package session

import (
	"errors"
	"sync"
)

var (
	ErrSessionDoesNotExist = errors.New("session does not exist")
	ErrSessionExists       = errors.New("session already exists")
)

type sessions map[string]string

type MemorySessionStore struct {
	sessions sessions
	mu       sync.RWMutex
}

func NewMemorySessionStore() Manager {
	return &MemorySessionStore{sessions: make(sessions)}
}

func (s *MemorySessionStore) Save(sessionKey, sessionData string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.sessions[sessionKey]
	if exists {
		return ErrSessionExists
	}

	s.sessions[sessionKey] = sessionData

	return nil
}

func (s *MemorySessionStore) Session(sessionKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.sessions[sessionKey]
	if !ok {
		return "", ErrSessionDoesNotExist
	}

	return item, nil
}

func (s *MemorySessionStore) DeleteSession(sessionKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.sessions[sessionKey]

	if !exists {
		return ErrSessionDoesNotExist
	}

	delete(s.sessions, sessionKey)
	return nil
}
