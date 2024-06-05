package user

import (
	"time"
)

type ID int

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

type User struct {
	ID        ID        `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
