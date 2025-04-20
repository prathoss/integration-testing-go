package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/service_picture/repository"
)

func TestPictureViewRepository_GetByProfileAndPicture(t *testing.T) {
	// Create a mock DB connection
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("could not create pool mock: %v", err)
	}

	var ppool repository.Pool = pool
	repo := repository.NewView(ppool)

	t.Run(
		"Success", func(t *testing.T) {
			// Mock expected query and response
			expectedQuery := "SELECT id, profile_id, picture_id, view_count, last_viewed_at FROM picture_views WHERE profile_id = \\$1 AND picture_id = \\$2"

			viewTime := time.Now().Round(time.Second)
			rows := pgxmock.NewRows([]string{"id", "profile_id", "picture_id", "view_count", "last_viewed_at"}).
				AddRow(1, 2, 3, 5, viewTime)

			pool.ExpectQuery(expectedQuery).
				WithArgs(uint(2), uint(3)).
				WillReturnRows(rows)

			// Call the method
			view, err := repo.GetByProfileAndPicture(t.Context(), 2, 3)
			if err != nil {
				t.Fatalf("GetByProfileAndPicture returned error: %v", err)
			}
			expected := repository.PictureView{
				ID:           1,
				ProfileID:    2,
				PictureID:    3,
				ViewCount:    5,
				LastViewedAt: viewTime,
			}
			if !equalPictureView(view, expected) {
				t.Fatalf("expected %+v, got %+v", expected, view)
			}
		},
	)

	t.Run(
		"Not Found", func(t *testing.T) {
			// Mock expected query and no rows response
			expectedQuery := "SELECT id, profile_id, picture_id, view_count, last_viewed_at FROM picture_views WHERE profile_id = \\$1 AND picture_id = \\$2"

			pool.ExpectQuery(expectedQuery).
				WithArgs(uint(4), uint(5)).
				WillReturnError(pgx.ErrNoRows)

			// Call the method
			_, err := repo.GetByProfileAndPicture(t.Context(), 4, 5)
			var enf *domain.ErrNotFound
			if !errors.As(err, &enf) {
				t.Fatalf("expected to get not found error, got %T", err)
			}
		},
	)

	// Ensure all expectations were met
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func equalPictureView(a, b repository.PictureView) bool {
	return a.ID == b.ID &&
		a.PictureID == b.PictureID &&
		a.ProfileID == b.ProfileID &&
		a.ViewCount == b.ViewCount &&
		a.LastViewedAt.Equal(b.LastViewedAt)
}
