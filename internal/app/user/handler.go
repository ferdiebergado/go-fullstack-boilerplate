package user

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/validation"
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/http/request"
)

type Handler struct {
	config       *config.Config
	router       *goexpress.Router
	service      AuthService
	htmlTemplate *html.Template
}

func NewHandler(cfg *config.Config, router *goexpress.Router, service AuthService, htmlTemplate *html.Template) *Handler {
	return &Handler{
		config:       cfg,
		router:       router,
		service:      service,
		htmlTemplate: htmlTemplate,
	}
}

func (h *Handler) HandleSignUp(w http.ResponseWriter, _ *http.Request) {
	h.htmlTemplate.Render(w, "signup.html", nil)
}

func (h *Handler) HandleSignUpForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[SignUpParams](r)

	if err != nil {
		response.RenderError(w, r, errtypes.BadRequest(err))
		return
	}

	user, err := h.service.SignUp(r.Context(), params)

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

	res := &response.APIResponse[User]{
		Message: "Sign up successful!",
		Data:    user,
	}

	slog.Debug("sending response", "message", res.Message, "data", res.Data)
	response.RenderJSON(w, http.StatusCreated, res)
}

func (h *Handler) HandleSignin(w http.ResponseWriter, _ *http.Request) {
	h.htmlTemplate.Render(w, "signin.html", nil)
}

func (h *Handler) HandleSignInForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[SignInParams](r)

	if err != nil {
		httpError := errtypes.BadRequest(err)
		response.RenderError(w, r, httpError)
		return
	}

	err = h.service.SignIn(r.Context(), params)

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

	res := &response.APIResponse[any]{
		Message: "Logged in.",
	}

	response.RenderJSON(w, http.StatusOK, res)
}
