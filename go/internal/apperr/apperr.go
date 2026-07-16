// Package apperr holds the business errors the HTTP layer maps to statuses,
// mirroring the Java service's custom exceptions and GlobalExceptionHandler.
package apperr

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// PlayerNotFound maps to 404, as PlayerNotFoundException does.
type PlayerNotFound struct{ PlayerID int64 }

func (e PlayerNotFound) Error() string {
	return fmt.Sprintf("Player not found with ID: %d", e.PlayerID)
}

// InsufficientBalance maps to 400, as InsufficientBalanceException does.
type InsufficientBalance struct {
	Balance decimal.Decimal
	Stake   decimal.Decimal
}

func (e InsufficientBalance) Error() string {
	return fmt.Sprintf("Insufficient balance: current=%s, requested stake=%s",
		e.Balance.StringFixed(2), e.Stake.StringFixed(2))
}

// DuplicateBet maps to 409, as DuplicateBetException does.
type DuplicateBet struct{ IdempotencyKey string }

func (e DuplicateBet) Error() string {
	return fmt.Sprintf("Duplicate bet submission: idempotencyKey=%s", e.IdempotencyKey)
}

// Validation carries the joined field errors that Spring's handler builds from
// a MethodArgumentNotValidException. Maps to 400.
type Validation struct{ Detail string }

func (e Validation) Error() string { return e.Detail }
