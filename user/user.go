package user

import (
	"time"
)

type ID int

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID        ID
	Name      string
	Email     string
	Password  string
	IsEnabled bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reader interface {
	Get(id ID) (*User, error)
	All() (*[]User, error)
}

type Writer interface {
	Create(u *UserCreate) (ID, error)
}

type Repository interface {
	Reader
	Writer
}

// UseCase Interface
type UseCase interface {
	Get(id ID) (*User, error)
	Create(user *UserCreate) (ID, error)

	All() (*[]User, error)
}

func (user User) CheckIsEnabled() (enabled bool) {
	return user.IsEnabled
}
