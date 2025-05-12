package testcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisTestContainer struct {
	testcontainers.Container
	DSN string
}

func StartRedisContainer(ctx context.Context) (*RedisTestContainer, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"REDIS_DATABASE":       "testing",
			"ALLOW_EMPTY_PASSWORD": "yes",
		},
		ExposedPorts: []string{"6379/tcp"},
		Image:        "bitnami/redis:latest",
		WaitingFor: wait.ForExec([]string{"redis-cli", "ping"}).
			WithPollInterval(500 * time.Millisecond).
			WithStartupTimeout(15 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		panic(err)
	}
	return &RedisTestContainer{
		Container: container,
		DSN:       fmt.Sprintf("redis://%s:%s", host, mappedPort.Port()),
	}, nil
}

func StartRedisContainerWithAuth(ctx context.Context) (*RedisTestContainer, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"REDIS_PASSWORD": "password",
			"REDIS_DATABASE": "testing",
		},
		ExposedPorts: []string{"6379/tcp"},
		Image:        "bitnami/redis:latest",
		WaitingFor: wait.ForLog(".*Ready to accept connections").
			WithPollInterval(2 * time.Second).AsRegexp(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		panic(err)
	}
	return &RedisTestContainer{
		Container: container,
		DSN:       fmt.Sprintf("redis://redis:testing@%s:%s", host, mappedPort.Port()),
	}, nil
}
