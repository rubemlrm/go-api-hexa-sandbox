package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/config"
)

func StartConnection(cfg *config.Config) *sql.DB {
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Schema)
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
