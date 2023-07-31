package user_postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rubemlrm/go-api-bootstrap/user"
)

type PostGresDB struct {
	db *sql.DB
}

func NewConnection(db *sql.DB) *PostGresDB {
	return &PostGresDB{
		db: db,
	}
}

func (r *PostGresDB) Create(u *user.User) (user.ID, error) {

	if r == nil {
		return 0, fmt.Errorf("db is null")
	}
	var id int
	query := `INSERT into users (name, email, password, is_enabled) values($1,$2,$3,$4) RETURNING id`
	err := r.db.QueryRow(query, u.Name, u.Email, u.Password, u.IsEnabled).Scan(&id)
	if err != nil {
		return 0, err
	}
	return user.ID(id), nil
}

func (r *PostGresDB) Get(id user.ID) (*user.User, error) {
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
