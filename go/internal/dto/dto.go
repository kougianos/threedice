// Package dto holds the request and response payloads, mirroring the Java
// service's records field for field.
package dto

import (
	"github.com/kougianos/threedice/go/internal/apitime"
	"github.com/kougianos/threedice/go/internal/money"
)

// PlaceBetRequest mirrors the Java record's nullable wrapper types: a pointer
// distinguishes "absent" from a supplied zero, the way Long/BigDecimal/Integer
// distinguish null from 0.
//
// The valid-product rule is deliberately absent from these tags. Bean Validation
// reports every failing constraint on a field, whereas go-playground stops at
// the first, so predictedValue=300 must report both the max and the
// valid-product failure. Validate applies that rule separately and merges.
type PlaceBetRequest struct {
	PlayerID       *int64        `json:"playerId"       validate:"required"`
	Stake          *money.Amount `json:"stake"          validate:"required,dmin=1.00,dmax=10000.00"`
	PredictedValue *int32        `json:"predictedValue" validate:"required,min=1,max=216"`
	IdempotencyKey string        `json:"idempotencyKey" validate:"required,max=64"`
}

type PlaceBetResponse struct {
	BetID        int64        `json:"betId"`
	DieOne       int32        `json:"dieOne"`
	DieTwo       int32        `json:"dieTwo"`
	DieThree     int32        `json:"dieThree"`
	ProductValue int32        `json:"productValue"`
	Status       string       `json:"status"`
	Winnings     money.Amount `json:"winnings"`
	BalanceAfter money.Amount `json:"balanceAfter"`
}

type PlayerResponse struct {
	ID        int64        `json:"id"`
	Username  string       `json:"username"`
	Balance   money.Amount `json:"balance"`
	CreatedAt apitime.Time `json:"createdAt"`
}

type BetHistoryEntry struct {
	BetID          int64        `json:"betId"`
	PredictedValue int32        `json:"predictedValue"`
	Stake          money.Amount `json:"stake"`
	Status         string       `json:"status"`
	DieOne         int32        `json:"dieOne"`
	DieTwo         int32        `json:"dieTwo"`
	DieThree       int32        `json:"dieThree"`
	ProductValue   int32        `json:"productValue"`
	CreatedAt      apitime.Time `json:"createdAt"`
}

type TransactionHistoryEntry struct {
	TransactionID int64        `json:"transactionId"`
	BetID         int64        `json:"betId"`
	Type          string       `json:"type"`
	Amount        money.Amount `json:"amount"`
	BalanceAfter  money.Amount `json:"balanceAfter"`
	CreatedAt     apitime.Time `json:"createdAt"`
}
