package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mirrors TransactionHistoryIT.

func setupTxnPlayer(t *testing.T) int64 {
	t.Helper()
	return resetWithPlayer(t, "txnplayer", "5000.00")
}

func TestTransactionHistory_Empty(t *testing.T) {
	playerID := setupTxnPlayer(t)

	resp := get(t, fmt.Sprintf("/api/transactions/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)
	assert.Empty(t, resp.array(t))
}

func TestTransactionHistory_DebitOnLoss(t *testing.T) {
	playerID := setupTxnPlayer(t)

	// 5*5*5 = 125, predicted 12 -> LOST -> DEBIT of the stake
	script(5, 5, 5)
	require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "30", "12")).status)

	resp := get(t, fmt.Sprintf("/api/transactions/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)

	entries := resp.array(t)
	require.Len(t, entries, 1)

	e := entries[0]
	assert.Equal(t, "DEBIT", e["type"])
	assert.EqualValues(t, 30, e["amount"])
	assert.EqualValues(t, 4970, e["balanceAfter"])
	assert.NotNil(t, e["betId"])
	assert.NotNil(t, e["transactionId"])
	assert.NotEmpty(t, e["createdAt"])
}

func TestTransactionHistory_CreditOnWin(t *testing.T) {
	playerID := setupTxnPlayer(t)

	// 2*3*4 = 24, predicted 24 -> WON at 5x: winnings 250, 5000 - 50 + 250 = 5200
	script(2, 3, 4)
	require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "50", "24")).status)

	resp := get(t, fmt.Sprintf("/api/transactions/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)

	entries := resp.array(t)
	require.Len(t, entries, 1)

	e := entries[0]
	assert.Equal(t, "CREDIT", e["type"])
	assert.EqualValues(t, 250, e["amount"])
	assert.EqualValues(t, 5200, e["balanceAfter"])
}

func TestTransactionHistory_LimitedTo10(t *testing.T) {
	playerID := setupTxnPlayer(t)
	script(1, 1, 1)

	for i := 0; i < 12; i++ {
		require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "5", "100")).status)
	}

	resp := get(t, fmt.Sprintf("/api/transactions/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)
	assert.Len(t, resp.array(t), 10)
}

func TestTransactionHistory_PlayerNotFound(t *testing.T) {
	setupTxnPlayer(t)

	resp := get(t, "/api/transactions/history/99999")
	require.Equal(t, http.StatusNotFound, resp.status)
	assert.Contains(t, resp.detail(t), "Player not found")
}
