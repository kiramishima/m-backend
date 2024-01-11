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

var _ svcport.UserService = (*UserService)(nil)

// UserProfile struct
type UserService struct {
	logger         *zap.SugaredLogger
	repository     repport.UserRepository
	contextTimeOut time.Duration
}

// NewUserService creates a new user service
func NewUserService(logger *zap.SugaredLogger, repo repport.UserRepository, timeout time.Duration) *UserService {
	return &UserService{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

// GetProfile method.
func (svc *UserService) GetProfile(c context.Context, uid int) (*domain.UserProfile, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	up, err := svc.repository.GetProfile(ctx, uid)

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

	return up, nil
}

// UpdateProfile repository method.
func (svc *UserService) UpdateProfile(c context.Context, data *domain.UserProfileRequest) (*domain.UserProfile, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()

	// data
	var up, err = svc.repository.GetProfile(ctx, *data.UserID)
	if err != nil {
		switch {
		case errors.Is(err, httpErrors.ErrNoRecords):
			return nil, httpErrors.ErrNoRecords
		default:
			return nil, httpErrors.InternalServerError
		}
	}
	if data.UserName != nil {
		up.UserName = *data.UserName
	}
	if data.UserPhoto != nil {
		up.UserPhoto = *data.UserPhoto
	}
	if data.Gender != nil {
		up.Gender = *data.Gender
	}

	// Call repository
	_, err = svc.repository.UpdateProfile(ctx, up)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return nil, httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrAlreadyExists) {
				return nil, httpErrors.ErrAlreadyExists
			} else if errors.Is(err, httpErrors.ErrUserNotFound) {
				return nil, httpErrors.ErrUserNotFound
			} else {
				return nil, httpErrors.ErrBadEmailOrPassword
			}
		}
	}

	return up, nil
}

// GetBonds repository method.
func (svc *UserService) GetBonds(c context.Context, uid int) ([]*domain.Bond, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	list, err := svc.repository.GetBonds(ctx, uid)

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

	return list, nil
}
