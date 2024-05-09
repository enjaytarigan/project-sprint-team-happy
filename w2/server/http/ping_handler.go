package httpserver

import (
	"net/http"
)

func (s *HttpServer) handlePing(w http.ResponseWriter, r *http.Request) {

	s.writeJSON(w, r, http.StatusOK, map[string]any{"message": "pong"})
}
