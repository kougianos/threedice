// Package handler is the HTTP layer: routing, request binding, and the error
// mapping that mirrors the Java GlobalExceptionHandler.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kougianos/threedice/go/internal/apperr"
	"github.com/kougianos/threedice/go/internal/problem"
)

// writeJSON sends a success payload. Errors after the header is written can
// only be logged.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("writing response", "err", err)
	}
}

// writeError maps a business error onto the status the Java handler uses.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		notFound     apperr.PlayerNotFound
		insufficient apperr.InsufficientBalance
		duplicate    apperr.DuplicateBet
		invalid      apperr.Validation
	)

	switch {
	case errors.As(err, &notFound):
		slog.Warn(err.Error())
		problem.Write(w, r, http.StatusNotFound, err.Error())
	case errors.As(err, &insufficient):
		slog.Warn(err.Error())
		problem.Write(w, r, http.StatusBadRequest, err.Error())
	case errors.As(err, &duplicate):
		slog.Warn(err.Error())
		problem.Write(w, r, http.StatusConflict, err.Error())
	case errors.As(err, &invalid):
		slog.Warn("Validation failed", "errors", err.Error())
		problem.Write(w, r, http.StatusBadRequest, err.Error())
	default:
		slog.Error("Unexpected error", "err", err)
		problem.Write(w, r, http.StatusInternalServerError, "An unexpected error occurred")
	}
}

// parsePlayerID mirrors @PathVariable @Positive on the Java controllers. The
// method-name prefix reproduces the property path Spring puts in the detail,
// e.g. "getPlayer.playerId: Player ID must be positive".
func parsePlayerID(w http.ResponseWriter, r *http.Request, raw, javaMethod string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		// Matches the Java handler for MethodArgumentTypeMismatchException,
		// which names the parameter that failed to convert.
		problem.Write(w, r, http.StatusBadRequest,
			"playerId must be a number: "+raw)
		return 0, false
	}
	if id <= 0 {
		problem.Write(w, r, http.StatusBadRequest,
			javaMethod+".playerId: Player ID must be positive")
		return 0, false
	}
	return id, true
}

// NotFound replaces chi's bare 404 with an RFC 7807 body.
func NotFound(w http.ResponseWriter, r *http.Request) {
	problem.Write(w, r, http.StatusNotFound, "No endpoint "+r.Method+" "+r.URL.Path)
}

// MethodNotAllowed replaces chi's bare 405 with an RFC 7807 body.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	problem.Write(w, r, http.StatusMethodNotAllowed,
		"Method "+r.Method+" is not supported for "+r.URL.Path)
}
