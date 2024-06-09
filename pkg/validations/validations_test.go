package validations

import (
	"errors"
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
			_, err := New("en", WithCustomTranslation(tt.translationTag, tt.translationMessage))
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRegisterCustomValidationRule(t *testing.T) {
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
			_, err := New("en", WithCustomValidationRule(tt.validationTag, tt.validatorFunc))

			assert.Equal(t, err, tt.expectedError)
		})
	}
}

func TestCheck(t *testing.T) {
	type Input struct {
		Email string `validate:"email"`
	}

	tests := []struct {
		name          string
		input         Input
		validatorFunc validator.Func
		expectedError error
	}{
		{
			name:          "validated struct with success",
			input:         Input{Email: faker.Email()},
			expectedError: nil,
		},
		{
			name:          "failed to validate struct",
			input:         Input{Email: ""},
			expectedError: errors.New("failed to validate: Key: 'Input.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl, err := New("en")
			assert.NoError(t, err)
			err = vl.Check(tt.input)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}

func TestCheckWithTranslations(t *testing.T) {
	type Input struct {
		Email string `validate:"email"`
	}

	tests := []struct {
		name               string
		input              Input
		translationTag     string
		translationMessage string
		expectedError      error
	}{
		{
			name:               "validated struct with success",
			input:              Input{Email: faker.Email()},
			translationTag:     "email",
			translationMessage: "email must following the email rfc",
			expectedError:      nil,
		},
		{
			name:               "failed to validate struct",
			input:              Input{Email: ""},
			translationTag:     "email",
			translationMessage: "email must following the email rfc",
			expectedError:      errors.New("failed to validate: map[Input.Email:email must following the email rfc]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl, err := New("en", WithCustomTranslation(tt.translationTag, tt.translationMessage))

			assert.NoError(t, err)
			err = vl.CheckWithTranslations(tt.input)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}
