package server

import "net/http"

func (s *Server) setupRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/health", s.getHealth())
	mux.Handle("GET /api/v1/profiles/{id}", s.getProfile())
}
