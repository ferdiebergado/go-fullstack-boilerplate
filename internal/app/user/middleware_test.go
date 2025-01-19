package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
)

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
				req = req.WithContext(user.WithUser(req.Context(), tt.userID))
			}

			// Set Content-Type if specified
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Wrap the handler with AuthMiddleware and serve the request
			handler := user.AuthMiddleware(mockHandler)
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
