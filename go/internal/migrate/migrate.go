// Package migrate applies the schema at startup, the way Flyway does on the
// Java side.
package migrate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
)

// Up applies all pending migrations from files. It is a no-op when the schema
// is already current, and it adopts a Flyway-managed schema if it finds one.
func Up(ctx context.Context, files fs.FS, databaseURL string, log *slog.Logger) error {
	if err := adoptFlywaySchema(ctx, databaseURL, log); err != nil {
		return err
	}

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

// adoptFlywaySchema lets this service take over a database the Java service
// built.
//
// Both services own the same schema but track it with different tools: Flyway
// records state in flyway_schema_history, golang-migrate in schema_migrations.
// Pointed at a Flyway-built database, golang-migrate would find no history of
// its own, conclude the database is empty, re-run the first migration and die
// on "relation clients already exists" -- leaving itself marked dirty.
//
// So when we find Flyway's history and none of our own, we record Flyway's
// version as ours and let Up carry on from there. A database with neither table
// is untouched and migrates from scratch as usual.
func adoptFlywaySchema(ctx context.Context, databaseURL string, log *slog.Logger) error {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connecting to check schema history: %w", err)
	}
	defer conn.Close(ctx)

	var ours bool
	if err := conn.QueryRow(ctx,
		`SELECT to_regclass('public.schema_migrations') IS NOT NULL`).Scan(&ours); err != nil {
		return fmt.Errorf("looking for schema_migrations: %w", err)
	}
	if ours {
		return nil // we already track this database
	}

	var flyway bool
	if err := conn.QueryRow(ctx,
		`SELECT to_regclass('public.flyway_schema_history') IS NOT NULL`).Scan(&flyway); err != nil {
		return fmt.Errorf("looking for flyway_schema_history: %w", err)
	}
	if !flyway {
		return nil // fresh database
	}

	// Flyway stores version as text and allows dotted versions; only whole
	// numbers map onto golang-migrate's bigint, and ours are V1 and V2.
	var version int64
	if err := conn.QueryRow(ctx, `
		SELECT COALESCE(MAX(version::bigint), 0)
		FROM flyway_schema_history
		WHERE success AND version ~ '^[0-9]+$'`).Scan(&version); err != nil {
		return fmt.Errorf("reading flyway version: %w", err)
	}
	if version == 0 {
		return nil // Flyway is present but has applied nothing
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin adopt tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op once committed

	// Matches the table golang-migrate's postgres driver would create itself.
	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint  NOT NULL PRIMARY KEY,
			dirty   boolean NOT NULL
		)`); err != nil {
		return fmt.Errorf("creating schema_migrations: %w", err)
	}

	// DO NOTHING so two instances starting together cannot collide.
	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)
		 ON CONFLICT (version) DO NOTHING`, version); err != nil {
		return fmt.Errorf("recording adopted version: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit adopt tx: %w", err)
	}

	log.Warn("Adopted a Flyway-managed schema; continuing from its version",
		"version", version)
	return nil
}
