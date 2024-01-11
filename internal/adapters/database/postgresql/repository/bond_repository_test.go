package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"kiramishima/m-backend/internal/core/domain"
	dbErrors "kiramishima/m-backend/pkg/errors"
	"testing"
	"time"
)

func TestListBonds(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()
	ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
	defer cancel()

	repo := NewBondRepository(sqlxDB, nil)

	var uuid1 = uuid.NewString()
	var uuid2 = uuid.NewString()
	var bonds = []*domain.Bond{
		{
			UUID:        uuid1,
			Name:        faker.Name(),
			Price:       10000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 1,
			Status:      "available",
			IsOwner:     false,
			CreatedAt:   time.Now(),
		},
		{
			UUID:        uuid2,
			Name:        faker.Name(),
			Price:       12000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 2,
			Status:      "available",
			IsOwner:     true,
			CreatedAt:   time.Now(),
		},
	}

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
			INNER JOIN users_profile up on b.created_by = up.user_id`

	rows := sqlmock.NewRows([]string{"uuid", "name", "price", "currency", "created_by", "created_by_id", "status", "created_at", "updated_at"}).
		AddRow(&bonds[0].UUID, bonds[0].Name, bonds[0].Price, bonds[0].Currency, bonds[0].CreatedBy, bonds[0].CreatedByID, bonds[0].Status, bonds[0].CreatedAt, bonds[0].UpdateAt).
		AddRow(&bonds[1].UUID, bonds[1].Name, bonds[1].Price, bonds[1].Currency, bonds[1].CreatedBy, bonds[1].CreatedByID, bonds[1].Status, bonds[1].CreatedAt, bonds[1].UpdateAt)

	t.Run("OK", func(t *testing.T) {

		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnRows(rows)

		list, err := repo.ListBonds(ctx, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, list)
		assert.Equal(t, len(list), 2)
		assert.Equal(t, list[0].UUID, uuid1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare(query).
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		list, err := repo.ListBonds(ctx, 1)
		t.Log("err", err, list)
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare(query).
			WillReturnError(dbErrors.ErrPrepareStatement)

		list, err := repo.ListBonds(ctx, 1)
		t.Log("err", err)
		assert.Error(t, err)
		assert.Nil(t, list)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateBond(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()
	ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
	defer cancel()

	repo := NewBondRepository(sqlxDB, nil)

	var bonds = []*domain.Bond{
		{
			UUID:        uuid.NewString(),
			Name:        faker.Name(),
			Price:       10000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 1,
			Status:      "on_hold",
			IsOwner:     false,
			CreatedAt:   time.Now(),
		},
		{
			UUID:        uuid.NewString(),
			Name:        faker.Name(),
			Price:       12000.0000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 2,
			Status:      "on_sell",
			IsOwner:     true,
			CreatedAt:   time.Now(),
		},
	}

	var n1 = faker.Name()
	var p1 float32 = 12000.00001
	var num1 = 5000
	var status1 = "on_hold"
	var bondNotExisting = &domain.BondRequest{
		Name:      &n1,
		Price:     &p1,
		Number:    &num1,
		CreatedBy: 1,
		Status:    &status1,
	}

	var bondExisting = &domain.BondRequest{
		Name:      &bonds[0].Name,
		Price:     &bonds[0].Price,
		Number:    &bonds[0].Number,
		CreatedBy: 1,
		Status:    &bonds[0].Status,
	}

	var query = `INSERT INTO bonds (uuid, name, number, price, currency_id, created_by, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	t.Run("OK", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(query).
			WithArgs(sqlmock.AnyArg(), bondNotExisting.Name, bondNotExisting.Number, bondNotExisting.Price, bondNotExisting.CurrencyID, bondNotExisting.CreatedBy, bondNotExisting.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.CreateBond(ctx, bondNotExisting)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(query).
			WithArgs(sqlmock.AnyArg(), bondExisting.Name, bondExisting.Number, bondExisting.Price, bondExisting.CurrencyID, bondExisting.CreatedBy, bondExisting.Status).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(dbErrors.ErrBondAlreadyExists)

		mock.ExpectCommit()

		err := repo.CreateBond(ctx, bondExisting)
		t.Log("err", err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErrors.ErrBondAlreadyExists)
		assert.Error(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(query).
			WithArgs(sqlmock.AnyArg(), bondExisting.Name, bondExisting.Number, bondExisting.Price, bondExisting.CurrencyID, bondExisting.CreatedBy, bondExisting.Status).
			WillReturnError(dbErrors.ErrExecuteQuery)
		mock.ExpectCommit()

		err := repo.CreateBond(ctx, bondExisting)
		t.Log("err", err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErrors.ErrExecuteQuery)
		assert.Error(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateBond(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	c := context.Background()
	ctx, cancel := context.WithTimeout(c, time.Duration(5)*time.Second)
	defer cancel()

	repo := NewBondRepository(sqlxDB, nil)

	var bonds = []*domain.Bond{
		{
			UUID:        uuid.NewString(),
			Name:        faker.Name(),
			Price:       10000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 1,
			Status:      "available",
			IsOwner:     false,
			CreatedAt:   time.Now(),
		},
		{
			UUID:        uuid.NewString(),
			Name:        faker.Name(),
			Price:       12000,
			Currency:    1,
			CreatedBy:   faker.Username(),
			CreatedByID: 2,
			Status:      "available",
			IsOwner:     true,
			CreatedAt:   time.Now(),
		},
	}

	var bond = bonds[0]
	bond.Name = faker.Name() + " " + faker.LastName()
	bond.Number = 1000
	bond.Status = "bought"

	var query = `UPDATE bonds SET name = ?, number = ?, price = ?, currency_id = ?, status = ? 
		WHERE uuid = ? AND created_by = ?`

	t.Run("OK", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(query).
			WithArgs(bond.Name, bond.Number, bond.Price, bond.Currency, bond.CreatedByID, bond.Status, bond.UUID, bond.CreatedByID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := repo.UpdateBond(ctx, bond)
		t.Log("err", err)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(query).
			WithArgs(bond.Name, bond.Number, bond.Price, bond.Currency, bond.CreatedBy, bond.Status, bond.UUID, bond.CreatedByID).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(dbErrors.ErrUpdatingRecord)

		mock.ExpectCommit()

		err := repo.UpdateBond(ctx, bond)
		t.Log("err", err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErrors.ErrUpdatingRecord)
		assert.Error(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(query).
			WithArgs(bond.Name, bond.Number, bond.Price, bond.Currency, bond.CreatedBy, bond.Status, bond.UUID, bond.CreatedByID).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(dbErrors.InternalServerError)
		mock.ExpectCommit()

		err := repo.UpdateBond(ctx, bond)
		t.Log("err", err)
		assert.Error(t, err)
		assert.ErrorIs(t, err, dbErrors.InternalServerError)
		assert.Error(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteBond(t *testing.T) {}
