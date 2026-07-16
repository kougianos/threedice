package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mirrors PlayerControllerIT.

func TestGetPlayer_Success(t *testing.T) {
	playerID := resetWithPlayer(t, "testplayer", "1000.00")

	resp := get(t, fmt.Sprintf("/api/players/%d", playerID))
	require.Equal(t, http.StatusOK, resp.status)

	body := resp.json(t)
	assert.EqualValues(t, playerID, body["id"])
	assert.Equal(t, "testplayer", body["username"])
	assert.EqualValues(t, 1000, body["balance"])
	assert.NotEmpty(t, body["createdAt"])
}

func TestGetPlayer_NotFound(t *testing.T) {
	resetWithPlayer(t, "testplayer", "1000.00")

	resp := get(t, "/api/players/99999")
	require.Equal(t, http.StatusNotFound, resp.status)
	assert.Contains(t, resp.detail(t), "Player not found")
}
