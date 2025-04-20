package picture_test

import (
	"cmp"
	"slices"
	"testing"
	"time"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/service_profile/picture"
	"github.com/prathoss/integration_testing/test/deps"
)

func TestClient_GetPicturesByAuthor(t *testing.T) {
	ctrl := deps.NewBuilder(t).
		WithPG(t).
		WithPictureService(t).
		Build()

	client := picture.NewClient(ctrl.GetPictureService().GetAddress(t))

	pictures, err := client.GetPicturesByAuthor(t.Context(), 1)
	if err != nil {
		t.Fatalf("GetPicturesByAuthor failed: %v", err)
	}

	expected := []domain.Picture{
		{
			ID:        1,
			AuthorID:  1,
			URL:       "https://example.com/images/sunset.jpg",
			ViewCount: 5,
			CreatedAt: time.Date(2025, 0o4, 14, 21, 9, 43, 0, time.UTC),
		},
		{
			ID:        2,
			AuthorID:  1,
			URL:       "https://example.com/images/mountain.jpg",
			ViewCount: 1,
			CreatedAt: time.Date(2025, 0o4, 14, 21, 10, 43, 0, time.UTC),
		},
	}
	slices.SortFunc(
		pictures, func(a, b domain.Picture) int {
			return cmp.Compare(a.ID, b.ID)
		},
	)

	if !slices.Equal(expected, pictures) {
		t.Fatalf("expected: %v, got: %v", expected, pictures)
	}
}
