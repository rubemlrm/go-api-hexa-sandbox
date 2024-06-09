package factories

import (
	"github.com/go-faker/faker/v4"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"golang.org/x/exp/rand"
)

type UserFactory struct{}

func (s *UserFactory) CreateUser() *user.User {
	return &user.User{
		ID:        user.ID(rand.Int()),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		IsEnabled: true,
	}
}

func (s *UserFactory) CreateUserCreate() *user.UserCreate {
	return &user.UserCreate{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}
}

func (s *UserFactory) CreateInvalidUserCreate() *user.UserCreate {
	return &user.UserCreate{
		Name:  faker.Name(),
		Email: faker.Password(),
	}
}
