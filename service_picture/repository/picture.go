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

type pictureModel struct {
	ID        uint
	URL       string
	AuthorID  uint
	CreatedAt time.Time
}

func (m pictureModel) toPicture() domain.Picture {
	return domain.Picture{
		ID:        m.ID,
		URL:       m.URL,
		AuthorID:  m.AuthorID,
		CreatedAt: m.CreatedAt,
	}
}

type pictureWithViewCount struct {
	pictureModel
	ViewCount uint
}

func (m pictureWithViewCount) toPicture() domain.Picture {
	p := m.pictureModel.toPicture()
	p.ViewCount = m.ViewCount
	return p
}

type Picture interface {
	GetByID(ctx context.Context, id uint) (domain.Picture, error)
	GetByAuthorID(ctx context.Context, authorID uint) ([]domain.Picture, error)
	Create(ctx context.Context, picture domain.Picture) (domain.Picture, error)
}

func NewPicture(pgConn *pgxpool.Pool) Picture {
	return &pictureRepository{
		conn: pgConn,
	}
}

var _ Picture = (*pictureRepository)(nil)

type pictureRepository struct {
	conn *pgxpool.Pool
}

func (p *pictureRepository) GetByID(ctx context.Context, id uint) (domain.Picture, error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM pictures WHERE id=$1", id)
	if err != nil {
		return domain.Picture{}, fmt.Errorf("selecting picture: %w", err)
	}
	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pictureModel])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Picture{}, domain.ErrNotFound{Msg: "picture not found"}
	}
	if err != nil {
		return domain.Picture{}, err
	}

	return model.toPicture(), nil
}

func (p *pictureRepository) GetByAuthorID(
	ctx context.Context,
	authorID uint,
) ([]domain.Picture, error) {
	rows, err := p.conn.Query(
		ctx, `
			SELECT p.*, pvc.view_count
			FROM pictures p
			JOIN (
				SELECT pv.picture_id, sum(pv.view_count) AS view_count
				FROM picture_views pv
				GROUP BY pv.picture_id
			) pvc ON p.id = pvc.picture_id
			WHERE p.author_id=$1
		`, authorID,
	)
	if err != nil {
		return nil, fmt.Errorf("selecting pictures: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[pictureWithViewCount])
	if err != nil {
		return nil, err
	}

	pictures := make([]domain.Picture, len(models))
	for i, model := range models {
		pictures[i] = model.toPicture()
	}

	return pictures, nil
}

func (p *pictureRepository) Create(
	ctx context.Context,
	picture domain.Picture,
) (domain.Picture, error) {
	rows, err := p.conn.Query(
		ctx,
		"INSERT INTO pictures (url, author_id) VALUES ($1, $2) RETURNING *",
		picture.URL,
		picture.AuthorID,
	)
	if err != nil {
		return domain.Picture{}, fmt.Errorf("inserting picture: %w", err)
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[pictureModel])
	if err != nil {
		return domain.Picture{}, err
	}

	return model.toPicture(), nil
}
