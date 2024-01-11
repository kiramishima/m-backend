package domain

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type MarketBondRequest struct {
	MarketBondID *int `json:"id,omitempty"`
	SellerID     int  `json:"seller_id,omitempty"`
	BuyerID      int  `json:"buyer_id,omitempty"`
	Order        *int `json:"order" validate:"required,gte=1, lte=10000"`
}

func (u *MarketBondRequest) Validate(v *validator.Validate) error {
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
