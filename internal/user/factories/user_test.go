package factories_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/factories"
	"github.com/stretchr/testify/assert"
)

func TestUserFactory_CreateUser(t *testing.T) {
	uf := factories.UserFactory{}
	u := uf.CreateUser()
	assert.NotNil(t, u)
	assert.NotEmpty(t, u.Name)
	assert.NotEmpty(t, u.Email)
	assert.NotEmpty(t, u.Password)
	assert.True(t, u.IsEnabled)
}

func TestUserFactory_CreateUserCreate(t *testing.T) {
	uf := factories.UserFactory{}
	uc := uf.CreateUserCreate()
	assert.NotNil(t, uc)
	assert.NotEmpty(t, uc.Name)
	assert.NotEmpty(t, uc.Email)
	assert.NotEmpty(t, uc.Password)
}

func TestUserFactory_CreateInvalidUserCreate(t *testing.T) {
	uf := factories.UserFactory{}
	uc := uf.CreateInvalidUserCreate()
	assert.NotNil(t, uc)
	assert.NotEmpty(t, uc.Name)
	assert.NotEmpty(t, uc.Email)
	// Email should not be a valid email, as it's generated with faker.Password()
	assert.NotEqual(t, uc.Email, faker.Email())
}

func TestUserFactory_CreateInvalidUserWithoutEmailCreate(t *testing.T) {
	uf := factories.UserFactory{}
	uc := uf.CreateInvalidUserWithoutEmailCreate()
	assert.NotNil(t, uc)
	assert.NotEmpty(t, uc.Name)
	assert.Empty(t, uc.Email)
	assert.NotEqual(t, uc.Email, faker.Email())
}
