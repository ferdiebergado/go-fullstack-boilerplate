package app

import (
	"database/sql"

	"github.com/ferdiebergado/goexpress"
)

// Set up the router
func SetupRouter(conn *sql.DB) *goexpress.Router {
	router := goexpress.New()
	application := New(conn, router)
	application.AddHandlers()
	return router
}
