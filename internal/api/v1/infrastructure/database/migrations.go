package database

import (
	"github.com/pressly/goose/v3"
)

func MigrateDB() {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(Db, "/scripts"); err != nil {
		panic(err)
	}
}