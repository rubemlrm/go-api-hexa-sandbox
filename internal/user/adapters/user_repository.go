package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/models"
	"log/slog"

	_ "github.com/lib/pq"
)

var _ models.Repository = (*PostgresDB)(nil)

type PostgresDB struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewConnection(db *sql.DB, logger *slog.Logger) *PostgresDB {
	return &PostgresDB{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresDB) Create(ctx context.Context, u *models.UserCreate) (models.ID, error) {
	var id int
	query := `INSERT into users (name, email, password) values($1,$2,$3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, u.Name, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return models.ID(id), nil
}

func (r *PostgresDB) Get(ctx context.Context, id models.ID) (*models.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled FROM users where id = $1`)
	if err != nil {
		return nil, err
	}

	var u models.User
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("not found result")
	}
	err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.IsEnabled)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresDB) All(ctx context.Context) (*[]models.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled from users`)
	if err != nil {
		return nil, err
	}

	var uu []models.User
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.IsEnabled)
		if err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return &uu, nil
}
