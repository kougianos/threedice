// Package problem renders RFC 7807 responses byte-compatible with Spring's
// ProblemDetail.
package problem

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const contentType = "application/problem+json"

// Problem mirrors the JSON Spring emits for ProblemDetail.forStatusAndDetail:
// {"type":"about:blank","title":"Not Found","status":404,"detail":"...","instance":"/api/players/9"}
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func New(status int, detail, instance string) Problem {
	return Problem{
		Type: "about:blank",
		// Matches the HTTP reason phrase Spring uses for the title.
		Title:    http.StatusText(status),
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

// Write sends an RFC 7807 response, taking instance from the request path the
// way Spring populates it from the servlet request URI.
func Write(w http.ResponseWriter, r *http.Request, status int, detail string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(New(status, detail, r.URL.Path)); err != nil {
		slog.Error("writing problem response", "err", err)
	}
}
