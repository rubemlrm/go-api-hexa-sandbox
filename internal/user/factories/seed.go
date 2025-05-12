package factories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"
)

func GenerateUsersOnDB(db *sql.DB, users []user.User) error {
	baseSQL := `INSERT INTO users (id, name, email, password) VALUES `

	var placeholders []string
	var args []interface{}
	for i, row := range users {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, row.ID, row.Name, row.Email, row.Password)
	}

	query := baseSQL + strings.Join(placeholders, ",")

	_, err := db.Exec(query, args...)
	return err
}

func GenerateUsers(total int) []user.User {
	uf := &UserFactory{}
	var uu []user.User
	for i := 1; i <= total; i++ {
		uu = append(uu, *uf.CreateUser())
	}

	return uu
}
