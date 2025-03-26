package query

import (
	"context"
	"database/sql"
	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/adapters"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"log/slog"
)

type ListUsersHandler decorator.QueryHandler[UserSearchFilters, *[]user.User]

type UserSearchFilters struct{}
type ListUsers struct {
	userRepository adapters.UserRepository
}

func NewListUsersHandler(db *sql.DB, logger *slog.Logger) decorator.QueryHandler[UserSearchFilters, *[]user.User] {
	return ListUsers{
		userRepository: adapters.NewUserRepository(db, logger),
	}
}

func (m ListUsers) Handle(ctx context.Context, _ UserSearchFilters) (*[]user.User, error) {
	uu, err := m.userRepository.All(ctx)
	if err != nil {
		return uu, err
	}
	return uu, nil
}
