package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	handlerPort "kiramishima/m-backend/internal/core/ports/handlers"
	svcports "kiramishima/m-backend/internal/core/ports/services"
	httpErrors "kiramishima/m-backend/pkg/errors"

	httpUtils "kiramishima/m-backend/pkg/utils"
	"net/http"
)

var _ handlerPort.AuthHandlers = (*AuthHandlers)(nil)

// NewAuthHandlers creates a instance of auth handlers
func NewAuthHandlers(r *chi.Mux, logger *zap.SugaredLogger, s svcports.AuthService, render *render.Render, validate *validator.Validate) {
	handler := &AuthHandlers{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/sign-in", handler.SignInHandler)
		r.Post("/sign-up", handler.SignUpHandler)
	})
}

type AuthHandlers struct {
	logger   *zap.SugaredLogger
	service  svcports.AuthService
	response *render.Render
	validate *validator.Validate
}

func (h *AuthHandlers) SignInHandler(w http.ResponseWriter, req *http.Request) {
	var form = &domain.AuthRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}
	h.logger.Info(form)
	// Validate Form
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}

	ctx := req.Context()

	resp, err := h.service.FindByCredentials(ctx, form)
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrInvalidRequestBody) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.BadQueryParams.Error()})
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadEmailOrPassword.Error()})
			} else if errors.Is(err, httpErrors.ErrBadPassword) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadPassword.Error()})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, resp); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// SignUpHandler for register new users
func (h *AuthHandlers) SignUpHandler(w http.ResponseWriter, req *http.Request) {
	var form = &domain.RegisterRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}
	h.logger.Info(form)
	// Validate form
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: fmt.Sprintf("Validation error: %s", err)})
		return
	}
	ctx := req.Context()

	err = h.service.Register(ctx, form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrInvalidRequestBody) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadEmailOrPassword.Error()})
			} else if errors.Is(err, httpErrors.ErrAlreadyExists) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrAlreadyExists.Error()})
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrUserNotFound.Error()})
			} else {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadEmailOrPassword.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, domain.SuccessResponse{Message: "Success. Check you email for activate your account."}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}
