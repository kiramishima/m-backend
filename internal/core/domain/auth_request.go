package domain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AuthRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
	Token    string `json:"token,omitempty"`
}

func (u *AuthRequest) Hash256Password(password string) string {
	buf := []byte(password)
	pwd := sha3.New256()
	pwd.Write(buf)
	return hex.EncodeToString(pwd.Sum(nil))
}

func (u *AuthRequest) BcryptPassword(password string) (string, error) {
	buf := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(buf, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func (u *AuthRequest) ValidateBcryptPassword(password, password2 string) bool {
	byteHash := []byte(password)
	buf := []byte(password2)
	err := bcrypt.CompareHashAndPassword(byteHash, buf)
	if err != nil {
		return false
	}
	return true
}

func (u *AuthRequest) Validate(v *validator.Validate) error {
	err := v.Struct(u)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return err
	}
	return nil
}
