package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kougianos/threedice/go/internal/service"
)

type TransactionHandler struct {
	transactions *service.TransactionService
}

func NewTransactionHandler(transactions *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactions: transactions}
}

func (h *TransactionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/history/{playerId}", h.TransactionHistory)
	return r
}

// TransactionHistory returns the last 10 transactions for a player, most recent first.
func (h *TransactionHandler) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	playerID, ok := parsePlayerID(w, r, chi.URLParam(r, "playerId"), "getTransactionHistory")
	if !ok {
		return
	}

	entries, err := h.transactions.TransactionHistory(r.Context(), playerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, entries)
}
