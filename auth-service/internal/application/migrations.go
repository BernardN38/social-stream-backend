package application

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func RunDatabaseMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
