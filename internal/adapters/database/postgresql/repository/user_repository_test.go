package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"kiramishima/m-backend/internal/core/domain"
	"testing"
	"time"
)

func TestGetProfile(t *testing.T) {
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

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ").
			ExpectQuery().
			WithArgs(user.Email).
			WillReturnRows(rows)

		userDB, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: user.Email, Password: user.Password})
		assert.NoError(t, err)
		assert.Equal(t, user, userDB)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ").
			ExpectQuery().
			WithArgs(user.Email).
			WillReturnError(sql.ErrConnDone)

		userProfile, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: user.Email, Password: user.Password})
		assert.Error(t, err)
		assert.Empty(t, userProfile)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ").
			WillReturnError(sql.ErrConnDone)

		userMock, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: user.Email, Password: user.Password})
		assert.Error(t, err)
		assert.Empty(t, userMock)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT id, email, password, created_at, updated_at FROM users WHERE email = ").
			ExpectQuery().
			WithArgs(user.Email).
			WillReturnError(sql.ErrNoRows)

		userProfile, err := repo.FindByCredentials(ctx, &domain.AuthRequest{Email: user.Email, Password: user.Password})
		assert.Error(t, err)
		assert.Empty(t, userProfile)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateProfile(t *testing.T) {}

func TestGetBonds(t *testing.T) {}
