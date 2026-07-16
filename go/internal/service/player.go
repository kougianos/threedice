package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/kougianos/threedice/go/internal/apitime"
	"github.com/kougianos/threedice/go/internal/apperr"
	"github.com/kougianos/threedice/go/internal/db"
	"github.com/kougianos/threedice/go/internal/dto"
	"github.com/kougianos/threedice/go/internal/money"
)

type PlayerService struct {
	queries *db.Queries
	log     *slog.Logger
}

func NewPlayerService(pool *pgxpool.Pool, log *slog.Logger) *PlayerService {
	return &PlayerService{queries: db.New(pool), log: log}
}

// SeedDemoPlayer creates the default player when no clients exist, mirroring
// PlayerService's CommandLineRunner.
func (s *PlayerService) SeedDemoPlayer(ctx context.Context, initialBalance decimal.Decimal) error {
	count, err := s.queries.CountClients(ctx)
	if err != nil {
		return fmt.Errorf("counting clients: %w", err)
	}
	if count > 0 {
		return nil
	}

	demo, err := s.queries.CreateClient(ctx, db.CreateClientParams{
		Username: "player1",
		Balance:  initialBalance,
	})
	if err != nil {
		return fmt.Errorf("creating demo player: %w", err)
	}

	s.log.Info("Created demo player",
		"id", demo.ID, "username", demo.Username, "balance", demo.Balance)
	return nil
}

func (s *PlayerService) GetPlayer(ctx context.Context, playerID int64) (dto.PlayerResponse, error) {
	client, err := s.queries.GetClient(ctx, playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.PlayerResponse{}, apperr.PlayerNotFound{PlayerID: playerID}
	}
	if err != nil {
		return dto.PlayerResponse{}, fmt.Errorf("loading client: %w", err)
	}

	return dto.PlayerResponse{
		ID:        client.ID,
		Username:  client.Username,
		Balance:   money.New(client.Balance),
		CreatedAt: apitime.New(client.CreatedAt),
	}, nil
}
