package server

import "net/http"

func (s *Server) setupRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/health", s.getHealth())
	mux.Handle("POST /api/v1/pictures", s.uploadPicture())
	mux.Handle("GET /api/v1/pictures", s.listPictures())
	mux.Handle("GET /api/v1/pictures/{id}", s.getPictureDetail())
}
