package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/http/request"
)

type Handler struct {
	config  *config.Config
	router  *goexpress.Router
	service Service
	tmpl    *html.Template
	sessMgr session.Manager
}

func NewHandler(cfg *config.Config, router *goexpress.Router, service Service, tmpl *html.Template, sessMgr session.Manager) *Handler {
	return &Handler{
		config:  cfg,
		router:  router,
		service: service,
		tmpl:    tmpl,
		sessMgr: sessMgr,
	}
}

func (h *Handler) HandleSignUp(w http.ResponseWriter, _ *http.Request) {
	h.tmpl.Render(w, "signup.html", nil)
}

func (h *Handler) HandleSignUpForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[SignUpParams](r)

	if err != nil {
		response.RenderError(w, r, errtypes.BadRequest(err))
		return
	}

	u, err := h.service.SignUp(r.Context(), params)

	if err != nil {
		var inputErr *validation.Error
		if errors.As(err, &inputErr) {
			valErr := errtypes.ValidationError(*inputErr)

			response.RenderError(w, r, valErr)
			return
		}

		var emailErr *EmailExistsError
		if errors.As(err, &emailErr) {
			valErr := errtypes.ValidationError(
				validation.Error{
					Errors: validation.Errors{
						"email": {
							emailErr.Error(),
						},
					},
				})

			response.RenderError(w, r, valErr)
			return
		}

		serverError := errtypes.ServerError(err)
		response.RenderError(w, r, serverError)
		return
	}

	res := &response.APIResponse[user.User]{
		Message: "Sign up successful!",
		Data:    u,
	}

	slog.Debug("sending response", "message", res.Message, "data", res.Data)
	response.RenderJSON(w, http.StatusCreated, res)
}

func (h *Handler) HandleSignin(w http.ResponseWriter, _ *http.Request) {
	h.tmpl.Render(w, "signin.html", nil)
}

func (h *Handler) HandleSignInForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[SignInParams](r)

	if err != nil {
		httpError := errtypes.BadRequest(err)
		response.RenderError(w, r, httpError)
		return
	}

	userID, err := h.service.SignIn(r.Context(), params)

	if err != nil {
		var inputErr *validation.Error
		if errors.As(err, &inputErr) {
			valErr := errtypes.ValidationError(*inputErr)
			response.RenderError(w, r, valErr)
			return
		}

		if errors.Is(err, ErrUserPassInvalid) {
			innerErr := errors.Unwrap(err)
			authErr := errtypes.AuthenticationError(innerErr)
			response.RenderError(w, r, authErr)
			return
		}

		serverError := errtypes.ServerError(err)
		response.RenderError(w, r, serverError)
		return
	}

	sessData, err := h.sessMgr.Get(r)

	if err != nil {
		serverError := errtypes.ServerError(err)
		response.RenderError(w, r, serverError)
		return
	}

	if sessData != nil {
		if err := h.sessMgr.Destroy(r.Context(), sessData.ID()); err != nil {
			serverError := errtypes.ServerError(err)
			response.RenderError(w, r, serverError)
			return
		}
	}

	sess, err := h.sessMgr.Start(r)

	if err != nil {
		serverError := errtypes.ServerError(err)
		response.RenderError(w, r, serverError)
		return
	}

	var redirectURL string

	if err := sess.SetUserID(userID); err != nil {
		serverError := errtypes.ServerError(err)
		response.RenderError(w, r, serverError)
		return
	}

	intendedURL, ok := sess.Flash("intendedUrl").(string)

	if !ok {
		redirectURL = "/dashboard"
	} else {
		redirectURL = intendedURL
	}

	res := &response.APIResponse[map[string]string]{
		Message: "Logged in.",
		Data: &map[string]string{
			"redirectUrl": redirectURL,
		},
	}

	response.RenderJSON(w, http.StatusOK, res)
}

func (h *Handler) HandleProfile(w http.ResponseWriter, _ *http.Request) {
	h.tmpl.Render(w, "profile.html", nil)
}
