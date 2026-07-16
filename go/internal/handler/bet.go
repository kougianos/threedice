package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kougianos/threedice/go/internal/apperr"
	"github.com/kougianos/threedice/go/internal/dto"
	"github.com/kougianos/threedice/go/internal/problem"
	"github.com/kougianos/threedice/go/internal/service"
	"github.com/kougianos/threedice/go/internal/validation"
)

type BetHandler struct {
	bets      *service.BetService
	validator *validation.Validator
}

func NewBetHandler(bets *service.BetService, v *validation.Validator) *BetHandler {
	return &BetHandler{bets: bets, validator: v}
}

func (h *BetHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.PlaceBet)
	r.Get("/history/{playerId}", h.BetHistory)
	return r
}

// PlaceBet rolls three dice, settles the outcome, updates the balance and
// returns the result.
func (h *BetHandler) PlaceBet(w http.ResponseWriter, r *http.Request) {
	var req dto.PlaceBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Spring surfaces an unreadable body as HttpMessageNotReadableException.
		problem.Write(w, r, http.StatusBadRequest, "Malformed JSON request body")
		return
	}

	if detail := h.validator.ValidatePlaceBet(req); detail != "" {
		writeError(w, r, apperr.Validation{Detail: detail})
		return
	}

	resp, err := h.bets.PlaceBet(r.Context(), service.PlaceBetCommand{
		PlayerID:       *req.PlayerID,
		Stake:          req.Stake.Decimal,
		PredictedValue: *req.PredictedValue,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// BetHistory returns the last 10 bets for a player, most recent first.
func (h *BetHandler) BetHistory(w http.ResponseWriter, r *http.Request) {
	playerID, ok := parsePlayerID(w, r, chi.URLParam(r, "playerId"), "getBetHistory")
	if !ok {
		return
	}

	entries, err := h.bets.BetHistory(r.Context(), playerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, entries)
}
