package services

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	"time"

	"kiramishima/m-backend/internal/adapters/database/postgresql/repository"
)

// Module services
var Module = fx.Module("services",
	fx.Provide(func(cfg *domain.Configuration, logger *zap.SugaredLogger, authrepo *repository.AuthRepository) *AuthService {
		return NewAuthService(logger, authrepo, time.Duration(cfg.ContextTimeout)*time.Second)
	}),
	fx.Provide(func(cfg *domain.Configuration, logger *zap.SugaredLogger, bondrepo *repository.BondRepository) *BondService {
		return NewBondService(logger, bondrepo, time.Duration(cfg.ContextTimeout)*time.Second)
	}),
	fx.Provide(func(cfg *domain.Configuration, logger *zap.SugaredLogger, mbondrepo *repository.MarketBondRepository) *MarketBondsService {
		return NewMarketBondsService(logger, mbondrepo, time.Duration(cfg.ContextTimeout)*time.Second)
	}),
	fx.Provide(func(cfg *domain.Configuration, logger *zap.SugaredLogger, urepo *repository.UserRepository) *UserService {
		return NewUserService(logger, urepo, time.Duration(cfg.ContextTimeout)*time.Second)
	}),
)
