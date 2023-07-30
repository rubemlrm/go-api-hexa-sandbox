package user

import (
	"time"
)

type ID int

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
}

type Writer interface {
	Create(u *User) (ID, error)
}

type Repository interface {
	Reader
	Writer
}

// UseCase Interface
type UseCase interface {
	Get(id ID) (*User, error)
	Create(user *User) (ID, error)
}

func (user User) CheckIsEnabled() (enabled bool) {
	return user.IsEnabled
}
