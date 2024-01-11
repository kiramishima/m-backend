package repository

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

type BondRepository interface {
	ListBonds(ctx context.Context, uid int) ([]*domain.Bond, error)
	GetBondByID(ctx context.Context, bond_id int) (*domain.Bond, error)
	CreateBond(ctx context.Context, data *domain.BondRequest) error
	UpdateBond(ctx context.Context, udata *domain.Bond) error
	DeleteBond(ctx context.Context, bond_id int) error
}
