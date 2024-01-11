package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/services"
)

// Module Handlers.
var Module = fx.Module("handlers",
	fx.Invoke(func(r *chi.Mux, logger *zap.SugaredLogger, svc *services.AuthService, render *render.Render, validate *validator.Validate) {
		NewAuthHandlers(r, logger, svc, render, validate)
	}),
	fx.Invoke(func(r *chi.Mux, logger *zap.SugaredLogger, svc *services.BondService, render *render.Render, validate *validator.Validate) {
		NewBondHandlers(r, logger, svc, render, validate)
	}),
	fx.Invoke(func(r *chi.Mux, logger *zap.SugaredLogger, svc *services.UserService, render *render.Render, validate *validator.Validate) {
		NewUserHandlers(r, logger, svc, render, validate)
	}),
	fx.Invoke(func(r *chi.Mux, logger *zap.SugaredLogger, svc *services.MarketBondsService, render *render.Render, validate *validator.Validate) {
		NewMarketBondsHandlers(r, logger, svc, render, validate)
	}),
)
