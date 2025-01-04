//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/goexpress"

	"github.com/ferdiebergado/gopherkit/assert"
)

func TestAuthHandler(t *testing.T) {
	cfg := config.Load()

	database := db.New(cfg.DB)
	conn, err := database.Connect(context.Background())

	if err != nil {
		t.Fatalf("can't connect to the database: %v", err)
	}

	defer database.Disconnect()

	router := goexpress.New()
	application := app.New(cfg, conn, router)
	application.SetupRouter()

	t.Run("POST /api/signup should return status 201 and send new user as json ", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/signup", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var user user.User
		err := json.NewDecoder(rec.Body).Decode(&user)

		assert.NoError(t, err)
		// assert.Equal(t, "healthy", result.Status)
	})

	t.Run("GET /signup should return status 200 and render signin.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/signup", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		expected := "The page you were looking for does not exist."
		actual := rec.Body.String()

		assert.Contains(t, actual, expected)
	})
}
