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

type profileModel struct {
	ID        uint
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m profileModel) toProfile() domain.Profile {
	return domain.Profile{
		ID:    m.ID,
		Email: m.Email,
		Name:  m.Name,
	}
}

type Profile interface {
	GetByID(ctx context.Context, id uint) (domain.Profile, error)
	GetByEmail(ctx context.Context, email string) (domain.Profile, error)
}

func NewProfile(pgConn *pgxpool.Pool) Profile {
	return &profileRepository{
		conn: pgConn,
	}
}

var _ Profile = (*profileRepository)(nil)

type profileRepository struct {
	conn *pgxpool.Pool
}

func (u *profileRepository) GetByEmail(ctx context.Context, email string) (domain.Profile, error) {
	rows, err := u.conn.Query(ctx, "SELECT * FROM profiles WHERE email=$1", email)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("selecting profile: %v", err)
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[profileModel])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Profile{}, domain.NewErrNotFound("profile not found")
	}
	if err != nil {
		return domain.Profile{}, err
	}

	return model.toProfile(), nil
}

func (u *profileRepository) GetByID(ctx context.Context, id uint) (domain.Profile, error) {
	rows, err := u.conn.Query(ctx, "SELECT * FROM profiles WHERE id=$1", id)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("selecting profile: %v", err)
	}
	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[profileModel])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Profile{}, domain.NewErrNotFound("profile not found")
	}
	if err != nil {
		return domain.Profile{}, err
	}

	return model.toProfile(), nil
}
