package command

import (
	"context"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type CreateUserHandler decorator.CommandHandler[*user.UserCreate, user.ID]

type UserCreateStore struct {
	userRepository user.UserRepository
}

func NewCreateUserHandler(repository user.UserRepository) decorator.CommandHandler[*user.UserCreate, user.ID] {
	return UserCreateStore{
		userRepository: repository,
	}
}

func (m UserCreateStore) Handle(ctx context.Context, cmd *user.UserCreate) (user.ID, error) {
	u, err := m.userRepository.Create(ctx, cmd)
	if err != nil {
		return 0, err
	}
	return u, nil
}
