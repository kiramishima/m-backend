package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	cache "kiramishima/m-backend/internal/adapters/cache/redis"
	"kiramishima/m-backend/internal/core/domain"
	rPort "kiramishima/m-backend/internal/core/ports/repository"
	dbErrors "kiramishima/m-backend/pkg/errors"
)

var _ rPort.MarketBondRepository = (*MarketBondRepository)(nil)

// MarketBondRepository struct
type MarketBondRepository struct {
	db    *sqlx.DB
	cache *cache.RedisCache
}

// NewMarketBondRepository Creates a new instance of BondRepository
func NewMarketBondRepository(conn *sqlx.DB, cache *cache.RedisCache) *MarketBondRepository {
	return &MarketBondRepository{
		db:    conn,
		cache: cache,
	}
}

// ListMarketBonds repository method for listing the bonds.
func (repo *MarketBondRepository) ListMarketBonds(ctx context.Context, uid int) ([]*domain.MarketBond, error) {
	var query = `SELECT
    		mb.id,
    		b.uuid,
    		b.name,
    		b.price,
    		mb.available,
    		c.currency,
    		up.username AS created_by,
    		b.created_by AS created_by_id,
    		b.status,
    		b.created_at,
    		b.updated_at
    	FROM market_bonds mb
			INNER JOIN bonds b on b.id = mb.bonds_id
			INNER JOIN currencies c on c.id = b.currency_id
			INNER JOIN users_profile up on b.created_by = up.user_id
		WHERE b.status = 'on_sell'`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrPrepareStatement
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbErrors.ErrNoRecords
		} else {
			return nil, dbErrors.ErrExecuteStatement
		}
	}
	var list = make([]*domain.MarketBond, 0)
	for rows.Next() {
		var createAt sql.NullTime
		var updatedAt sql.NullTime
		var item = &domain.MarketBond{}
		err = rows.Scan(&item.ID, &item.UUID, &item.Name, &item.Price, &item.Available, &item.Currency, &item.CreatedBy, &item.CreatedByID, &item.Status, &createAt, &updatedAt)
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
		return nil, dbErrors.ErrScanData
	}
	return list, nil
}

// GetMarketBondByUUID repository method for listing the bonds.
func (repo *MarketBondRepository) GetMarketBondByID(ctx context.Context, market_bond_id int) (*domain.MarketBond, error) {
	var query = `SELECT
    		mb.id,
    		b.uuid,
    		b.name,
    		b.price,
    		mb.available,
    		c.currency,
    		up.username AS created_by,
    		b.created_by AS created_by_id,
    		b.status,
    		b.created_at,
    		b.updated_at
    	FROM market_bonds mb
			INNER JOIN bonds b on b.id = mb.bonds_id
			INNER JOIN currencies c on c.id = b.currency_id
			INNER JOIN users_profile up on b.created_by = up.user_id
		WHERE b.status = 'on_sell' AND mb.status = 'available' AND mb.deleted_at IS NOT NULL AND b.deleted_at IS NOT NULL AND b.uuid = ?`

	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrPrepareStatement
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, market_bond_id)

	var createAt sql.NullTime
	var updatedAt sql.NullTime
	var item = &domain.MarketBond{}
	err = row.Scan(&item.ID, &item.UUID, &item.Name, &item.Price, &item.Available, &item.Currency, &item.CreatedBy, &item.CreatedByID, &item.Status, &createAt, &updatedAt)
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

// BuyMarketBond repository method
func (repo *MarketBondRepository) BuyMarketBond(ctx context.Context, order *domain.MarketBondRequest) error {
	var mbond = struct {
		BondID    int `db:"bond_id"`
		Available int `db:"available"`
	}{
		BondID:    0,
		Available: 0,
	}
	var query = `SELECT id, bond_id, available FROM market_bonds WHERE id = ? LIMIT 1`
	err := repo.db.SelectContext(ctx, &mbond, query, order.MarketBondID)
	if err != nil {
		return dbErrors.ErrExecuteQuery
	}

	if mbond.Available < *order.Order {
		return dbErrors.ErrNoAvailableBonds
	}

	// Init TX
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return dbErrors.ErrBeginTransaction
	}

	op := mbond.Available - *order.Order
	status := "available"
	if op < 0 {
		tx.Rollback()
		return dbErrors.ErrNoAvailableBonds
	} else if op == 0 {
		status = "bought"
	}

	query = `INSERT INTO transactions (seller_id, buyer_id, bond_id, total_acquired, status)
		VALUES(?, ?, ?, ?, ?)`
	result, err := repo.db.ExecContext(ctx, query, order.SellerID, order.BuyerID, mbond.BondID, order.Order, 0)
	LastInsID, _ := result.LastInsertId()

	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, dbErrors.ErrNoRecords):
			return dbErrors.ErrDeleteBond
		default:
			return dbErrors.ErrDeleteBond
		}
	}
	// Update Market Bonds

	query = `UPDATE market_bonds SET available = ?, status = ? WHERE id = ?`
	_, err = repo.db.ExecContext(ctx, query, op, status, order.MarketBondID)
	if err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, dbErrors.ErrNoRecords):
			return dbErrors.ErrDeleteBond
		default:
			return dbErrors.ErrDeleteBond
		}
	}

	// Update Transaction Status
	query = `UPDATE transactions SET status = ? WHERE id = ?`
	_, err = repo.db.ExecContext(ctx, query, 1, LastInsID)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return dbErrors.ErrCommit
	}

	return nil
}

func (repo *MarketBondRepository) SellMarketBond(ctx context.Context, data *domain.MarketSellRequest) error {
	var available int
	var query = `SELECT number FROM bonds WHERE id = ? AND created_by = ?`
	err := repo.db.SelectContext(ctx, &available, query, data.BondID, data.SellerID)
	if err != nil {
		return dbErrors.ErrExecuteQuery
	}

	// Init TX
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		return dbErrors.ErrBeginTransaction
	}

	// Num Bonds to sell
	op := 0
	if (available - *data.Num) == 0 {
		op = available
	} else if (available - *data.Num) > 0 {
		op = *data.Num
	} else {
		tx.Rollback()
		return dbErrors.ErrNoAvailableBonds
	}

	query = `INSERT INTO market_bonds (bond_id, available)
		VALUES(?, ?)`
	_, err = repo.db.ExecContext(ctx, query, data.BondID, op)

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
