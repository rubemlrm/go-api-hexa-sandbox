package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/pkg/validations"
	"github.com/rubemlrm/go-api-bootstrap/user/factories"

	"github.com/rubemlrm/go-api-bootstrap/user"
	"github.com/stretchr/testify/assert"
)

func mockValidationFuncError(_ string, _ ...validations.Option) (*validations.Validator, error) {
	return &validations.Validator{}, errors.New("validation initialization error")
}

func TestUserEnabled(t *testing.T) {
	u := user.User{
		Name:      "testing",
		Email:     "teste@teste.com",
		Password:  "ChangeMe",
		IsEnabled: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isEnabled := u.CheckIsEnabled()
	assert.Equal(t, isEnabled, true)
}

func TestUserValidate(t *testing.T) {
	uf := factories.UserFactory{}
	tests := []struct {
		name           string
		user           *user.UserCreate
		validationFunc validations.ValidationFunc
		expectedErr    error
	}{
		{
			name:           "fail validation initialization error",
			user:           &user.UserCreate{},
			validationFunc: mockValidationFuncError,
			expectedErr:    errors.New("validation initialization error"),
		},
		{
			name:           "user validated with success",
			user:           uf.CreateUserCreate(),
			validationFunc: validations.New,
			expectedErr:    nil,
		},
		{
			name:           "user failed to validate",
			user:           uf.CreateInvalidUserCreate(),
			validationFunc: validations.New,
			expectedErr:    errors.New("failed to validate: map[UserCreate.email:Key: 'UserCreate.email' Error:Field validation for 'email' failed on the 'email' tag]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate(tt.validationFunc)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
