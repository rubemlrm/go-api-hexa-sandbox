package testcontainers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestContainer struct {
	testcontainers.Container
	DSN string
}

func StartPostgresContainer(ctx context.Context) (*PostgresTestContainer, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "postgres",
		},
		ExposedPorts: []string{"5432/tcp"},
		Image:        "postgres:14.3",
		WaitingFor: wait.ForExec([]string{"pg_isready"}).
			WithPollInterval(2 * time.Second).
			WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
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

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "user", "password", host, mappedPort.Port(), "postgres")

	return &PostgresTestContainer{
		Container: container,
		DSN:       dsn,
	}, nil
}
