package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/common/config"
)

func StartConnection(cfg *config.Config, logger *slog.Logger) *sql.DB {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Schema)
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		logger.Error(err.Error())
	}

	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
