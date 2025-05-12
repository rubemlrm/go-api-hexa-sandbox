package query

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type UserSearch struct {
	ID user.ID
}

type GetUserHandler decorator.QueryHandler[UserSearch, *user.User]
type GetUser struct {
	userRepository adapters.UserRepository
}

func NewGetUserHandler(db *sql.DB, logger *slog.Logger) decorator.QueryHandler[UserSearch, *user.User] {
	return GetUser{
		userRepository: adapters.NewUserRepository(db, logger),
	}
}

func (m GetUser) Handle(ctx context.Context, q UserSearch) (*user.User, error) {
	u, err := m.userRepository.Get(ctx, q.ID)
	if err != nil {
		return u, err
	}
	return u, nil
}
