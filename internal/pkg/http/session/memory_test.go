package session_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
)

func TestMemorySessionStore(t *testing.T) {
	stopChan := make(chan struct{})
	ttl := 2 * time.Second
	store := session.NewMemorySessionStore(ttl, stopChan)

	t.Run("Save and Retrieve Session", func(t *testing.T) {
		if err := store.Save("key1", "value1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		value, err := store.Session("key1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if value != "value1" {
			t.Errorf("expected value 'value1', got '%s'", value)
		}
	})

	t.Run("Prevent Duplicate Sessions", func(t *testing.T) {
		if err := store.Save("key1", "value2"); !errors.Is(err, session.ErrSessionExists) {
			t.Errorf("expected error '%v', got '%v'", session.ErrSessionExists, err)
		}
	})

	t.Run("Retrieve Non-Existent Session", func(t *testing.T) {
		_, err := store.Session("non-existent-key")
		if !errors.Is(err, session.ErrSessionDoesNotExist) {
			t.Errorf("expected error '%v', got '%v'", session.ErrSessionDoesNotExist, err)
		}
	})

	t.Run("Session Expiry", func(t *testing.T) {
		if err := store.Save("key2", "value2"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		time.Sleep(ttl + 500*time.Millisecond) // Wait for the session to expire

		_, err := store.Session("key2")
		if !errors.Is(err, session.ErrSessionExpired) {
			t.Errorf("expected error '%v', got '%v'", session.ErrSessionExpired, err)
		}
	})

	t.Run("Delete Session", func(t *testing.T) {
		if err := store.Save("key3", "value3"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := store.Destroy("key3"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err := store.Session("key3")
		if !errors.Is(err, session.ErrSessionDoesNotExist) {
			t.Errorf("expected error '%v', got '%v'", session.ErrSessionDoesNotExist, err)
		}
	})

	t.Run("Delete Non-Existent Session", func(t *testing.T) {
		if err := store.Destroy("non-existent-key"); !errors.Is(err, session.ErrSessionDoesNotExist) {
			t.Errorf("expected error '%v', got '%v'", session.ErrSessionDoesNotExist, err)
		}
	})

	close(stopChan)
}

func TestMemorySessionStoreConcurrentAccess(t *testing.T) {
	stopChan := make(chan struct{})
	store := session.NewMemorySessionStore(1*time.Second, stopChan)

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrently save sessions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sessionKey := "key" + string(rune(i))
			err := store.Save(sessionKey, "value"+string(rune(i)))
			if err != nil && !errors.Is(err, session.ErrSessionExists) {
				t.Errorf("unexpected error during save: %v", err)
			}
		}(i)
	}

	// Concurrently retrieve sessions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sessionKey := "key" + string(rune(i))
			_, err := store.Session(sessionKey)
			if err != nil && !errors.Is(err, session.ErrSessionDoesNotExist) {
				t.Errorf("unexpected error during retrieve: %v", err)
			}
		}(i)
	}

	// Concurrently delete sessions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sessionKey := "key" + string(rune(i))
			err := store.Destroy(sessionKey)
			if err != nil && !errors.Is(err, session.ErrSessionDoesNotExist) {
				t.Errorf("unexpected error during delete: %v", err)
			}
		}(i)
	}

	wg.Wait()
	close(stopChan)
}
