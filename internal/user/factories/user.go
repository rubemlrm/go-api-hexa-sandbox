package factories

import (
	"crypto/rand"
	"math/big"

	"github.com/go-faker/faker/v4"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type UserFactory struct{}

func (s *UserFactory) CreateUser() *user.User {
	mv := big.NewInt(1<<63 - 1)
	num, err := rand.Int(rand.Reader, mv)
	if err != nil {
		panic("failed to generate random number for user ID")
	}

	return &user.User{
		ID:        user.ID(num.Int64()),
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

func (s *UserFactory) CreateInvalidUserWithoutEmailCreate() *user.UserCreate {
	return &user.UserCreate{
		Name:     faker.Name(),
		Password: faker.Password(),
	}
}
