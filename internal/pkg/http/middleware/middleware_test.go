package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/middleware"
)

type MockSessionManager struct {
	SessionFunc       func(sessionID string) (string, error)
	SaveFunc          func(sessionKey, sessionData string) error
	DeleteSessionFunc func(sessionKey string) error
}

func (m *MockSessionManager) Save(sessionKey, sessionData string) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(sessionKey, sessionData)
	}
	return nil
}

func (m *MockSessionManager) Session(sessionID string) (string, error) {
	if m.SessionFunc != nil {
		return m.SessionFunc(sessionID)
	}
	return "", nil
}

func (m *MockSessionManager) Destroy(sessionKey string) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(sessionKey)
	}
	return nil
}

func TestSessionMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		sessionCookie  *http.Cookie
		sessionManager func(string) (string, error)
		expectedUserID string
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "No session cookie",
			sessionCookie:  nil,
			sessionManager: nil,
			expectedUserID: "",
			expectedStatus: http.StatusOK,
			wantErr:        true,
		},
		{
			name:           "Invalid session cookie",
			sessionCookie:  &http.Cookie{Name: "session", Value: "invalid_session"},
			sessionManager: func(sessionID string) (string, error) { return "", nil },
			expectedUserID: "",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "Valid session cookie",
			sessionCookie:  &http.Cookie{Name: "session", Value: "valid_session"},
			sessionManager: func(sessionID string) (string, error) { return "12345", nil },
			expectedUserID: "12345",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock session manager with the provided function
			mockSessMgr := &MockSessionManager{
				SessionFunc: tt.sessionManager,
			}

			// Set up the middleware
			cfg := config.HTTPServerConfig{
				SessionName: "session",
			}
			middleware := middleware.SessionMiddleware(cfg, mockSessMgr)

			// Create a test handler to check if the middleware works correctly
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, err := user.FromContext(r.Context())
				if (err != nil) != tt.wantErr {
					t.Errorf("FromContext() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if userID != tt.expectedUserID {
					t.Errorf("expected userID %s, got %s.", tt.expectedUserID, userID)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Wrap the handler with the middleware
			middlewareHandler := middleware(handler)

			// Create a request
			req := httptest.NewRequest("GET", "/", nil)
			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the handler
			middlewareHandler.ServeHTTP(rr, req)

			// Assert the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
