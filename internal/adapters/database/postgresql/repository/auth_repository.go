package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	my "github.com/go-mysql/errors"
	"github.com/jmoiron/sqlx"
	"kiramishima/m-backend/internal/core/domain"
	rPort "kiramishima/m-backend/internal/core/ports/repository"
	dbErrors "kiramishima/m-backend/pkg/errors"
	"log"
)

var _ rPort.AuthRepository = (*AuthRepository)(nil)

// AuthRepository struct
type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository Creates a new instance of AuthRepository
func NewAuthRepository(conn *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: conn,
	}
}

// FindByCredentials Repository method for sign in
func (repo *AuthRepository) FindByCredentials(ctx context.Context, data *domain.AuthRequest) (*domain.User, error) {
	var query = `SELECT id,
		   email,
		   password,
		   created_at,
		   updated_at
	FROM users
	WHERE email = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, dbErrors.ErrPrepareStatement
	}
	defer stmt.Close()

	u := &domain.User{}

	row := stmt.QueryRowContext(ctx, data.Email)
	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	err = row.Scan(&u.ID, &u.Email, &u.Password, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dbErrors.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("%s: %w", dbErrors.ErrScanData, err)
		}
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return u, nil
}

// Register repository method for create a new user.
func (repo *AuthRepository) Register(ctx context.Context, registerReq *domain.RegisterRequest) error {
	var query = `INSERT INTO users(email, password) VALUES(?, ?)`
	stmt, err := repo.db.PreparexContext(ctx, query)
	if err != nil {
		return dbErrors.ErrPrepareStatement
	}
	defer stmt.Close()

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return dbErrors.ErrBeginTransaction
	}

	res, err := stmt.ExecContext(ctx, registerReq.Email, registerReq.Password)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		log.Println("Code ", my.MySQLErrorCode(err))
		// log.Println("Code 2 ", errors.Is(err, my.ErrDupeKey))
		if err.Error() == "Error 1062 (23000): Duplicate entry 'mail3@mail.com' for key 'email'" {
			return dbErrors.ErrAlreadyExists
		}
		if errors.Is(err, my.ErrDupeKey) {
			return dbErrors.ErrAlreadyExists
		} else if errors.Is(err, sql.ErrConnDone) {
			return sql.ErrConnDone
		} else {
			return dbErrors.ErrScanData
		}
	}

	LastInsID, _ := res.LastInsertId()
	query = `INSERT INTO users_profile(user_id, username) VALUES(?, ?)`
	_, err = repo.db.ExecContext(ctx, query, LastInsID, registerReq.Name)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		if errors.Is(err, my.ErrDupeKey) {
			return dbErrors.ErrAlreadyExists
		} else if errors.Is(err, sql.ErrConnDone) {
			return sql.ErrConnDone
		} else {
			return dbErrors.ErrScanData
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return dbErrors.ErrCommit
	}

	return nil
}
