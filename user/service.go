package user

import (
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(user *User) (ID, error) {
	u := User{
		ID:        ID(1),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		IsEnabled: user.IsEnabled,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	id, err := s.repo.Create(&u)
	if err != nil {
		return 0, fmt.Errorf("Error creating user")
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
