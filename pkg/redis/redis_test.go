package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/pkg/redis"
	"github.com/rubemlrm/go-api-bootstrap/tests/testcontainers"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	_, err := testcontainers.StartRedisContainer(ctx)
	assert.NoError(t, err)

	cl, err := redis.New()
	assert.Error(t, err)
	assert.Nil(t, cl)
}

func TestRedisClientAddress(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	container, err := testcontainers.StartRedisContainer(ctx)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		address       string
		expectedError error
	}{
		{
			name:          "test client and connection with valid address",
			address:       container.DSN,
			expectedError: nil,
		},
		{
			name:          "failed to setup client because of empty address",
			address:       "",
			expectedError: errors.New("address can't be empty"),
		},
		{
			name:          "failed to setup client because of invalid address",
			address:       "testingtesting",
			expectedError: errors.New("invalid connection string invalid redis URL scheme: "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = redis.New(redis.WithAddr(tt.address))
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRedisClientWithPassword(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	container, err := testcontainers.StartRedisContainerWithAuth(ctx)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "test client and connection with valid password",
			password:      "password",
			expectedError: nil,
		},
		{
			name:          "failed to setup client because of empty password",
			password:      "",
			expectedError: errors.New("password can't be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = redis.New(
				redis.WithAddr(container.DSN),
				redis.WithAuthentication(tt.password))
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
