package user_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/user"
	user_mocks "github.com/rubemlrm/go-api-bootstrap/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServiceGet(t *testing.T) {
	t.Run("user found", func(t *testing.T) {
		// arrange
		u := &user.User{
			ID:        1,
			Name:      "foo",
			Email:     "foo@bar.zsx",
			Password:  "changeme",
			IsEnabled: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		repo := user_mocks.NewMockRepository(t)

		repo.On("Get", user.ID(1)).Return(u, nil).Once()

		service := user.NewService(repo, nil)
		// act
		userFound, err := service.Get(user.ID(1))

		// assert
		assert.Nil(t, err)
		assert.Equal(t, u, userFound)
	})

	t.Run("user not userFound", func(t *testing.T) {
		// arrange
		repo := user_mocks.NewMockRepository(t)
		repo.On("Get", user.ID(2)).Return(nil, fmt.Errorf("not found")).Once()

		service := user.NewService(repo, nil)

		// act
		userFound, err := service.Get(user.ID(2))

		//
		assert.Nil(t, userFound)
		assert.Errorf(t, err, "not found")
	})
}

func TestUserCreation(t *testing.T) {
	t.Run(" Create user with success", func(t *testing.T) {
		// arrange
		u := user.UserCreate{
			Name:     "foo",
			Email:    "foo@bar.xyz",
			Password: "changeme",
		}

		repo := user_mocks.NewMockRepository(t)
		repo.On("Create", &u).Return(user.ID(1), nil).Once()
		service := user.NewService(repo, nil)
		// act
		userCreate, err := service.Create(&u)

		// assert
		assert.Equal(t, 1, userCreate)
		assert.Nil(t, err)
	})
}
