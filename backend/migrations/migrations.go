package migrations

import (
	"database/sql"
	"embed"

	"github.com/bongofriend/bongo-notes/backend/lib/config"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationsFS embed.FS

func ApplyMigrations(c config.Config) error {
	dbConfig := c.Db
	database, err := sql.Open(dbConfig.Driver, dbConfig.Path)
	if err != nil {
		return err
	}
	defer database.Close()
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect(dbConfig.Driver); err != nil {
		return err
	}
	if err := goose.Up(database, "."); err != nil {
		return err
	}
	return nil
}
