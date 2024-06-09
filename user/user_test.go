package user_test

import (
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/user"
	"github.com/stretchr/testify/assert"
)

func TestUserEnabled(t *testing.T) {
	user := user.User{
		Name:      "testing",
		Email:     "teste@teste.com",
		Password:  "ChangeMe",
		IsEnabled: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	isEnabled := user.CheckIsEnabled()
	assert.Equal(t, isEnabled, true)
}

func TestUserCreateValidate(t *testing.T) {
	user := user.UserCreate{
		Name:     "testing",
		Email:    "teste@teste.pt",
		Password: "awesome",
	}

	err := user.Validate()
	assert.NoError(t, err)
}
