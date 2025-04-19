package server

import (
	"net/http"

	"github.com/prathoss/integration_testing/xhttp"
)

func (s *Server) getHealth() http.Handler {
	return xhttp.Handler(
		func(w http.ResponseWriter, r *http.Request) error {
			var n int
			if err := s.dbConnPool.QueryRow(r.Context(), "SELECT 1").Scan(&n); err != nil {
				return err
			}
			w.WriteHeader(http.StatusOK)
			return nil
		},
	)
}
