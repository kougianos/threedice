package integration

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	threedice "github.com/kougianos/threedice/go"
	"github.com/kougianos/threedice/go/internal/migrate"
)

// These cover taking over the database the Java service built, which is what
// happens when a deployment is switched from Spring Boot to Go in place.

func quietLog() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// newDatabase creates an empty database on the shared container and returns its
// URL.
func newDatabase(t *testing.T, name string) string {
	t.Helper()
	ctx := context.Background()

	admin, err := pgx.Connect(ctx, containerDSN)
	require.NoError(t, err)
	defer admin.Close(ctx)

	_, err = admin.Exec(ctx, `DROP DATABASE IF EXISTS `+pgx.Identifier{name}.Sanitize()+` WITH (FORCE)`)
	require.NoError(t, err)
	_, err = admin.Exec(ctx, `CREATE DATABASE `+pgx.Identifier{name}.Sanitize())
	require.NoError(t, err)

	t.Cleanup(func() {
		c, err := pgx.Connect(context.Background(), containerDSN)
		if err != nil {
			return
		}
		defer c.Close(context.Background())
		_, _ = c.Exec(context.Background(), `DROP DATABASE IF EXISTS `+pgx.Identifier{name}.Sanitize()+` WITH (FORCE)`)
	})

	return strings.Replace(containerDSN, "/"+containerDB+"?", "/"+name+"?", 1)
}

// seedFlywaySchema builds the schema the way the Java service leaves it: the
// tables, plus Flyway's own history table recording V1 and V2.
func seedFlywaySchema(t *testing.T, dsn string) {
	t.Helper()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	defer conn.Close(ctx)

	// The same DDL Flyway would have applied, taken from the shared migrations.
	for _, name := range []string{"000001_initial_schema.up.sql", "000002_add_idempotency_key.up.sql"} {
		body, err := fs.ReadFile(threedice.Migrations(), name)
		require.NoError(t, err)
		_, err = conn.Exec(ctx, string(body))
		require.NoError(t, err, "applying %s", name)
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE flyway_schema_history (
			installed_rank INT           NOT NULL PRIMARY KEY,
			version        VARCHAR(50),
			description    VARCHAR(200)  NOT NULL,
			type           VARCHAR(20)   NOT NULL,
			script         VARCHAR(1000) NOT NULL,
			checksum       INT,
			installed_by   VARCHAR(100)  NOT NULL,
			installed_on   TIMESTAMP     NOT NULL DEFAULT NOW(),
			execution_time INT           NOT NULL,
			success        BOOLEAN       NOT NULL
		);
		INSERT INTO flyway_schema_history VALUES
			(1, '1', 'initial schema',      'SQL', 'V1__initial_schema.sql',      1, 'threedice', NOW(), 10, true),
			(2, '2', 'add idempotency key', 'SQL', 'V2__add_idempotency_key.sql', 2, 'threedice', NOW(), 10, true);
	`)
	require.NoError(t, err)
}

func TestMigrate_AdoptsFlywaySchema(t *testing.T) {
	ctx := context.Background()
	dsn := newDatabase(t, "adopt_flyway_test")
	seedFlywaySchema(t, dsn)

	// Without adoption this fails with `relation "clients" already exists`.
	require.NoError(t, migrate.Up(ctx, threedice.Migrations(), dsn, quietLog()))

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	defer conn.Close(ctx)

	var version int64
	var dirty bool
	require.NoError(t, conn.QueryRow(ctx,
		`SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty))
	assert.EqualValues(t, 2, version, "should adopt Flyway's version")
	assert.False(t, dirty, "adoption must not leave the schema dirty")

	// Running again must stay a no-op rather than re-adopting.
	require.NoError(t, migrate.Up(ctx, threedice.Migrations(), dsn, quietLog()))

	var rows int
	require.NoError(t, conn.QueryRow(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&rows))
	assert.Equal(t, 1, rows, "golang-migrate expects exactly one history row")
}

func TestMigrate_AdoptionPreservesExistingData(t *testing.T) {
	ctx := context.Background()
	dsn := newDatabase(t, "adopt_data_test")
	seedFlywaySchema(t, dsn)

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	defer conn.Close(ctx)

	// A player mid-game, as the live database would have.
	_, err = conn.Exec(ctx,
		`INSERT INTO clients (username, balance) VALUES ('player1', 1234.56)`)
	require.NoError(t, err)

	require.NoError(t, migrate.Up(ctx, threedice.Migrations(), dsn, quietLog()))

	var username, balance string
	require.NoError(t, conn.QueryRow(ctx,
		`SELECT username, balance::text FROM clients WHERE id = 1`).Scan(&username, &balance))
	assert.Equal(t, "player1", username)
	assert.Equal(t, "1234.56", balance, "the live balance must survive the switch")
}

func TestMigrate_FreshDatabaseStillMigratesFromScratch(t *testing.T) {
	ctx := context.Background()
	dsn := newDatabase(t, "fresh_migrate_test")

	require.NoError(t, migrate.Up(ctx, threedice.Migrations(), dsn, quietLog()))

	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	defer conn.Close(ctx)

	var version int64
	require.NoError(t, conn.QueryRow(ctx, `SELECT version FROM schema_migrations`).Scan(&version))
	assert.EqualValues(t, 2, version)

	for _, table := range []string{"clients", "bets", "draws", "transactions"} {
		var exists bool
		require.NoError(t, conn.QueryRow(ctx,
			`SELECT to_regclass('public.'||$1) IS NOT NULL`, table).Scan(&exists))
		assert.True(t, exists, "table %s should have been created", table)
	}
}
