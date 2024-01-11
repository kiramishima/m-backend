package domain

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type BondRequest struct {
	UUID       *string  `json:"uuid,omitempty"`
	Name       *string  `json:"name" validate:"required,gte=4"`
	Number     *int     `json:"number" validate:"required,gte=1, lte=10000"`
	Price      *float32 `json:"price" validate:"required,gte=0.0000, lte=1000000000.0000"`
	CurrencyID *int     `json:"currency_id" validate:"required"`
	CreatedBy  int
	Status     *string `json:"status"`
}

func (u *BondRequest) Validate(v *validator.Validate) error {
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
