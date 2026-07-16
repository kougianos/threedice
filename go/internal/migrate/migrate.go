// Package migrate applies the schema at startup, the way Flyway does on the
// Java side.
package migrate

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Up applies all pending migrations from files. It is a no-op when the schema
// is already current.
func Up(files fs.FS, databaseURL string, log *slog.Logger) error {
	src, err := iofs.New(files, ".")
	if err != nil {
		return fmt.Errorf("reading migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, databaseURL)
	if err != nil {
		return fmt.Errorf("opening migrator: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("Schema is up to date")
		return nil
	}
	if err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	version, _, _ := m.Version()
	log.Info("Migrations applied", "version", version)
	return nil
}
