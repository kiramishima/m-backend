package domain

import (
	appErr "kiramishima/m-backend/pkg/errors"
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// NewUser crea un nuevo usuario
func NewUser(username, password, email string) (*User, error) {
	user := &User{
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate valida al usuario
func (user *User) Validate() error {
	if user.Password == "" || user.Email == "" {
		return appErr.ErrEmptyUserField
	}

	if strings.ContainsAny(user.Password, " \t\r\n") {
		return appErr.ErrFieldWithSpaces
	}

	if len(user.Password) < 6 {
		return appErr.ErrShortPassword
	}

	if len(user.Password) > 72 {
		return appErr.ErrLongPassword
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return appErr.ErrInvalidEmail
	}

	return nil
}
