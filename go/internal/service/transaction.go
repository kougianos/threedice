package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kougianos/threedice/go/internal/apitime"
	"github.com/kougianos/threedice/go/internal/apperr"
	"github.com/kougianos/threedice/go/internal/db"
	"github.com/kougianos/threedice/go/internal/dto"
	"github.com/kougianos/threedice/go/internal/money"
)

type TransactionService struct {
	queries *db.Queries
	log     *slog.Logger
}

func NewTransactionService(pool *pgxpool.Pool, log *slog.Logger) *TransactionService {
	return &TransactionService{queries: db.New(pool), log: log}
}

// TransactionHistory returns the last 10 transactions for a player, most recent first.
func (s *TransactionService) TransactionHistory(ctx context.Context, playerID int64) ([]dto.TransactionHistoryEntry, error) {
	exists, err := s.queries.ClientExists(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("checking client: %w", err)
	}
	if !exists {
		return nil, apperr.PlayerNotFound{PlayerID: playerID}
	}

	rows, err := s.queries.ListTransactionHistory(ctx, db.ListTransactionHistoryParams{
		ClientID: playerID,
		Limit:    historyLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("listing transaction history: %w", err)
	}

	entries := make([]dto.TransactionHistoryEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, dto.TransactionHistoryEntry{
			TransactionID: r.ID,
			BetID:         r.BetID,
			Type:          r.Type,
			Amount:        money.New(r.Amount),
			BalanceAfter:  money.New(r.BalanceAfter),
			CreatedAt:     apitime.New(r.CreatedAt),
		})
	}
	return entries, nil
}
