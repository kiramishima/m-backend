package services

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	repport "kiramishima/m-backend/internal/core/ports/repository"
	svcport "kiramishima/m-backend/internal/core/ports/services"
	httpErrors "kiramishima/m-backend/pkg/errors"
	"time"
)

var _ svcport.MarketBondsService = (*MarketBondsService)(nil)

type MarketBondsService struct {
	logger         *zap.SugaredLogger
	repository     repport.MarketBondRepository
	contextTimeOut time.Duration
}

// NewMarketBondsService creates a new auth service
func NewMarketBondsService(logger *zap.SugaredLogger, repo repport.MarketBondRepository, timeout time.Duration) *MarketBondsService {
	return &MarketBondsService{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

// ListBonds return list of available bonds
func (svc *MarketBondsService) ListMarketBonds(c context.Context, uid int) ([]*domain.MarketBond, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.ListMarketBonds(ctx, uid)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return nil, httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				return nil, httpErrors.ErrNoRecords
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return nil, httpErrors.ErrExecuteStatement
			} else {
				return nil, httpErrors.InternalServerError
			}
		}
	}

	return data, nil
}

// GetMarketBondByID service method
func (svc *MarketBondsService) GetMarketBondByID(c context.Context, uid int, market_bond_id int) (*domain.MarketBond, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.GetMarketBondByID(ctx, market_bond_id)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return nil, httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				return nil, httpErrors.ErrNoRecords
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return nil, httpErrors.ErrExecuteStatement
			} else {
				return nil, httpErrors.InternalServerError
			}
		}
	}

	if data.CreatedByID == uid {
		data.IsOwner = true
	}

	return data, nil
}

// BuyBond repository method
func (svc *MarketBondsService) BuyMarketBond(c context.Context, order *domain.MarketBondRequest) error {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()

	data, err := svc.repository.GetMarketBondByID(ctx, *order.MarketBondID)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				return httpErrors.ErrNoRecords
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return httpErrors.ErrExecuteStatement
			} else {
				return httpErrors.InternalServerError
			}
		}
	}

	if data.CreatedByID == order.BuyerID {
		return httpErrors.ErrNoAvailableBonds
	}

	err = svc.repository.BuyMarketBond(ctx, order)
	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				return httpErrors.ErrNoRecords
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return httpErrors.ErrExecuteStatement
			} else {
				return httpErrors.InternalServerError
			}
		}
	}

	return nil
}

// SellBond repository method
func (svc *MarketBondsService) SellMarketBond(c context.Context, data *domain.MarketSellRequest) error {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()

	err := svc.repository.SellMarketBond(ctx, data)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrNoRecords) {
				return httpErrors.ErrNoRecords
			} else if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return httpErrors.ErrExecuteStatement
			} else {
				return httpErrors.InternalServerError
			}
		}
	}

	return nil
}
