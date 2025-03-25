package goose

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"

	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	var sqlMigrations *sql.DB
	sqlMigrations, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	a, filename, b, _ := runtime.Caller(0)
	print(fmt.Sprintf("%d %d", a, b))
	dir := path.Join(path.Dir(filename), "../../../migrations")

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

func RollbackMigrations(dsn string) error {
	var sqlMigrations *sql.DB
	sqlMigrations, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	a, filename, b, _ := runtime.Caller(0)
	print(fmt.Sprintf("%d %d", a, b))
	dir := path.Join(path.Dir(filename), "../../../migrations")

	files := os.DirFS(dir)
	goose.SetBaseFS(files)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Down(sqlMigrations, "."); err != nil {
		panic(err)
	}
	return nil
}
