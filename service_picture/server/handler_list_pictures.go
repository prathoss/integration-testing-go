package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/xhttp"
)

func (s *Server) listPictures() http.Handler {
	return xhttp.Handler(
		func(w http.ResponseWriter, r *http.Request) error {
			authorIDStr := r.URL.Query().Get("author")
			if authorIDStr == "" {
				return domain.ErrInvalid{Msg: "author query parameter is required"}
			}
			authorID, err := strconv.ParseUint(authorIDStr, 10, 64)
			if err != nil {
				return domain.ErrInvalid{Msg: "author query parameter must be a number"}
			}
			pictures, err := s.pictureRepository.GetByAuthorID(r.Context(), uint(authorID))
			if err != nil {
				return err
			}
			return json.NewEncoder(w).Encode(pictures)
		},
	)
}
