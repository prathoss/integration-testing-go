package repository_test

import (
	"cmp"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/service_picture/repository"
	"github.com/prathoss/integration_testing/test/deps"
)

func TestPictureRepository(t *testing.T) {
	for i := range 100 {
		t.Run(
			fmt.Sprintf("query by author: %d", i), func(t *testing.T) {
				t.Parallel()
				GetByAuthorID(t)
			},
		)
	}
}

func GetByAuthorID(t *testing.T) {
	ctrl := deps.NewBuilder(t).
		WithPG(t).
		Build()

	pool, err := pgxpool.New(t.Context(), ctrl.GetPG().GetAddress(t))
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	pictureRepo := repository.NewPicture(pool)
	pictures, err := pictureRepo.GetByAuthorID(t.Context(), 1)
	if err != nil {
		t.Fatalf("could not get pictures by author: %v", err)
	}
	expected := []domain.Picture{
		{
			ID:        1,
			AuthorID:  1,
			URL:       "https://example.com/images/sunset.jpg",
			ViewCount: 5,
			CreatedAt: time.Date(2025, 4, 14, 21, 9, 43, 0, time.UTC),
		},
		{
			ID:        2,
			AuthorID:  1,
			URL:       "https://example.com/images/mountain.jpg",
			ViewCount: 1,
			CreatedAt: time.Date(2025, 4, 14, 21, 10, 43, 0, time.UTC),
		},
	}
	slices.SortFunc(
		pictures, func(a, b domain.Picture) int {
			return cmp.Compare(a.ID, b.ID)
		},
	)
	if !slices.Equal(expected, pictures) {
		t.Errorf("expected: %v, got: %v", expected, pictures)
	}
}
