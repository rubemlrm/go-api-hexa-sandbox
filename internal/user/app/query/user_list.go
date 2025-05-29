package query

import (
	"context"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type ListUsersHandler decorator.QueryHandler[UserSearchFilters, *[]user.User]

type UserSearchFilters struct{}
type ListUsers struct {
	userRepository user.UserRepository
}

func NewListUsersHandler(repository user.UserRepository, l *slog.Logger) ListUsersHandler {
	return decorator.ApplyQueryDecorators[UserSearchFilters, *[]user.User](
		ListUsers{userRepository: repository},
		l,
	)
}

func (m ListUsers) Handle(ctx context.Context, _ UserSearchFilters) (*[]user.User, error) {
	uu, err := m.userRepository.All(ctx)
	if err != nil {
		return uu, err
	}
	return uu, nil
}
