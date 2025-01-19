package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/auth"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
)

type MockSessionManager struct {
	FetchFunc         func(*http.Request) (*session.Data, error)
	SaveFunc          func(context.Context, string, session.Data) error
	DeleteSessionFunc func(*http.Request) error
	SessionIDFunc     func(*http.Request) (string, error)
}

func (m *MockSessionManager) Save(ctx context.Context, sessionID string, sessionData session.Data) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, sessionID, sessionData)
	}
	return nil
}

func (m *MockSessionManager) Fetch(r *http.Request) (*session.Data, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(r)
	}
	return &session.Data{}, nil
}

func (m *MockSessionManager) Destroy(r *http.Request) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(r)
	}
	return nil
}

func (m *MockSessionManager) SessionID(r *http.Request) (string, error) {
	if m.SessionIDFunc != nil {
		return m.SessionIDFunc(r)
	}
	return "", nil
}

func TestSessionMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		sessionCookie  *http.Cookie
		sessionManager func(*http.Request) (*session.Data, error)
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
			sessionManager: func(r *http.Request) (*session.Data, error) { return &session.Data{}, nil },
			expectedUserID: "",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "Valid session cookie",
			sessionCookie:  &http.Cookie{Name: "session", Value: "valid_session"},
			sessionManager: func(r *http.Request) (*session.Data, error) { return &session.Data{UserID: "12345"}, nil },
			expectedUserID: "12345",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock session manager with the provided function
			mockSessMgr := &MockSessionManager{
				FetchFunc: tt.sessionManager,
			}

			// Set up the middleware
			cfg := config.SessionConfig{
				SessionName: "session",
			}
			middleware := auth.SessionMiddleware(cfg, mockSessMgr)

			// Create a test handler to check if the middleware works correctly
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, err := auth.FromContext(r.Context())
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

// Mock response package
type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

func RenderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Test AuthMiddleware
func TestAuthMiddleware(t *testing.T) {
	// Mock handler to be wrapped by the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	cfg := config.Load()

	conn, err := db.Connect(context.Background(), cfg.DB)

	if err != nil {
		t.Fatal("failed to connect to the database")
	}

	tests := []struct {
		name           string
		withUser       bool   // Whether to include user in context
		userID         string // User ID to set in context
		contentType    string // Content-Type of the request
		expectedStatus int    // Expected HTTP status code
		expectedBody   string // Expected response body
		expectedHeader string // Expected redirect location header
	}{
		{
			name:           "Authorized request",
			withUser:       true,
			userID:         "12345",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			expectedHeader: "",
		},
		{
			name:           "Unauthorized JSON request",
			withUser:       false,
			contentType:    "application/json",
			expectedStatus: http.StatusFound,
			expectedBody:   `{"message":"Login is required to access this resource.","data":{"redirect_path":"/signin"}}`,
			expectedHeader: "",
		},
		{
			name:           "Unauthorized HTML request",
			withUser:       false,
			contentType:    "",
			expectedStatus: http.StatusSeeOther,
			expectedBody:   "",
			expectedHeader: "/signin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Add user to context if specified
			if tt.withUser {
				req = req.WithContext(auth.WithUser(req.Context(), tt.userID))
			}

			// Set Content-Type if specified
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Wrap the handler with AuthMiddleware and serve the request
			sessMgr := session.NewDatabaseSession(cfg.Session, conn)
			handler := auth.RequireUserMiddleware(sessMgr)(mockHandler)
			handler.ServeHTTP(rr, req)

			// Assert status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Assert response body for JSON requests
			if tt.contentType == "application/json" {
				var resp APIResponse[map[string]string]
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				if err != nil {
					t.Fatalf("failed to unmarshal JSON response: %v", err)
				}
				var expectedResp APIResponse[map[string]string]
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedResp)
				if err != nil {
					t.Fatalf("failed to unmarshal expected JSON: %v", err)
				}
				if resp.Message != expectedResp.Message {
					t.Errorf("expected response %+v, got %+v", expectedResp.Message, resp.Message)
				}
			}

			// Assert redirect location header for HTML requests
			if tt.contentType == "" && tt.expectedStatus == http.StatusSeeOther {
				location := rr.Header().Get("Location")
				if location != tt.expectedHeader {
					t.Errorf("expected Location header %q, got %q", tt.expectedHeader, location)
				}
			}

			// Assert response body for successful requests
			if tt.expectedStatus == http.StatusOK && rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
