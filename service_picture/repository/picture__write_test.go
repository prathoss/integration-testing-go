package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/service_picture/repository"
	"github.com/prathoss/integration_testing/test/deps"
)

func TestPictureRepository_Create(t *testing.T) {
	ctrl := deps.NewBuilder(t).
		WithPG(t).
		Build()

	pool, err := pgxpool.New(t.Context(), ctrl.GetPG().GetAddress(t))
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	for i := range 100 {
		t.Run(
			fmt.Sprintf("create: %d", i), func(t *testing.T) {
				testCreatePicture(t, pool)
			},
		)
	}
}

func testCreatePicture(t *testing.T, pool *pgxpool.Pool) {
	createdAt := time.Now()
	pictureRepository := repository.NewPicture(pool)
	picture, err := pictureRepository.Create(
		t.Context(), domain.Picture{
			URL:       "http://test.gopic.io/test.jpeg",
			AuthorID:  1,
			CreatedAt: createdAt,
			ViewCount: 0,
		},
	)
	if err != nil {
		t.Fatalf("failed to create picture: %v", err)
	}

	expected := domain.Picture{
		ID:        9,
		URL:       "http://test.gopic.io/test.jpeg",
		AuthorID:  1,
		CreatedAt: createdAt,
		ViewCount: 0,
	}
	if !equalPicture(picture, expected) {
		t.Fatalf("expected %+v, got %+v", expected, picture)
	}
}
