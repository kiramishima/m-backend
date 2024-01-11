package repository

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

// MarketBondRepository interface
type MarketBondRepository interface {
	ListMarketBonds(ctx context.Context, uid int) ([]*domain.MarketBond, error)
	GetMarketBondByID(ctx context.Context, market_bond_id int) (*domain.MarketBond, error)
	BuyMarketBond(ctx context.Context, order *domain.MarketBondRequest) error
	SellMarketBond(ctx context.Context, data *domain.MarketSellRequest) error
}
