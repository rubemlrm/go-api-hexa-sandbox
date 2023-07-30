package user_postgres

import (
	"database/sql"
	"fmt"

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
	query, err := r.db.Prepare(`INSERT into user (name, email, password, isEnabled) values(?,?,?)`)
	defer query.Close()
	if err != nil {
		return 0, err
	}
	res, err := query.Exec(u.Name, u.Email, u.Password, u.IsEnabled)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return user.ID(id), nil
}

func (r *PostGresDB) Get(id user.ID) (*user.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, password, isEnabled where id = ?`)
	if err != nil {
		return nil, err
	}

	var u user.User
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("not found")
	}
	err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.IsEnabled)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
