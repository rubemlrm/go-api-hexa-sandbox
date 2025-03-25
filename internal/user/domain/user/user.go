package user

import "time"

type User struct {
	ID        ID        `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsEnabled bool      `json:"is_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ID int

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required,min=3"`
}
