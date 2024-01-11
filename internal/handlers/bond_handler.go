package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	handlerPort "kiramishima/m-backend/internal/core/ports/handlers"
	svcports "kiramishima/m-backend/internal/core/ports/services"
	httpErrors "kiramishima/m-backend/pkg/errors"
	httpUtils "kiramishima/m-backend/pkg/utils"
	"net/http"
	"os"
	"strconv"
)

var _ handlerPort.BondHandlers = (*BondHandlers)(nil)

// NewBondHandlers creates an instance of bond handlers
func NewBondHandlers(r *chi.Mux, logger *zap.SugaredLogger, s svcports.BondService, render *render.Render, validate *validator.Validate) {
	var tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SecretKey")), nil)

	handler := &BondHandlers{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/bonds", func(r chi.Router) {
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Get("/", handler.ListBondsHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/", handler.CreateBondHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Patch("/{id}", handler.UpdateBondHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Delete("/{id}", handler.DeleteBondHandler)
	})
}

type BondHandlers struct {
	logger   *zap.SugaredLogger
	service  svcports.BondService
	response *render.Render
	validate *validator.Validate
}

// ListBondsHandler
func (h *BondHandlers) ListBondsHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)

	h.logger.Info("UserID: ", UserID)
	ctx := req.Context()

	resp, err := h.service.ListBonds(ctx, UserID)
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, domain.WrapResponse[[]*domain.Bond]{Data: make([]*domain.Bond, 0)})
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadRequest})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, domain.WrapResponse[[]*domain.Bond]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// GetBondByUUIDHandler
func (h *BondHandlers) GetBondByUUIDHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)
	var BondID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)

	ctx := req.Context()

	resp, err := h.service.GetBondByID(ctx, UserID, int(BondID))
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, domain.WrapResponse[[]*domain.Bond]{Data: make([]*domain.Bond, 0)})
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadRequest})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, domain.WrapResponse[*domain.Bond]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// CreateBondHandler for register new users
func (h *BondHandlers) CreateBondHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)

	var form = &domain.BondRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}
	form.CreatedBy = UserID
	h.logger.Info(form)
	// Validate form
	err = form.Validate(h.validate)
	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: err.Error()})
		return
	}
	// context
	ctx := req.Context()

	err = h.service.CreateBond(ctx, form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrBeginTransaction) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			} else if errors.Is(err, httpErrors.ErrBondAlreadyExists) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBondAlreadyExists.Error()})
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrUserNotFound.Error()})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: err.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusAccepted, domain.SuccessResponse{Message: "Bond success created"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// UpdateBondHandler for register new users
func (h *BondHandlers) UpdateBondHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)
	var BondID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)
	var form = &domain.BondRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}

	form.CreatedBy = UserID
	h.logger.Info(form)
	ctx := req.Context()

	err = h.service.UpdateBond(ctx, int(BondID), form)
	if err != nil {
		// h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrInvalidRequestBody) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
			} else if errors.Is(err, httpErrors.ErrBeginTransaction) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			} else if errors.Is(err, httpErrors.ErrCommit) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusAccepted, domain.SuccessResponse{Message: "Success. Check you email for activate your account."}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// DeleteBondHandler for register new users
func (h *BondHandlers) DeleteBondHandler(w http.ResponseWriter, req *http.Request) {
	var BondID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)

	ctx := req.Context()

	err := h.service.DeleteBond(ctx, int(BondID))
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrInvalidRequestBody) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
			} else if errors.Is(err, httpErrors.ErrDeleteBond) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrDeleteBond.Error()})
			} else if errors.Is(err, httpErrors.ErrBeginTransaction) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBeginTransaction.Error()})
			} else if errors.Is(err, httpErrors.ErrCommit) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBeginTransaction.Error()})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusAccepted, domain.SuccessResponse{Message: "The bond was deleted"}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}
