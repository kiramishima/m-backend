package services

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

// MarketBondsService interface
type MarketBondsService interface {
	ListMarketBonds(ctx context.Context, uid int) ([]*domain.MarketBond, error)
	GetMarketBondByID(ctx context.Context, uid int, market_bond_id int) (*domain.MarketBond, error)
	BuyMarketBond(ctx context.Context, order *domain.MarketBondRequest) error
	SellMarketBond(ctx context.Context, data *domain.MarketSellRequest) error
}
