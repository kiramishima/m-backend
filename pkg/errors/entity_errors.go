package errors

import "errors"

// Entity Errors
var (
	// User
	ErrEmptyUserField  = errors.New("username, password and email can't be empty")
	ErrFieldWithSpaces = errors.New("username and password can't have spaces")
	ErrShortPassword   = errors.New("password shorter than 6 characters")
	ErrLongPassword    = errors.New("password longer than 72 characters")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrUserNotFound    = errors.New("user not found")
)
