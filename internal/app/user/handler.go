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
	params, err := request.JSON[SignUpParams](r)

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

	user, err := h.service.SignUp(r.Context(), params)

	if err != nil {
		var inputErr *validation.Error
		if errors.As(err, &inputErr) {
			valErr := errtypes.ValidationError(*inputErr)

			res := &response.APIResponse[any]{
				Message: valErr.Description,
				Errors:  inputErr.Errors,
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

			res := &response.APIResponse[validation.Error]{
				Message: "Invalid input!",
				Errors: validation.Errors{
					"email": {
						valErr.Description,
					},
				},
			}

			response.RenderError(w, valErr, res)
			return
		}

		serverError := errtypes.ServerError(err)

		res := &response.APIResponse[any]{
			Message: serverError.Error(),
		}

		response.RenderError(w, serverError, res)
		return
	}

	res := &response.APIResponse[*User]{
		Message: "Sign up successful!",
		Data:    &user,
	}

	slog.Debug("sending response", "message", res.Message, "data", res.Data)
	response.RenderJSON(w, http.StatusCreated, res)
}

func (h *Handler) HandleSignin(w http.ResponseWriter, _ *http.Request) {
	h.htmlTemplate.Render(w, nil, "signin.html")
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

	err = h.service.SignIn(r.Context(), params)

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

			res := &response.APIResponse[validation.Error]{
				Message: "Invalid input!",
				Errors: validation.Errors{
					"email": {
						valErr.Description,
					},
				},
			}

			response.RenderError(w, valErr, res)
			return
		}

		var inputErr *validation.Error
		if errors.As(err, &inputErr) {
			valErr := &errtypes.HTTPError{
				AppError: &errtypes.AppError{
					Description: inputErr.Error(),
					Err:         inputErr,
					Severity:    errtypes.Low,
				},
				Code: http.StatusUnauthorized,
			}

			res := &response.APIResponse[validation.Error]{
				Message: "Invalid input!",
				Errors:  inputErr.Errors,
			}

			response.RenderError(w, valErr, res)
			return
		}

		serverError := errtypes.ServerError(err)

		res := &response.APIResponse[any]{
			Message: serverError.Error(),
		}

		response.RenderError(w, serverError, res)
		return
	}

	res := &response.APIResponse[any]{
		Message: "Logged in.",
	}

	response.RenderJSON(w, http.StatusOK, res)
}
