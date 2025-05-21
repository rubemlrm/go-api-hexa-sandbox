package query

import (
	"context"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/decorator"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

type ListUsersHandler decorator.QueryHandler[UserSearchFilters, *[]user.User]

type UserSearchFilters struct{}
type ListUsers struct {
	userRepository user.UserRepository
}

func NewListUsersHandler(repository user.UserRepository) decorator.QueryHandler[UserSearchFilters, *[]user.User] {
	return ListUsers{
		userRepository: repository,
	}
}

func (m ListUsers) Handle(ctx context.Context, _ UserSearchFilters) (*[]user.User, error) {
	uu, err := m.userRepository.All(ctx)
	if err != nil {
		return uu, err
	}
	return uu, nil
}
