package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/rubemlrm/go-api-bootstrap/user"
)

var _ user.Repository = (*PostgresDB)(nil)

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

func (r *PostgresDB) Create(u *user.UserCreate) (user.ID, error) {
	var id int
	query := `INSERT into users (name, email, password) values($1,$2,$3) RETURNING id`
	err := r.db.QueryRow(query, u.Name, u.Email, u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return user.ID(id), nil
}

func (r *PostgresDB) Get(id user.ID) (*user.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled FROM users where id = $1`)
	if err != nil {
		return nil, err
	}

	var u user.User
	rows, err := stmt.Query(id)
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

func (r *PostgresDB) All() (*[]user.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, is_enabled from users`)
	if err != nil {
		return nil, err
	}

	var uu []user.User
	rows, err := stmt.Query()
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
