package user

import (
	"fmt"
	"log/slog"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

var _ UseCase = (*Service)(nil)

func NewService(r Repository, l *slog.Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}

func (s *Service) Create(user *UserCreate) (ID, error) {
	id, err := s.repo.Create(user)
	if err != nil {
		return 0, fmt.Errorf("error creating user")
	}

	return id, nil
}

func (s *Service) Get(id ID) (*User, error) {
	u, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (s *Service) All() (*[]User, error) {
	u, err := s.repo.All()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users")
	}
	return u, nil
}
