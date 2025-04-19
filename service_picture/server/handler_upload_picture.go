package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/logging"
	"github.com/prathoss/integration_testing/xhttp"
)

func (s *Server) uploadPicture() http.Handler {
	return xhttp.Handler(
		func(w http.ResponseWriter, r *http.Request) error {
			authorIDStr := r.URL.Query().Get("author_id")
			authorID, err := strconv.ParseUint(authorIDStr, 10, 64)
			if err != nil {
				return domain.ErrInvalid{Msg: "author_id is required"}
			}

			// Parse multipart form
			err = r.ParseMultipartForm(10 << 20) // 10 MB max
			if err != nil {
				return fmt.Errorf("failed to parse form: %w", err)
			}

			// Get file from form
			file, header, err := r.FormFile("picture")
			if err != nil {
				return fmt.Errorf("failed to get file: %w", err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					slog.ErrorContext(r.Context(), "failed to close file", logging.Err(err))
				}
			}()

			// Check file type
			contentType := header.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				return fmt.Errorf("invalid file type: %s", contentType)
			}

			// Generate unique filename
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)

			// In a real implementation, this would upload to S3
			// For now, we'll simulate by generating a URL
			url := fmt.Sprintf("http://s3.gopic.dev/%s", filename)

			// Create picture in database
			pic, err := s.pictureRepository.Create(
				r.Context(), domain.Picture{
					URL:      url,
					AuthorID: uint(authorID),
				},
			)
			if err != nil {
				return fmt.Errorf("failed to create picture: %w", err)
			}

			// Return picture details
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			return json.NewEncoder(w).Encode(
				map[string]interface{}{
					"id":  pic.ID,
					"url": pic.URL,
				},
			)
		},
	)
}
