package auth

import (
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
	"github.com/ferdiebergado/goexpress"
)

func RegisterAuthRoutes(router *goexpress.Router, handler *Handler, sessMgr session.Manager) {
	router.Get("/signup", handler.HandleSignUp)
	router.Get("/signin", handler.HandleSignin)
	router.Get("/profile", handler.HandleProfile, goexpress.Middleware(RequireUserMiddleware(sessMgr)))

	router.Post("/api/signup", handler.HandleSignUpForm)
	router.Post("/api/signin", handler.HandleSignInForm)
}
