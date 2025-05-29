package user

import (
	"context"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
)

var _ validations.Validater[UserCreate] = (*UserCreate)(nil)

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
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

// UseCase Interface
type UseCase interface {
	Get(ctx context.Context, id ID) (*User, error)
	Create(ctx context.Context, user *UserCreate) (ID, error)
	All(ctx context.Context) (*[]User, error)
}

func (u *UserCreate) Check(vf validations.ValidationFunc) ([]map[string]string, error) {
	vl, err := vf("en",
		validations.WithCustomFieldLabel("json"),
		validations.WithCustomTranslation("required", "{0} must have a value!"),
	)

	if err != nil {
		return nil, err
	}

	failedValidations, err := vl.ValidateInput(u)
	if err != nil {
		return nil, err
	}

	if len(failedValidations) > 0 {
		return failedValidations, nil
	}
	return nil, nil
}

func (user User) CheckIsEnabled() (enabled bool) {
	return user.IsEnabled
}
