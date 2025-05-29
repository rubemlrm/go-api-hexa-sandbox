package validations_test

import (
	"errors"
	validations "github.com/rubemlrm/go-api-bootstrap/internal/common/validations"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/validations/mocks"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCustomTranslation(t *testing.T) {
	tests := []struct {
		name               string
		translationTag     string
		translationMessage string
		expectedError      error
	}{
		{
			name:               "register translation with success",
			translationTag:     "testing-tag",
			translationMessage: "error validating x",
			expectedError:      nil,
		},
		{
			name:               "register translation failed because of empty tag",
			translationTag:     "",
			translationMessage: "error validating x",
			expectedError:      errors.New("tag name can't be empty"),
		},
		{
			name:               "register translation failed because of empty message",
			translationTag:     "testing-tag",
			translationMessage: "",
			expectedError:      errors.New("message can't be empty"),
		},
		{
			name:               "register translation failed",
			translationTag:     "testing-tag",
			translationMessage: "",
			expectedError:      errors.New("message can't be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validations.New("en", validations.WithCustomTranslation(tt.translationTag, tt.translationMessage))
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestWithCustomFieldLabel(t *testing.T) {
	type input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"-" validate:"required"`
	}

	tests := []struct {
		name                string
		labelValue          string
		expectedError       error
		expectErrorOnNew    bool
		input               *input
		expectedValidations []map[string]string
	}{
		{
			name:                "failed to define custom field label",
			labelValue:          "",
			expectedError:       errors.New("custom field label is required"),
			input:               &input{},
			expectErrorOnNew:    true,
			expectedValidations: nil,
		},
		{
			name:          "check input name after failed struct validation",
			labelValue:    "json",
			expectedError: nil,
			input: &input{
				Email:    "test",
				Password: faker.Password(),
			},
			expectErrorOnNew: false,
			expectedValidations: []map[string]string{
				{"error": "Key: 'input.email' Error:Field validation for 'email' failed on the 'email' tag", "field": "email"},
			},
		},
		{
			name:          "check input omit after failed struct validation",
			labelValue:    "json",
			expectedError: nil,
			input: &input{
				Email:    faker.Email(),
				Password: "",
			},
			expectErrorOnNew: false,
			expectedValidations: []map[string]string{
				{"error": "Key: 'input.Password' Error:Field validation for 'Password' failed on the 'required' tag", "field": "Password"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl, err := validations.New("en", validations.WithCustomFieldLabel(tt.labelValue))
			if tt.expectErrorOnNew == true && tt.expectedValidations == nil {
				assert.Equal(t, tt.expectedError, err)
				assert.Error(t, err)
				return
			}
			failedValidations, err := vl.ValidateInput(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedValidations, failedValidations)
		})
	}
}

func TestWithCustomValidationRule(t *testing.T) {
	tests := []struct {
		name          string
		validationTag string
		validatorFunc validator.Func
		expectedError error
	}{
		{
			name:          "register translation with success",
			validationTag: "testing-tag",
			validatorFunc: func(fl validator.FieldLevel) bool { return false },
			expectedError: nil,
		},
		{
			name:          "register translation failed because of empty tag",
			validationTag: "",
			validatorFunc: func(fl validator.FieldLevel) bool { return false },
			expectedError: errors.New("tag name can't be empty"),
		},
		{
			name:          "register translation failed because of empty function",
			validationTag: "testing-tag",
			validatorFunc: nil,
			expectedError: errors.New("function can't be nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validations.New("en", validations.WithCustomValidationRule(tt.validationTag, tt.validatorFunc))

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestCheck(t *testing.T) {
	type Input struct {
		Email string `validate:"email"`
	}

	tests := []struct {
		name                string
		input               Input
		validatorFunc       validator.Func
		expectedError       error
		mockedError         error
		expectedValidations []map[string]string
	}{
		{
			name:                "validated struct with success",
			input:               Input{Email: faker.Email()},
			mockedError:         nil,
			expectedError:       nil,
			expectedValidations: nil,
		},
		{
			name:                "failed to validate struct",
			input:               Input{Email: ""},
			mockedError:         errors.New("key: 'Input.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
			expectedError:       errors.New("failed to validate: key: 'Input.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
			expectedValidations: []map[string]string(nil),
		},
		{
			name:          "failed to validate struct fields with empty email",
			input:         Input{Email: ""},
			mockedError:   nil,
			expectedError: nil,
			expectedValidations: []map[string]string{
				{"error": "Key: 'Input.Email' Error:Field validation for 'Email' failed on the 'email' tag", "field": "Email"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vl *validations.Validator
			var err error
			if tt.mockedError != nil {
				mockStructValidator := mocks.NewMockStructValidator(t)
				vl = &validations.Validator{
					Validate:     mockStructValidator,
					Translations: nil,
				}
				mockStructValidator.On("Struct", tt.input).Return(tt.mockedError)
			} else {
				vl, err = validations.New("en")
				assert.NoError(t, err)
			}

			failedValidations, err := vl.ValidateInput(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedValidations, failedValidations)
		})
	}
}

func TestCheckWithTranslations(t *testing.T) {
	type Input struct {
		Email string `validate:"email"`
	}

	tests := []struct {
		name                string
		input               Input
		translationTag      string
		translationMessage  string
		expectedError       error
		expectedValidations []map[string]string
	}{
		{
			name:                "validated struct with success",
			input:               Input{Email: faker.Email()},
			translationTag:      "email",
			translationMessage:  "email must following the email rfc",
			expectedError:       nil,
			expectedValidations: nil,
		},
		{
			name:               "failed to validate struct",
			input:              Input{Email: ""},
			translationTag:     "email",
			translationMessage: "email must following the email rfc",
			expectedError:      nil,
			expectedValidations: []map[string]string{
				{"error": "email must following the email rfc", "field": "Email"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl, err := validations.New("en", validations.WithCustomTranslation(tt.translationTag, tt.translationMessage))

			assert.NoError(t, err)
			failedValidations, err := vl.ValidateInput(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedValidations, failedValidations)
		})
	}
}
