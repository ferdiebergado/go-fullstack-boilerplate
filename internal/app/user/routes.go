package user

import (
	"github.com/ferdiebergado/goexpress"
)

func RegisterAuthRoutes(router *goexpress.Router, handler *Handler) {
	router.Get("/signup", handler.HandleSignUp)
	router.Post("/api/signup", handler.HandleSignUpForm)
	router.Post("/api/signin", handler.HandleSignInForm)
}
