package command

import (
	"context"
	"database/sql"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"log/slog"
)

type CreateUserHandler decorator.CommandHandler[*user.UserCreate, user.ID]

type UserCreateStore struct {
	userRepository adapters.UserRepository
}

func NewCreateUserHandler(db *sql.DB, logger *slog.Logger) decorator.CommandHandler[*user.UserCreate, user.ID] {
	return UserCreateStore{
		userRepository: adapters.NewUserRepository(db, logger),
	}
}

func (m UserCreateStore) Handle(ctx context.Context, cmd *user.UserCreate) (user.ID, error) {
	u, err := m.userRepository.Create(ctx, cmd)
	if err != nil {
		return 0, err
	}
	return u, nil
}
