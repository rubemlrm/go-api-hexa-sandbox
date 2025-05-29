package app

import (
	user "github.com/rubemlrm/go-api-bootstrap/internal/user/app"
)

type Application struct {
	UserModule user.UserModule
}
