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
