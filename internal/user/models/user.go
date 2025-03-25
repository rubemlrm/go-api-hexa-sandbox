package models

import (
	"context"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
)

func (u *UserCreate) Validate(vf validations.ValidationFunc) error {
	vl, err := vf("en",
		validations.WithCustomFieldLabel("json"),
		validations.WithCustomValidationRule("is-awesome", ValidateMyVal),
		validations.WithCustomTranslation("is-awesome", "{0} must have a value!"),
	)

	if err != nil {
		return err
	}

	err = vl.CheckWithTranslations(u)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMyVal(fl validator.FieldLevel) bool {
	return fl.Field().String() == "awesome"
}

// UseCase Interface
type UseCase interface {
	Get(ctx context.Context, id ID) (*User, error)
	Create(ctx context.Context, user *UserCreate) (ID, error)
	All(ctx context.Context) (*[]User, error)
}

func (user User) CheckIsEnabled() (enabled bool) {
	return user.IsEnabled
}
