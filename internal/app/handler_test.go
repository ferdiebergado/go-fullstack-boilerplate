package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/gopherkit/env"
)

func TestBaseHandler(t *testing.T) {
	if err := env.Load("../../.env"); err != nil {
		t.Fatal("missing .env file")
	}

	// Connect to the database.
	conn, err := db.Connect(context.Background(), config.Load().DB)

	if err != nil {
		t.Fatalf("db connect: %v", err)
	}

	// Close the database connection after running the application
	defer db.Disconnect(conn)

	router := SetupRouter(conn)

	t.Run("GET / should return status 200 and render home.html", func(t *testing.T) {
		t.Skip()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}

		expected := "Welcome!"
		actual := rec.Body.String()

		if !strings.Contains(actual, expected) {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})

	t.Run("GET /dbstats should return status 200 and render dbstats.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dbstats", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected %d but got %d", http.StatusOK, rec.Code)
		}

		expected := "Database Statistics"
		actual := rec.Body.String()

		if !strings.Contains(actual, expected) {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})

	t.Run("GET /nonexistent should return status 404 and render 404.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected %d but got %d", http.StatusNotFound, rec.Code)
		}

		expected := "The page you are looking for does not exist."
		actual := rec.Body.String()

		if !strings.Contains(actual, expected) {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	})
}
