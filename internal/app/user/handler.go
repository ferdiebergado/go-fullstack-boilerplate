package user

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

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
	h.htmlTemplate.Render(w, nil, "signup.html")
}

func (h *Handler) HandleSignUpForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[BasicAuthParams](r)

	if err != nil {
		httpError := errtypes.HTTPError{
			AppError: &errtypes.AppError{
				Description: "cannot read auth request",
				Err:         err,
				Severity:    errtypes.Low,
			},
			Code: http.StatusBadRequest,
		}
		response.RenderError[any](w, &httpError, nil)
		return
	}

	authData := AuthData{
		Email:                params.Email,
		Password:             params.Password,
		PasswordConfirmation: params.PasswordConfirmation,
		AuthMethod:           Basic,
	}

	user, err := h.service.SignUp(r.Context(), authData)

	if err != nil {
		var inputErr *validation.InputError
		if errors.As(err, &inputErr) {
			valErr := errtypes.ValidationError(*inputErr)

			res := &response.APIResponse[validation.InputError]{
				Success: false,
				Message: valErr.Description,
				Data:    inputErr,
			}

			response.RenderError(w, valErr, res)
			return
		}

		if errors.Is(err, ErrEmailExists) {
			valErr := &errtypes.HTTPError{
				AppError: &errtypes.AppError{
					Description: strings.Split(err.Error(), ":")[0],
					Err:         errors.Unwrap(err),
					Severity:    errtypes.Low,
				},
				Code: http.StatusUnprocessableEntity,
			}

			res := &response.APIResponse[validation.InputError]{
				Success: false,
				Message: "Invalid input!",
				Data: &validation.InputError{
					Errors: map[string][]string{
						"email": {
							valErr.Description,
						},
					},
				},
			}

			response.RenderError(w, valErr, res)
			return
		}

		serverError := errtypes.ServerError(err)

		res := &response.APIResponse[any]{
			Success: false,
			Message: serverError.Error(),
		}

		response.RenderError(w, serverError, res)
		return
	}
	res := &response.APIResponse[User]{
		Success: true,
		Message: "Sign up successful!",
		Data:    &user,
	}

	slog.Debug("sending response", slog.Bool("success", res.Success), "message", res.Message, "data", res.Data)
	response.RenderJSON(w, http.StatusCreated, res)
}

func (h *Handler) HandleSignInForm(w http.ResponseWriter, r *http.Request) {
	params, err := request.JSON[SignInParams](r)

	if err != nil {
		httpError := errtypes.HTTPError{
			AppError: &errtypes.AppError{
				Description: "cannot read auth request",
				Err:         err,
				Severity:    errtypes.Low,
			},
			Code: http.StatusBadRequest,
		}
		response.RenderError[any](w, &httpError, nil)
		return
	}

	authData := AuthData{
		Email:      params.Email,
		Password:   params.Password,
		AuthMethod: Basic,
	}

	err = h.service.SignIn(r.Context(), authData)

	if err != nil {
		if errors.Is(err, ErrUserPassInvalid) {
			innerErr := errors.Unwrap(err)
			valErr := &errtypes.HTTPError{
				AppError: &errtypes.AppError{
					Description: innerErr.Error(),
					Err:         innerErr,
					Severity:    errtypes.Low,
				},
				Code: http.StatusUnauthorized,
			}

			res := &response.APIResponse[validation.InputError]{
				Success: false,
				Message: "Invalid input!",
				Data: &validation.InputError{
					Errors: map[string][]string{
						"email": {
							valErr.Description,
						},
					},
				},
			}
			response.RenderError(w, valErr, res)
			return
		}

		serverError := errtypes.ServerError(err)

		res := &response.APIResponse[any]{
			Success: false,
			Message: serverError.Error(),
		}

		response.RenderError(w, serverError, res)
		return
	}

	res := &response.APIResponse[any]{
		Success: true,
		Message: "Logged in.",
	}

	response.RenderJSON(w, http.StatusOK, res)
}
