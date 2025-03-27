package command_test

import (
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
	"testing"
)

func TestNewCreateUserHandler(t *testing.T) {
	tests := []struct {
		name             string
		expectedRequest  *user.UserCreate
		expectedStatus   int
		expectedResponse string
		mockError        error
		mockUserID       int
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
