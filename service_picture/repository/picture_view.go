package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prathoss/integration_testing/domain"
)

type PictureView struct {
	ID           uint
	ProfileID    uint
	PictureID    uint
	ViewCount    int
	LastViewedAt time.Time
}

type View interface {
	GetByProfileAndPicture(ctx context.Context, profileID, pictureID uint) (PictureView, error)
	IncrementViewCount(ctx context.Context, profileID, pictureID uint) (uint, error)
}

type pictureViewRepository struct {
	conn *pgxpool.Pool
}

func NewView(conn *pgxpool.Pool) View {
	return &pictureViewRepository{
		conn: conn,
	}
}

func (r *pictureViewRepository) GetByProfileAndPicture(
	ctx context.Context,
	profileID, pictureID uint,
) (PictureView, error) {
	var view PictureView
	err := r.conn.QueryRow(
		ctx,
		"SELECT id, profile_id, picture_id, view_count, last_viewed_at FROM picture_views WHERE profile_id = $1 AND picture_id = $2",
		profileID,
		pictureID,
	).Scan(&view.ID, &view.ProfileID, &view.PictureID, &view.ViewCount, &view.LastViewedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return PictureView{}, domain.ErrNotFound{Msg: "picture view not found"}
	}
	if err != nil {
		return PictureView{}, fmt.Errorf("selecting picture view: %w", err)
	}

	return view, nil
}

func (r *pictureViewRepository) IncrementViewCount(
	ctx context.Context,
	profileID, pictureID uint,
) (uint, error) {
	// Try to get existing view
	view, err := r.GetByProfileAndPicture(ctx, profileID, pictureID)
	if err != nil {
		// If not found, create new view
		if errors.Is(err, domain.ErrNotFound{}) {
			var viewCount uint
			err = r.conn.QueryRow(
				ctx,
				"INSERT INTO picture_views (profile_id, picture_id, view_count, last_viewed_at) VALUES ($1, $2, 1, now()) RETURNING view_count",
				profileID,
				pictureID,
			).Scan(&viewCount)
			if err != nil {
				return 0, fmt.Errorf("creating picture view: %w", err)
			}
			return viewCount, nil
		}
		return 0, err
	}

	// Update existing view
	var viewCount uint
	err = r.conn.QueryRow(
		ctx,
		"UPDATE picture_views SET view_count = view_count + 1, last_viewed_at = now() WHERE id = $1 RETURNING view_count",
		view.ID,
	).Scan(&viewCount)
	if err != nil {
		return 0, fmt.Errorf("updating picture view: %w", err)
	}

	return viewCount, nil
}
