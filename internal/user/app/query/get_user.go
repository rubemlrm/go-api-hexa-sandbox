package query

import (
	"context"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type UserSearch struct {
	ID user.ID
}

type GetUserHandler decorator.QueryHandler[UserSearch, *user.User]

type GetUser struct {
	userRepository user.UserRepository
}

func NewGetUserHandler(repository user.UserRepository, l *slog.Logger, tracer decorator.RecordTracer) GetUserHandler {
	return decorator.ApplyQueryDecorators[UserSearch, *user.User](
		GetUser{userRepository: repository},
		l,
		tracer,
	)
}

func (m GetUser) Handle(ctx context.Context, q UserSearch) (*user.User, error) {
	u, err := m.userRepository.Get(ctx, q.ID)
	if err != nil {
		return u, err
	}
	return u, nil
}
