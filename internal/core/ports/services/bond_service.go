package services

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

// BondService interface
type BondService interface {
	ListBonds(c context.Context, uid int) ([]*domain.Bond, error)
	GetBondByID(ctx context.Context, uid int, bond_id int) (*domain.Bond, error)
	CreateBond(ctx context.Context, data *domain.BondRequest) error
	UpdateBond(ctx context.Context, bond_id int, udata *domain.BondRequest) error
	DeleteBond(ctx context.Context, bond_id int) error
}
