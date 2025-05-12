package user

import "context"

type UserReader interface {
	Get(ctx context.Context, id ID) (*User, error)
	All(ctx context.Context) (*[]User, error)
}

type UserWriter interface {
	Create(ctx context.Context, u *UserCreate) (ID, error)
}

type UserRepository interface {
	UserReader
	UserWriter
}
