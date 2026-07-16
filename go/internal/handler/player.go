package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kougianos/threedice/go/internal/service"
)

type PlayerHandler struct {
	players *service.PlayerService
}

func NewPlayerHandler(players *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{players: players}
}

func (h *PlayerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{playerId}", h.GetPlayer)
	return r
}

// GetPlayer returns the player's details including their current balance.
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID, ok := parsePlayerID(w, r, chi.URLParam(r, "playerId"), "getPlayer")
	if !ok {
		return
	}

	player, err := h.players.GetPlayer(r.Context(), playerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, player)
}
