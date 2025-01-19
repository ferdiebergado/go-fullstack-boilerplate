package user

import (
	"github.com/ferdiebergado/goexpress"
)

func RegisterAuthRoutes(router *goexpress.Router, handler *Handler) {
	router.Get("/signup", handler.HandleSignUp)
	router.Get("/signin", handler.HandleSignin)
	router.Get("/profile", handler.HandleProfile, AuthMiddleware)

	router.Post("/api/signup", handler.HandleSignUpForm)
	router.Post("/api/signin", handler.HandleSignInForm)
}
