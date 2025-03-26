package user

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
	"time"
)

type User struct {
	ID        ID        `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ID int

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required,min=3"`
}

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
