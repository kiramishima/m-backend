package services

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

type UserService interface {
	GetProfile(c context.Context, uid int) (*domain.UserProfile, error)
	UpdateProfile(c context.Context, data *domain.UserProfileRequest) (*domain.UserProfile, error)
	GetBonds(c context.Context, uid int) ([]*domain.Bond, error)
}
