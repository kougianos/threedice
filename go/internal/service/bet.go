// Package service holds the business logic, mirroring the Java @Service classes.
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
	"github.com/kougianos/threedice/go/internal/domain"
	"github.com/kougianos/threedice/go/internal/dto"
	"github.com/kougianos/threedice/go/internal/money"
)

const (
	statusWon  = "WON"
	statusLost = "LOST"

	typeDebit  = "DEBIT"
	typeCredit = "CREDIT"

	historyLimit = 10
)

type BetService struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	roller  domain.Roller
	log     *slog.Logger
}

func NewBetService(pool *pgxpool.Pool, roller domain.Roller, log *slog.Logger) *BetService {
	return &BetService{pool: pool, queries: db.New(pool), roller: roller, log: log}
}

// PlaceBetCommand is the validated form of dto.PlaceBetRequest.
type PlaceBetCommand struct {
	PlayerID       int64
	Stake          decimal.Decimal
	PredictedValue int32
	IdempotencyKey string
}

// PlaceBet deducts the stake, rolls the dice, settles the outcome, updates the
// balance and records the bet, draw and ledger entry -- all in one transaction.
func (s *BetService) PlaceBet(ctx context.Context, cmd PlaceBetCommand) (dto.PlaceBetResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op once committed

	q := s.queries.WithTx(tx)

	// Reject duplicate submissions early.
	duplicate, err := q.BetExistsByIdempotencyKey(ctx, cmd.IdempotencyKey)
	if err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("checking idempotency key: %w", err)
	}
	if duplicate {
		return dto.PlaceBetResponse{}, apperr.DuplicateBet{IdempotencyKey: cmd.IdempotencyKey}
	}

	// Pessimistic lock: SELECT ... FOR UPDATE prevents concurrent balance reads.
	client, err := q.GetClientForUpdate(ctx, cmd.PlayerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.PlaceBetResponse{}, apperr.PlayerNotFound{PlayerID: cmd.PlayerID}
	}
	if err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("locking client: %w", err)
	}

	// Validate sufficient balance.
	if client.Balance.LessThan(cmd.Stake) {
		return dto.PlaceBetResponse{}, apperr.InsufficientBalance{
			Balance: client.Balance,
			Stake:   cmd.Stake,
		}
	}

	// Deduct stake immediately.
	balance := client.Balance.Sub(cmd.Stake)

	// Roll three dice. Kept as separate statements so the roller is called in a
	// guaranteed order -- the tests depend on the sequence.
	d1 := s.roller.Roll()
	d2 := s.roller.Roll()
	d3 := s.roller.Roll()
	productValue := d1 * d2 * d3

	// Determine outcome.
	won := productValue == cmd.PredictedValue
	status := statusLost
	winnings := decimal.Zero
	if won {
		status = statusWon
		winnings = cmd.Stake.Mul(domain.Odds(productValue))
		balance = balance.Add(winnings)
	}

	if err := q.UpdateClientBalance(ctx, db.UpdateClientBalanceParams{
		ID:      client.ID,
		Balance: balance,
	}); err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("updating balance: %w", err)
	}

	bet, err := q.CreateBet(ctx, db.CreateBetParams{
		ClientID:       client.ID,
		PredictedValue: cmd.PredictedValue,
		Stake:          cmd.Stake,
		IdempotencyKey: cmd.IdempotencyKey,
		Status:         status,
	})
	if err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("creating bet: %w", err)
	}

	if _, err := q.CreateDraw(ctx, db.CreateDrawParams{
		BetID:        bet.ID,
		DieOne:       d1,
		DieTwo:       d2,
		DieThree:     d3,
		ProductValue: productValue,
	}); err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("creating draw: %w", err)
	}

	// On a loss the ledger records a DEBIT of the stake; on a win a CREDIT of
	// the winnings, each with a balance snapshot.
	txType, amount := typeDebit, cmd.Stake
	if won {
		txType, amount = typeCredit, winnings
	}
	if _, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
		BetID:        bet.ID,
		ClientID:     client.ID,
		Type:         txType,
		Amount:       amount,
		BalanceAfter: balance,
	}); err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("creating transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.PlaceBetResponse{}, fmt.Errorf("commit: %w", err)
	}

	s.log.Info("Bet placed",
		"playerId", client.ID,
		"predicted", cmd.PredictedValue,
		"actual", productValue,
		"status", status,
		"balance", balance)

	return dto.PlaceBetResponse{
		BetID:        bet.ID,
		DieOne:       d1,
		DieTwo:       d2,
		DieThree:     d3,
		ProductValue: productValue,
		Status:       status,
		Winnings:     money.New(winnings),
		BalanceAfter: money.New(balance),
	}, nil
}

// BetHistory returns the last 10 bets for a player, most recent first.
func (s *BetService) BetHistory(ctx context.Context, playerID int64) ([]dto.BetHistoryEntry, error) {
	exists, err := s.queries.ClientExists(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("checking client: %w", err)
	}
	if !exists {
		return nil, apperr.PlayerNotFound{PlayerID: playerID}
	}

	rows, err := s.queries.ListBetHistory(ctx, db.ListBetHistoryParams{
		ClientID: playerID,
		Limit:    historyLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("listing bet history: %w", err)
	}

	// Non-nil so an empty history marshals as [] rather than null.
	entries := make([]dto.BetHistoryEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, dto.BetHistoryEntry{
			BetID:          r.ID,
			PredictedValue: r.PredictedValue,
			Stake:          money.New(r.Stake),
			Status:         r.Status,
			DieOne:         r.DieOne,
			DieTwo:         r.DieTwo,
			DieThree:       r.DieThree,
			ProductValue:   r.ProductValue,
			CreatedAt:      apitime.New(r.CreatedAt),
		})
	}
	return entries, nil
}
