package domain

import (
	"encoding/hex"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,gte=6"`
	Name     string `json:"name" validate:"required,alphanum,gte=6"`
}

func (u *RegisterRequest) Hash256Password(password string) string {
	buf := []byte(password)
	pwd := sha3.New256()
	pwd.Write(buf)
	return hex.EncodeToString(pwd.Sum(nil))
}

func (u *RegisterRequest) BcryptPassword(password string) (string, error) {
	buf := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(buf, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func (u *RegisterRequest) ValidateBcryptPassword(password, password2 string) bool {
	byteHash := []byte(password)
	buf := []byte(password2)
	err := bcrypt.CompareHashAndPassword(byteHash, buf)
	if err != nil {
		return false
	}
	return true
}

func (u *RegisterRequest) Validate(v *validator.Validate) error {
	err := v.Struct(u)
	if err != nil {
		errormsg := ""
		for _, err := range err.(validator.ValidationErrors) {
			errormsg = fmt.Sprintf("Field: %s, Error: %s", err.Field(), err.Tag())
		}

		// from here you can create your own error messages in whatever language you wish
		return fmt.Errorf(errormsg)
	}
	return nil
}
