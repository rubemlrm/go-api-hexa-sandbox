package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"

	_ "github.com/lib/pq"
)

var _ user.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r UserRepository) Create(ctx context.Context, u *user.UserCreate) (user.ID, error) {
	var id int
	query := `INSERT into users (name, email, password) values($1,$2,$3) RETURNING id`
	err := r.db.QueryRow(query, u.Name, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return user.ID(id), nil
}

func (r UserRepository) Get(ctx context.Context, id user.ID) (*user.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled FROM users where id = $1`)
	if err != nil {
		return nil, err
	}

	var u user.User
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

func (r UserRepository) All(ctx context.Context) (*[]user.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled from users`)
	if err != nil {
		return nil, err
	}

	var uu []user.User
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	for rows.Next() {
		var u user.User
		err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.IsEnabled)
		if err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return &uu, nil
}
