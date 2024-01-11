package repository

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

type UserRepository interface {
	GetProfile(ctx context.Context, uid int) (*domain.UserProfile, error)
	UpdateProfile(ctx context.Context, data *domain.UserProfile) (*domain.UserProfile, error)
	GetBonds(ctx context.Context, uid int) ([]*domain.Bond, error)
}
