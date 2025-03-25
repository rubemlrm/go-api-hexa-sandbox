package user

import "context"

type Reader interface {
	Get(ctx context.Context, id ID) (*User, error)
	All(ctx context.Context) (*[]User, error)
}

type Writer interface {
	Create(ctx context.Context, u *UserCreate) (ID, error)
}

type Repository interface {
	Reader
	Writer
}
