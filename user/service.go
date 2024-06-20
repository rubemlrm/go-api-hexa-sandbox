package user

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

var _ UseCase = (*Service)(nil)

var tracer = otel.Tracer("gin-server")

func NewService(r Repository, l *slog.Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}

func (s *Service) Create(ctx context.Context, user *UserCreate) (ID, error) {
	_, span := tracer.Start(ctx, "create user", oteltrace.WithAttributes(attribute.String("test", "123")))
	defer span.End()
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("error creating user")
	}

	return id, nil
}

func (s *Service) Get(ctx context.Context, id ID) (*User, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (s *Service) All(ctx context.Context) (*[]User, error) {
	u, err := s.repo.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users")
	}
	return u, nil
}
