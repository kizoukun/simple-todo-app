package helper

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitHelper() {
	Validate = validator.New()
}

func ValidateRequest(structType any) error {
	err := Validate.Struct(structType)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return ve
		}
		return err
	}

	return nil
}
