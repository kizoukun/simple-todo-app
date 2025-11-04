package helper

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/kizoukun/codingtest/entity"
)

var Validate *validator.Validate

type ctxKey string

const UserCtxKey ctxKey = "user"

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

func UserFromContext(ctx context.Context) (*entity.User, bool) {
	u, ok := ctx.Value(UserCtxKey).(*entity.User)
	return u, ok
}
