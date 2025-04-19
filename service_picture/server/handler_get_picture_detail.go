package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/xhttp"
)

func (s *Server) getPictureDetail() http.Handler {
	return xhttp.Handler(
		func(w http.ResponseWriter, r *http.Request) error {
			profileIDStr := r.URL.Query().Get("profile_id")
			profileID, err := strconv.ParseUint(profileIDStr, 10, 64)
			if err != nil {
				return domain.NewErrInvalid("profile_id is required")
			}

			// Get picture ID from URL
			idStr := r.PathValue("id")
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid picture ID: %w", err)
			}

			// Increment view count using the user view repository
			viewCount, err := s.pictureViewRepository.IncrementViewCount(
				r.Context(),
				uint(profileID),
				uint(id),
			)
			if err != nil {
				return fmt.Errorf("failed to increment view count: %w", err)
			}

			// Get picture from database
			pic, err := s.pictureRepository.GetByID(r.Context(), uint(id))
			if err != nil {
				return err
			}
			pic.ViewCount = viewCount

			// Return picture details with view count
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(pic)
		},
	)
}
