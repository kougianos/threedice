package handler

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter wires the API and serves the embedded frontend at /, the way Spring
// serves src/main/resources/static.
func NewRouter(
	bets *BetHandler,
	players *PlayerHandler,
	transactions *TransactionHandler,
	static fs.FS,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.NotFound(NotFound)
	r.MethodNotAllowed(MethodNotAllowed)

	r.Route("/api", func(api chi.Router) {
		api.Mount("/bets", bets.Routes())
		api.Mount("/players", players.Routes())
		api.Mount("/transactions", transactions.Routes())
	})

	r.Handle("/*", staticHandler(static))

	return r
}
