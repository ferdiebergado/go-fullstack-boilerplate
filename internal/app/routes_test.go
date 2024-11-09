package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	router "github.com/ferdiebergado/go-express"
)

func TestAddRoutes(t *testing.T) {
	r := router.NewRouter() // Create a new instance of your custom Router
	AddRoutes(r)            // Add your routes to the Router

	t.Run("GET / should return status 200 and render home.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}

		// Check if the response contains content from home.html
		expected := "Welcome!"
		actual := rec.Body.String()

		if !strings.Contains(actual, expected) {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})

	t.Run("GET /nonexistent should return status 404 and render 404.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected %d but got %d", http.StatusNotFound, rec.Code)
		}

		expected := "The page you are looking for does not exist."
		actual := rec.Body.String()

		// Check if the response contains content from 404.html
		if !strings.Contains(actual, expected) {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})
}
