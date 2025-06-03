package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/XSAM/otelsql"
	pq "github.com/lib/pq"
)

type PostgresWrapper struct {
	*sql.DB
	host         string
	username     string
	password     string
	port         string
	schema       string
	sslmode      string
	maxOpenConns int
	maxIddleCons int
	maxIddleTime time.Duration
	tracerName   string
}

type PostgresOption func(*PostgresWrapper)

func NewConnection(logger *slog.Logger, options ...PostgresOption) (*sql.DB, error) {
	d := &PostgresWrapper{
		maxOpenConns: 10,
		maxIddleCons: 10,
		maxIddleTime: 5 * time.Minute,
	}

	// Apply options
	for _, option := range options {
		option(d)
	}

	driverName, err := otelsql.Register("postgres",
		otelsql.WithSpanNameFormatter(func(ctx context.Context, _ otelsql.Method, _ string) string {
			return d.tracerName
		}))
	if err != nil {
		return nil, err
	}

	dbURI, err := d.generateConnectionString()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	db, err := sql.Open(driverName, dbURI)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Check if connection is valid
	if err := db.Ping(); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "28P01" {
				// handle authentication error
				logger.Error("Authentication failed", "error", pqErr.Message)
				return nil, fmt.Errorf("failed to connect to database: invalid credentials")
			} else if pqErr.Code == "3D000" {
				logger.Error("Database not found", "error", pqErr.Message)
				return nil, fmt.Errorf("database not found")
			}
			fmt.Println("Postgres error code:", pqErr.Code)
		}
		logger.Error("Failed to connect to database", "error", err.Error())
		return nil, err
	}

	db.SetConnMaxIdleTime(d.maxIddleTime)
	db.SetMaxOpenConns(d.maxOpenConns)
	db.SetMaxIdleConns(d.maxIddleCons)
	return db, nil
}

func WithUsername(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.username = v
	}
}

func WithPassword(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.password = v
	}
}

func WithHost(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.host = v
	}
}

func WithSchema(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.schema = v
	}
}

func WithPort(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.port = v
	}
}

func WithSSLMode(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.sslmode = v
	}
}

func WithMaxOpenConns(v int) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.maxOpenConns = v
	}
}

func WithMaxIddleCons(v int) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.maxIddleCons = v
	}
}

func WithMaxIddleTime(v time.Duration) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.maxIddleTime = v
	}
}

func WithTracerName(v string) PostgresOption {
	return func(pw *PostgresWrapper) {
		pw.tracerName = v
	}
}

func (d *PostgresWrapper) generateConnectionString() (string, error) {
	if d.username != "" && d.password != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", d.username, d.password, d.host, d.port, d.schema, d.sslmode), nil
	}
	return "", fmt.Errorf("username and password must be provided")
}
