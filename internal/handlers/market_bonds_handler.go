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

var _ handlerPort.MarketBondsHandlers = (*MarketBondsHandlers)(nil)

// NewMarketBondsHandlers creates an instance of market bonds handlers
func NewMarketBondsHandlers(r *chi.Mux, logger *zap.SugaredLogger, s svcports.MarketBondsService, render *render.Render, validate *validator.Validate) {
	var tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SecretKey")), nil)

	handler := &MarketBondsHandlers{
		logger:   logger,
		service:  s,
		response: render,
		validate: validate,
	}

	r.Route("/v1/market", func(r chi.Router) {
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Get("/", handler.ListMarketBondsHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/{id}", handler.GetMarketBondByIDHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/{id}/buy", handler.BuyMarketBondHandler)
		r.With(jwtauth.Verifier(tokenAuth)).With(jwtauth.Authenticator(tokenAuth)).Post("/sell", handler.SellMarketBondHandler)
	})
}

type MarketBondsHandlers struct {
	logger   *zap.SugaredLogger
	service  svcports.MarketBondsService
	response *render.Render
	validate *validator.Validate
}

// ListBondsHandler
func (h *MarketBondsHandlers) ListMarketBondsHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)

	h.logger.Info("UserID: ", UserID)
	ctx := req.Context()

	resp, err := h.service.ListMarketBonds(ctx, UserID)
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, domain.WrapResponse[[]*domain.MarketBond]{Data: make([]*domain.MarketBond, 0)})
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadRequest})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, domain.WrapResponse[[]*domain.MarketBond]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// GetBondByUUIDHandler
func (h *MarketBondsHandlers) GetMarketBondByIDHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)
	var MarketBondID, _ = strconv.ParseInt(chi.URLParam(req, "id"), 10, 10)

	ctx := req.Context()

	resp, err := h.service.GetMarketBondByID(ctx, UserID, int(MarketBondID))
	if err != nil {
		h.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			_ = h.response.JSON(w, http.StatusGatewayTimeout, domain.ErrorResponse{ErrorMessage: httpErrors.ErrTimeout.Error()})
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				_ = h.response.JSON(w, http.StatusOK, domain.WrapResponse[*domain.MarketBond]{Data: nil})
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrBadRequest})
			} else {
				_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
			}
		}
		return
	}

	if err := h.response.JSON(w, http.StatusOK, domain.WrapResponse[*domain.MarketBond]{Data: resp}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// BuyMarketBondHandler for register new users
func (h *MarketBondsHandlers) BuyMarketBondHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)
	var form = &domain.MarketBondRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}

	form.BuyerID = UserID
	h.logger.Info(form)
	ctx := req.Context()

	err = h.service.BuyMarketBond(ctx, form)
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

	if err := h.response.JSON(w, http.StatusAccepted, domain.SuccessResponse{Message: "Success. Your bought is in process."}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}

// SellMarketBondHandler for register new users
func (h *MarketBondsHandlers) SellMarketBondHandler(w http.ResponseWriter, req *http.Request) {
	var UserID = httpUtils.GetUserIDInJWTHeader(req)
	var form = &domain.MarketSellRequest{}

	err := httpUtils.ReadJSON(w, req, &form)

	if err != nil {
		h.logger.Error(err.Error())
		_ = h.response.JSON(w, http.StatusBadRequest, domain.ErrorResponse{ErrorMessage: httpErrors.ErrInvalidRequestBody.Error()})
		return
	}

	form.SellerID = UserID
	h.logger.Info(form)
	ctx := req.Context()

	err = h.service.SellMarketBond(ctx, form)
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

	if err := h.response.JSON(w, http.StatusAccepted, domain.SuccessResponse{Message: "Success. Your bond is available to sell in the market."}); err != nil {
		h.logger.Error(err)
		_ = h.response.JSON(w, http.StatusInternalServerError, domain.ErrorResponse{ErrorMessage: httpErrors.InternalServerError.Error()})
		return
	}
}
