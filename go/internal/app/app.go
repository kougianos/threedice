// Package app wires the service together, so the binary and the integration
// tests build the identical stack and differ only in the Roller they inject.
package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	threedice "github.com/kougianos/threedice/go"
	"github.com/kougianos/threedice/go/internal/config"
	"github.com/kougianos/threedice/go/internal/domain"
	"github.com/kougianos/threedice/go/internal/handler"
	"github.com/kougianos/threedice/go/internal/migrate"
	"github.com/kougianos/threedice/go/internal/postgres"
	"github.com/kougianos/threedice/go/internal/service"
	"github.com/kougianos/threedice/go/internal/validation"
)

type App struct {
	Handler http.Handler
	Pool    *pgxpool.Pool
}

// New migrates the schema, opens the pool, seeds the demo player and builds the
// router.
func New(ctx context.Context, cfg config.Config, roller domain.Roller, log *slog.Logger) (*App, error) {
	if err := migrate.Up(ctx, threedice.Migrations(), cfg.DatabaseURL, log); err != nil {
		return nil, err
	}

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	players := service.NewPlayerService(pool, log)
	if err := players.SeedDemoPlayer(ctx, cfg.InitialBalance); err != nil {
		pool.Close()
		return nil, err
	}

	validator, err := validation.New()
	if err != nil {
		pool.Close()
		return nil, err
	}

	bets := service.NewBetService(pool, roller, log)
	transactions := service.NewTransactionService(pool, log)

	return &App{
		Handler: handler.NewRouter(
			handler.NewBetHandler(bets, validator),
			handler.NewPlayerHandler(players),
			handler.NewTransactionHandler(transactions),
			threedice.Static(),
		),
		Pool: pool,
	}, nil
}

func (a *App) Close() { a.Pool.Close() }
