package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"kiramishima/m-backend/internal/core/domain"
	dbErrors "kiramishima/m-backend/pkg/errors"
	"testing"
	"time"
)

func TestFindByCredentials(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()
	repo := NewAuthRepository(sqlxDB)

	user := &domain.User{
		ID:        "1",
		Email:     "gini@mail.com",
		Password:  "12356",
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	form := &domain.AuthRequest{
		Email:    "gini@mail.com",
		Password: "12356",
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?").
			ExpectQuery().
			WithArgs(form.Email).
			WillReturnRows(rows)

		userDB, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: form.Email, Password: form.Password})
		assert.NoError(t, err)
		assert.Equal(t, user, userDB)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?").
			ExpectQuery().
			WithArgs(form.Email).
			WillReturnError(sql.ErrConnDone)

		userProfile, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: form.Email, Password: form.Password})
		assert.Error(t, err)
		assert.Empty(t, userProfile)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?").
			WillReturnError(sql.ErrConnDone)

		userMock, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: form.Email, Password: form.Password})
		assert.Error(t, err)
		assert.Empty(t, userMock)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?").
			ExpectQuery().
			WithArgs(form.Email).
			WillReturnError(sql.ErrNoRows)

		userProfile, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: form.Email, Password: form.Password})
		assert.Error(t, err)
		assert.Empty(t, userProfile)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRegister(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()
	repo := NewAuthRepository(sqlxDB)

	var users = []*domain.User{
		{
			ID:        "1",
			Email:     "gini@mail.com",
			Password:  "12356",
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
		{
			ID:        "2",
			Email:     "gin@mail.com",
			Password:  "12356",
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
		},
	}

	form := &domain.AuthRequest{
		Email:    "gini2@mail.com",
		Password: "12356",
	}

	form2 := &domain.AuthRequest{
		Email:    "gini@mail.com",
		Password: "12356",
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(users[0].ID, users[0].Email, users[0].Password, users[0].CreatedAt, users[0].UpdatedAt).
			AddRow(users[1].ID, users[1].Email, users[1].Password, users[1].CreatedAt, users[1].UpdatedAt)

		mock.ExpectPrepare("INSERT INTO users(email, password) VALUES(?, ?)").
			ExpectQuery().
			WithArgs(form.Email, form.Password).
			WillReturnRows(rows)

		err := repo.Register(ctx, &domain.RegisterRequest{Email: form.Email, Password: form.Password})
		assert.NoError(t, err)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Duplicate record", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users(email, password) VALUES(?, ?)").
			ExpectQuery().
			WithArgs(form2.Email, form2.Password).
			WillReturnError(dbErrors.ErrAlreadyExists)

		err := repo.Register(ctx, &domain.RegisterRequest{Email: form2.Email, Password: form2.Password})
		// t.Log(err)
		assert.ErrorIs(t, err, dbErrors.ErrAlreadyExists)
		// assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users(email, password) VALUES(?, ?)").
			WillReturnError(sql.ErrConnDone)

		err := repo.Register(ctx, &domain.RegisterRequest{Email: form2.Email, Password: form2.Password})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users(email, password) VALUES(?, ?)").
			WillReturnError(dbErrors.ErrPrepareStatement)

		err := repo.Register(ctx, &domain.RegisterRequest{Email: form.Email, Password: form.Password})
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
