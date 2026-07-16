package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mirrors BetControllerIT.

func setupBetPlayer(t *testing.T) int64 {
	t.Helper()
	return resetWithPlayer(t, "testplayer", "1000.00")
}

// ───── Core Bet Placement ─────

func TestPlaceBet_Win(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(1, 3, 4)

	resp := postBet(t, betBody(playerID, "10", "12"))
	require.Equal(t, http.StatusCreated, resp.status, "body: %s", resp.body)

	body := resp.json(t)
	assert.EqualValues(t, 1, body["dieOne"])
	assert.EqualValues(t, 3, body["dieTwo"])
	assert.EqualValues(t, 4, body["dieThree"])
	assert.EqualValues(t, 12, body["productValue"])
	assert.Equal(t, "WON", body["status"])
	assert.EqualValues(t, 50, body["winnings"])
	assert.EqualValues(t, 1040, body["balanceAfter"])
}

func TestPlaceBet_Lose(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(2, 3, 5)

	resp := postBet(t, betBody(playerID, "10", "12"))
	require.Equal(t, http.StatusCreated, resp.status, "body: %s", resp.body)

	body := resp.json(t)
	assert.Equal(t, "LOST", body["status"])
	assert.EqualValues(t, 30, body["productValue"])
	assert.EqualValues(t, 0, body["winnings"])
	assert.EqualValues(t, 990, body["balanceAfter"])
}

func TestPlaceBet_MultipleRounds(t *testing.T) {
	playerID := setupBetPlayer(t)

	// 6*6*6 = 216, predicted 12 -> LOST, 1000 - 50 = 950
	script(6, 6, 6)
	resp := postBet(t, betBody(playerID, "50", "12"))
	require.Equal(t, http.StatusCreated, resp.status)
	assert.Equal(t, "LOST", resp.json(t)["status"])
	assert.EqualValues(t, 950, resp.json(t)["balanceAfter"])

	// 2*2*2 = 8, predicted 8 -> WON at 2x, 950 - 100 + 200 = 1050
	script(2, 2, 2)
	resp = postBet(t, betBody(playerID, "100", "8"))
	require.Equal(t, http.StatusCreated, resp.status)
	body := resp.json(t)
	assert.Equal(t, "WON", body["status"])
	assert.EqualValues(t, 200, body["winnings"])
	assert.EqualValues(t, 1050, body["balanceAfter"])
}

// ───── Balance Consistency ─────

func TestPlaceBet_RapidSequentialBets_BalanceConsistent(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(6, 6, 6)

	for i := 0; i < 10; i++ {
		resp := postBet(t, betBody(playerID, "50", "12"))
		require.Equal(t, http.StatusCreated, resp.status, "bet %d body: %s", i, resp.body)
	}

	resp := get(t, fmt.Sprintf("/api/players/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)
	assert.EqualValues(t, 500, resp.json(t)["balance"])
}

// ───── Request Validation ─────

func TestPlaceBet_PlayerNotFound(t *testing.T) {
	setupBetPlayer(t)

	resp := postBet(t, betBody(99999, "10", "12"))
	require.Equal(t, http.StatusNotFound, resp.status)
	assert.Contains(t, resp.detail(t), "Player not found")
}

func TestPlaceBet_InsufficientBalance(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "5000", "12"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Insufficient balance")
}

func TestPlaceBet_ValidationError_MissingFields(t *testing.T) {
	setupBetPlayer(t)

	resp := postBet(t, `{}`)
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Player ID is required")
}

func TestPlaceBet_ValidationError_PredictedValueOutOfRange(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "10", "300"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Predicted value must be at most 216")
}

func TestPlaceBet_ValidationError_NegativeStake(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "-5", "12"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Minimum stake is $1.00")
}

func TestPlaceBet_ValidationError_ImpossibleProduct(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "10", "7"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "not a possible product of three dice")
}

func TestPlaceBet_ValidEdgeCaseProducts(t *testing.T) {
	playerID := setupBetPlayer(t)

	script(1, 1, 1)
	resp := postBet(t, betBody(playerID, "1", "1"))
	assert.Equal(t, http.StatusCreated, resp.status, "product 1 body: %s", resp.body)

	script(6, 6, 6)
	resp = postBet(t, betBody(playerID, "1", "216"))
	assert.Equal(t, http.StatusCreated, resp.status, "product 216 body: %s", resp.body)
}

// ───── Stake Limits ─────

func TestStake_BelowMinimum_Returns400(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "0.50", "12"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Minimum stake is $1.00")
}

func TestStake_AboveMaximum_Returns400(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, betBody(playerID, "10001", "12"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Maximum stake is $10,000.00")
}

func TestStake_AtMinimum_Succeeds(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(1, 1, 1)

	resp := postBet(t, betBody(playerID, "1.00", "12"))
	require.Equal(t, http.StatusCreated, resp.status, "body: %s", resp.body)
	assert.EqualValues(t, 999, resp.json(t)["balanceAfter"])
}

func TestStake_AtMaximum_RejectedByBalance(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(1, 1, 1)

	resp := postBet(t, betBody(playerID, "10000", "12"))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Insufficient balance")
}

// ───── Idempotency Key ─────

func TestIdempotency_DuplicateBet_Returns409(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(1, 1, 1)

	key := uniqueKey()
	body := fmt.Sprintf(
		`{"playerId":%d,"stake":10,"predictedValue":12,"idempotencyKey":%q}`, playerID, key)

	resp := postBet(t, body)
	require.Equal(t, http.StatusCreated, resp.status, "first bet body: %s", resp.body)

	resp = postBet(t, body)
	require.Equal(t, http.StatusConflict, resp.status)
	assert.Contains(t, resp.detail(t), "Duplicate bet submission")
}

func TestIdempotency_DifferentKeys_BothSucceed(t *testing.T) {
	playerID := setupBetPlayer(t)
	script(1, 1, 1)

	require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "10", "12")).status)
	require.Equal(t, http.StatusCreated, postBet(t, betBody(playerID, "10", "12")).status)
}

func TestIdempotency_MissingKey_Returns400(t *testing.T) {
	playerID := setupBetPlayer(t)

	resp := postBet(t, fmt.Sprintf(
		`{"playerId":%d,"stake":10,"predictedValue":12}`, playerID))
	require.Equal(t, http.StatusBadRequest, resp.status)
	assert.Contains(t, resp.detail(t), "Idempotency key is required")
}
