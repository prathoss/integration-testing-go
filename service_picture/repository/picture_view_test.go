package repository_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prathoss/integration_testing/service_picture/repository"
	"github.com/prathoss/integration_testing/test/deps"
)

func TestPictureViewRepository_GetByProfileAndPicture(t *testing.T) {
	for range 100 {
		t.Run(
			"Success", func(t *testing.T) {
				// Create a mock DB connection
				ctrl := deps.NewBuilder(t).
					WithPG(t).
					Build()

				pool, err := pgxpool.New(t.Context(), ctrl.GetPG().GetAddress(t))
				if err != nil {
					t.Fatalf("could not connect to database: %v", err)
				}
				repo := repository.NewView(pool)

				view, err := repo.GetByProfileAndPicture(t.Context(), 3, 4)
				if err != nil {
					t.Fatalf("GetByProfileAndPicture returned error: %v", err)
				}
				expected := repository.PictureView{
					ID:           5,
					ProfileID:    3,
					PictureID:    4,
					ViewCount:    5,
					LastViewedAt: time.Date(2025, 4, 14, 22, 13, 43, 0, time.UTC),
				}
				if !equalPictureView(view, expected) {
					t.Fatalf("expected %+v, got %+v", expected, view)
				}
			},
		)
	}
}

func equalPictureView(a, b repository.PictureView) bool {
	return a.ID == b.ID &&
		a.PictureID == b.PictureID &&
		a.ProfileID == b.ProfileID &&
		a.ViewCount == b.ViewCount &&
		a.LastViewedAt.Equal(b.LastViewedAt)
}
