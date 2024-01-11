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

var _ svcport.BondService = (*BondService)(nil)

type BondService struct {
	logger         *zap.SugaredLogger
	repository     repport.BondRepository
	contextTimeOut time.Duration
}

// NewAuthService creates a new auth service
func NewBondService(logger *zap.SugaredLogger, repo repport.BondRepository, timeout time.Duration) *BondService {
	return &BondService{
		logger:         logger,
		repository:     repo,
		contextTimeOut: timeout,
	}
}

// ListBonds return list of available bonds
func (svc *BondService) ListBonds(c context.Context, uid int) ([]*domain.Bond, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.ListBonds(ctx, uid)

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

// GetBondById service method
func (svc *BondService) GetBondByID(c context.Context, uid int, bond_id int) (*domain.Bond, error) {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	data, err := svc.repository.GetBondByID(ctx, bond_id)

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

// CreateBond repository method for create a new user.
func (svc *BondService) CreateBond(c context.Context, data *domain.BondRequest) error {
	// add

	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	// Call repository
	err := svc.repository.CreateBond(ctx, data)

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
				return httpErrors.ErrExecuteStatement
			}
		}
	}

	return nil
}

// UpdateBond repository method for update a bond.
func (svc *BondService) UpdateBond(c context.Context, bond_id int, udata *domain.BondRequest) error {
	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	// data
	var data, err = svc.repository.GetBondByID(ctx, bond_id)
	if err != nil {
		switch {
		case errors.Is(err, httpErrors.ErrNoRecords):
			return httpErrors.ErrNoRecords
		default:
			return httpErrors.InternalServerError
		}
	}
	if udata.Name != nil {
		data.Name = *udata.Name
	}
	if udata.Number != nil {
		data.Number = *udata.Number
	}
	if udata.Price != nil {
		data.Price = *udata.Price
	}
	// Call repository
	err = svc.repository.UpdateBond(ctx, data)

	if err != nil {
		svc.logger.Error(err.Error())

		select {
		case <-ctx.Done():
			return httpErrors.ErrTimeout
		default:
			if errors.Is(err, httpErrors.ErrExecuteStatement) {
				return httpErrors.ErrExecuteStatement
			} else if errors.Is(err, httpErrors.ErrItemNotFound) {
				return httpErrors.ErrItemNotFound
			} else {
				return httpErrors.ErrUpdatingRecord
			}
		}
	}

	return nil
}

// DeleteBond repository method for delete logical a bond.
func (svc *BondService) DeleteBond(c context.Context, bond_id int) error {

	// context
	ctx, cancel := context.WithTimeout(c, svc.contextTimeOut)
	defer cancel()
	// Call repository
	err := svc.repository.DeleteBond(ctx, bond_id)

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
