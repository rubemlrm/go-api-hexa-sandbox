package factories

import (
	"github.com/go-faker/faker/v4"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/models"
	"golang.org/x/exp/rand"
)

type UserFactory struct{}

func (s *UserFactory) CreateUser() *models.User {
	return &models.User{
		ID:        models.ID(rand.Int()),
		Name:      faker.Name(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		IsEnabled: true,
	}
}

func (s *UserFactory) CreateUserCreate() *models.UserCreate {
	return &models.UserCreate{
		Name:     faker.Name(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}
}

func (s *UserFactory) CreateInvalidUserCreate() *models.UserCreate {
	return &models.UserCreate{
		Name:  faker.Name(),
		Email: faker.Password(),
	}
}
