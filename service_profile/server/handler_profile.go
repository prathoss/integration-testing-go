package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/xhttp"
)

func (s *Server) getProfile() http.Handler {
	return xhttp.Handler(
		func(w http.ResponseWriter, r *http.Request) error {
			profileIDStr := r.PathValue("id")
			if profileIDStr == "" {
				return domain.NewErrInvalid("profile id is required")
			}

			profileID, err := strconv.ParseUint(profileIDStr, 10, 64)
			if err != nil {
				return domain.NewErrInvalid(err.Error())
			}

			profile, err := s.profileRepository.GetByID(r.Context(), uint(profileID))
			if err != nil {
				return err
			}

			pictures, err := s.pictureClient.GetPicturesByAuthor(r.Context(), uint(profileID))
			if err != nil {
				return err
			}

			return json.NewEncoder(w).Encode(
				domain.ProfileFeed{
					Profile:  profile,
					Pictures: pictures,
				},
			)
		},
	)
}
