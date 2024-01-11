package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	cache "kiramishima/m-backend/internal/adapters/cache/redis"
	"kiramishima/m-backend/internal/core/domain"
	rPort "kiramishima/m-backend/internal/core/ports/repository"
	dbErrors "kiramishima/m-backend/pkg/errors"
)

var _ rPort.UserRepository = (*UserRepository)(nil)

// UserRepository struct
type UserRepository struct {
	db    *sqlx.DB
	cache *cache.RedisCache
}

// NewUserRepository Creates a new instance of BondRepository
func NewUserRepository(conn *sqlx.DB, cache *cache.RedisCache) *UserRepository {
	return &UserRepository{
		db:    conn,
		cache: cache,
	}
}

// GetProfile repository method.
func (repo *UserRepository) GetProfile(ctx context.Context, uid int) (*domain.UserProfile, error) {
	var query = `SELECT
		username,
		photo,
		gender
	FROM users_profile
	WHERE user_id = ?`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrExecuteStatement
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, uid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbErrors.ErrNoRecords
		} else {
			return nil, dbErrors.ErrExecuteStatement
		}
	}
	var item = &domain.UserProfile{}
	err = row.Scan(item.UserName, item.UserPhoto, item.Gender)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbErrors.ErrScanData, err)
	}
	return item, nil
}

// Register repository method for create a new user.
func (repo *UserRepository) UpdateProfile(ctx context.Context, data *domain.UserProfile) (*domain.UserProfile, error) {

	var query = `UPDATE users_profile SET username = ?, photo = ?, gender = ?, updated_at = NOW() WHERE user_id = ?`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbErrors.ErrPrepareStatement, err)
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, dbErrors.ErrBeginTransaction
	}
	// Exec
	_, err = stmt.ExecContext(ctx, data.UserName, data.UserPhoto, data.Gender, data.UserID)

	if err != nil {
		switch {
		case errors.Is(err, dbErrors.ErrNoRecords):
			return nil, dbErrors.ErrNoRecords
		default:
			return nil, dbErrors.InternalServerError
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, dbErrors.ErrCommit
	}

	return nil, nil
}

func (repo *UserRepository) GetBonds(ctx context.Context, uid int) ([]*domain.Bond, error) {
	var query = `SELECT
    		b.uuid,
    		b.name,
    		b.price,
    		c.currency,
    		up.username AS created_by,
    		b.created_by AS created_by_id,
    		b.status,
    		b.created_at,
    		b.updated_at
    	FROM bonds b
			INNER JOIN currencies c on c.id = b.currency_id
			INNER JOIN users_profile up on b.created_by = up.user_id
		WHERE created_by = ?`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbErrors.ErrPrepareStatement, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, uid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbErrors.ErrNoRecords
		} else {
			return nil, dbErrors.ErrExecuteStatement
		}
	}
	var list = make([]*domain.Bond, 0)
	for rows.Next() {
		var createAt sql.NullTime
		var updatedAt sql.NullTime
		var item = &domain.Bond{}
		err = rows.Scan(item.UUID, item.Name, item.Price, item.Currency, item.CreatedBy, item.CreatedByID, item.Status, &createAt, &updatedAt)
		if err != nil {
			break
		}
		if createAt.Valid {
			item.CreatedAt = createAt.Time
		}
		if updatedAt.Valid {
			item.UpdateAt = updatedAt.Time
		}
		if item.CreatedByID == uid {
			item.IsOwner = true
		}
		list = append(list, item)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbErrors.ErrScanData, err)
	}
	return list, nil
}
