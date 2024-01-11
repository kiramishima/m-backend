package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	my "github.com/go-mysql/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	cache "kiramishima/m-backend/internal/adapters/cache/redis"
	"kiramishima/m-backend/internal/core/domain"
	rPort "kiramishima/m-backend/internal/core/ports/repository"
	dbErrors "kiramishima/m-backend/pkg/errors"
	"log"
)

var _ rPort.BondRepository = (*BondRepository)(nil)

// BondRepository struct
type BondRepository struct {
	db    *sqlx.DB
	cache *cache.RedisCache
}

// NewBondRepository Creates a new instance of BondRepository
func NewBondRepository(conn *sqlx.DB, cache *cache.RedisCache) *BondRepository {
	return &BondRepository{
		db:    conn,
		cache: cache,
	}
}

// ListBonds repository method for listing the bonds.
func (repo *BondRepository) ListBonds(ctx context.Context, uid int) ([]*domain.Bond, error) {
	var query = `SELECT
    		b.id,
    		b.uuid,
    		b.name,
    		b.price,
    		b.number,
    		c.currency,
    		up.username AS created_by,
    		b.created_by AS created_by_id,
    		b.status,
    		(SELECT COUNT(*) FROM market_bonds WHERE bond_id = b.id) on_sale,
    		b.created_at,
    		b.updated_at
    	FROM bonds b
			INNER JOIN currencies c on c.id = b.currency_id
			INNER JOIN users_profile up on b.created_by = up.user_id
		WHERE b.deleted_at IS NOT NULL AND b.created_by = ?`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrPrepareStatement
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
		err = rows.Scan(&item.ID, &item.UUID, &item.Name, &item.Price, &item.Number, &item.Currency, &item.CreatedBy, &item.CreatedByID, &item.Status, &createAt, &updatedAt)
		if err != nil {
			break
		}
		if createAt.Valid {
			item.CreatedAt = createAt.Time
		}
		if updatedAt.Valid {
			item.UpdateAt = updatedAt.Time
		}
		item.IsOwner = true

		list = append(list, item)
	}
	if err != nil {
		return nil, dbErrors.ErrScanData
	}
	return list, nil
}

// GetBondById repository method for listing the bonds.
func (repo *BondRepository) GetBondByID(ctx context.Context, bond_id int) (*domain.Bond, error) {
	var query = `SELECT
    		b.id,
    		b.uuid,
    		b.name,
    		b.price,
    		b.number,
    		c.currency,
    		up.username AS created_by,
    		b.created_by AS created_by_id,
    		b.status,
    		(SELECT COUNT(*) FROM market_bonds WHERE bond_id = b.id) on_sale,
    		b.created_at,
    		b.updated_at
    	FROM bonds b
			INNER JOIN currencies c on c.id = b.currency_id
			INNER JOIN users_profile up on b.created_by = up.user_id
		WHERE b.deleted_at IS NOT NULL AND b.id = ?`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrPrepareStatement
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, bond_id)

	var createAt sql.NullTime
	var updatedAt sql.NullTime
	var item = &domain.Bond{}
	err = row.Scan(&item.ID, &item.UUID, &item.Name, &item.Price, &item.Number, &item.Currency, &item.CreatedBy, &item.CreatedByID, &item.Status, &createAt, &updatedAt)
	if err != nil {
		return nil, dbErrors.ErrScanData
	}
	if createAt.Valid {
		item.CreatedAt = createAt.Time
	}
	if updatedAt.Valid {
		item.UpdateAt = updatedAt.Time
	}

	if err != nil {
		return nil, dbErrors.ErrScanData
	}
	return item, nil
}

// CreateBond repository method for create a new bond.
func (repo *BondRepository) CreateBond(ctx context.Context, data *domain.BondRequest) error {
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		// tx.Rollback()
		return dbErrors.ErrBeginTransaction
	}
	// query
	var query = `INSERT INTO bonds (uuid, name, number, price, currency_id, created_by, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	// uuid
	var uid = uuid.NewString()
	_, err = repo.db.ExecContext(ctx, query, uid, data.Name, data.Number, data.Price, data.CurrencyID, data.CreatedBy, data.Status)

	if err != nil {
		if ok, myerr := my.Error(err); ok {
			if errors.Is(myerr, my.ErrDupeKey) {
				return dbErrors.ErrBondAlreadyExists
			} else {
				log.Println("[ERROR SQL] ", myerr.Error())
				return dbErrors.ErrUserNotFound
			}
		} else {
			return fmt.Errorf("%s: %w", dbErrors.ErrScanData, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return my.ErrQueryKilled
	}

	return nil
}

// Register repository method for create a new user.
func (repo *BondRepository) UpdateBond(ctx context.Context, udata *domain.Bond) error {
	var query = `UPDATE bonds SET name = ?, number = ?, price = ?, currency_id = ?, status = ? 
             WHERE uuid = ? AND created_by = ?`

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return dbErrors.ErrBeginTransaction
	}
	// uuid
	_, err = repo.db.ExecContext(ctx, query, udata.Name, udata.Number, udata.Price, udata.Currency, udata.CreatedBy, udata.Status, udata.UUID, udata.CreatedByID)

	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, dbErrors.ErrNoRecords):
			return dbErrors.ErrUpdatingRecord
		default:
			return dbErrors.InternalServerError
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return dbErrors.ErrCommit
	}

	return nil
}

// DeleteBond repository method
func (repo *BondRepository) DeleteBond(ctx context.Context, bond_id int) error {
	var query = `UPDATE bonds SET deleted_at=NOW() WHERE id = ?`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", dbErrors.ErrPrepareStatement, err)
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return dbErrors.ErrBeginTransaction
	}

	_, err = stmt.ExecContext(ctx, bond_id)

	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, dbErrors.ErrNoRecords):
			return dbErrors.ErrDeleteBond
		default:
			return dbErrors.ErrDeleteBond
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return dbErrors.ErrCommit
	}

	return nil
}
