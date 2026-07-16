package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mirrors BetHistoryIT.

func setupHistoryPlayer(t *testing.T) int64 {
	t.Helper()
	return resetWithPlayer(t, "historyplayer", "5000.00")
}

func TestBetHistory_Empty(t *testing.T) {
	playerID := setupHistoryPlayer(t)

	resp := get(t, fmt.Sprintf("/api/bets/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)
	assert.Empty(t, resp.array(t))
}

func TestBetHistory_SingleBet(t *testing.T) {
	playerID := setupHistoryPlayer(t)

	// 2*3*4 = 24, predicted 24 -> WIN
	script(2, 3, 4)
	require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "20", "24")).status)

	resp := get(t, fmt.Sprintf("/api/bets/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)

	entries := resp.array(t)
	require.Len(t, entries, 1)

	e := entries[0]
	assert.EqualValues(t, 24, e["predictedValue"])
	assert.EqualValues(t, 20, e["stake"])
	assert.Equal(t, "WON", e["status"])
	assert.EqualValues(t, 2, e["dieOne"])
	assert.EqualValues(t, 3, e["dieTwo"])
	assert.EqualValues(t, 4, e["dieThree"])
	assert.EqualValues(t, 24, e["productValue"])
	assert.NotNil(t, e["betId"])
	assert.NotEmpty(t, e["createdAt"])
}

func TestBetHistory_LimitedTo10_MostRecentFirst(t *testing.T) {
	playerID := setupHistoryPlayer(t)
	script(1, 1, 1) // product 1, predicted 100 -> always LOST

	for i := 0; i < 12; i++ {
		require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "5", "100")).status)
	}

	resp := get(t, fmt.Sprintf("/api/bets/history/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)
	assert.Len(t, resp.array(t), 10)
}

func TestBetHistory_PlayerNotFound(t *testing.T) {
	setupHistoryPlayer(t)

	resp := get(t, "/api/bets/history/99999")
	require.Equal(t, http.StatusNotFound, resp.status)
	assert.Contains(t, resp.detail(t), "Player not found")
}
