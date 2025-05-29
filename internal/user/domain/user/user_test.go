package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/factories"
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
		name                string
		user                *user.UserCreate
		validationFunc      validations.ValidationFunc
		expectedErr         error
		expectedValidations []map[string]string
	}{
		{
			name:                "fail validation initialization error",
			user:                &user.UserCreate{},
			validationFunc:      mockValidationFuncError,
			expectedErr:         errors.New("validation initialization error"),
			expectedValidations: nil,
		},
		{
			name:                "user validated with success",
			user:                uf.CreateUserCreate(),
			validationFunc:      validations.New,
			expectedErr:         nil,
			expectedValidations: nil,
		},
		{
			name:           "user failed to validate",
			user:           uf.CreateInvalidUserWithoutEmailCreate(),
			validationFunc: validations.New,
			expectedErr:    nil,
			expectedValidations: []map[string]string{
				{"error": "email must have a value!", "field": "email"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failedValidations, err := tt.user.Check(tt.validationFunc)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedValidations, failedValidations)
		})
	}
}
