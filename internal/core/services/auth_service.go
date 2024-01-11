package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	repport "kiramishima/m-backend/internal/core/ports/repository"
	svcport "kiramishima/m-backend/internal/core/ports/services"
	httpErrors "kiramishima/m-backend/pkg/errors"
	"kiramishima/m-backend/pkg/utils"
	"time"
)

var _ svcport.AuthService = (*AuthService)(nil)

type AuthService struct {
	logger         *zap.SugaredLogger
	repository     repport.AuthRepository
	contextTimeOut time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(logger *zap.SugaredLogger, repo repport.AuthRepository, timeout time.Duration) *AuthService {
	return &AuthService{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

// FindByCredentials To Login users
func (svc *AuthService) FindByCredentials(c context.Context, data *domain.AuthRequest) (*domain.AuthResponse, error) {
	data.Password = data.Hash256Password(data.Password)
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	user, err := svc.repository.FindByCredentials(ctx, data)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return nil, httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrInvalidRequestBody) {
				return nil, httpErrors.BadQueryParams
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				return nil, httpErrors.ErrBadEmailOrPassword
			} else {
				return nil, httpErrors.InternalServerError
			}
		}
	}

	// Check Password
	if !data.ValidateBcryptPassword(user.Password, data.Password) {
		return nil, httpErrors.ErrBadPassword
	}

	// Generate Token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		svc.logger.Error(err.Error(), fmt.Sprintf("%T", err))
		return nil, jwt.ErrSignatureInvalid
	}

	return &domain.AuthResponse{Token: token}, nil
}

// Register repository method for create a new user.
func (svc *AuthService) Register(c context.Context, registerReq *domain.RegisterRequest) error {
	// Hash password
	// todo agregar el secret
	registerReq.Password = registerReq.Hash256Password(registerReq.Password)
	// Hash Bcrypt
	registerReq.Password, _ = registerReq.BcryptPassword(registerReq.Password)

	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	// Call repository
	err := svc.repository.Register(ctx, registerReq)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrAlreadyExists) {
				return httpErrors.ErrAlreadyExists
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				return httpErrors.ErrUserNotFound
			} else {
				return httpErrors.ErrBadEmailOrPassword
			}
		}
	}

	return nil
}
