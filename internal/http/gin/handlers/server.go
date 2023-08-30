package handlers

import (
	"github.com/rubemlrm/go-api-bootstrap/internal/http/gin/openapi"
)

type server struct{}

func NewServer() openapi.ServerInterface {
	return &server{}
}
