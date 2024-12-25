package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/assert"
)

func TestBaseHandler(t *testing.T) {
	cfg := config.Load()

	database := db.New(cfg.DB)
	conn, err := database.Connect(context.Background())

	if err != nil {
		t.Fatalf("can't connect to the database: %v", err)
	}

	defer database.Disconnect()

	router := goexpress.New()
	application := New(cfg, conn, router)
	application.SetupRouter()

	t.Run("GET /dbstats should return status 200 and render dbstats.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/dbstats", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		expected := "Database Statistics"
		actual := rec.Body.String()

		assert.Contains(t, actual, expected)
	})

	t.Run("GET /nonexistent should return status 404 and render 404.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		expected := "The page you are looking for does not exist."
		actual := rec.Body.String()

		assert.Contains(t, actual, expected)
	})
}
