package goose

import (
	"database/sql"
	"os"
	"path"
	"runtime"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	var sqlMigrations *sql.DB
	sqlMigrations, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../migrations")

	files := os.DirFS(dir)
	goose.SetBaseFS(files)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlMigrations, "."); err != nil {
		panic(err)
	}
	return nil
}
