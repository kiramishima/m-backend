package services

import (
	"context"
	"kiramishima/m-backend/internal/core/domain"
)

type AuthService interface {
	FindByCredentials(ctx context.Context, data *domain.AuthRequest) (*domain.AuthResponse, error)
	Register(ctx context.Context, registerReq *domain.RegisterRequest) error
}
