package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/store"
	"github.com/google/uuid"

	p "github.com/go-jet/jet/v2/postgres"
)

var _ store.RatingStore = (*PgRatingStore)(nil)

type PgRatingStore struct {
	*PgStore
}

func NewPgRatingStore(ps *PgStore) *PgRatingStore {
	return &PgRatingStore{
		PgStore: ps,
	}
}

func (s *PgRatingStore) BatchCreate(ctx context.Context, in []dto.RatingChangeset) ([]model.Rating, error) {
	ms := make([]model.Rating, 0, len(in))
	for _, x := range in {
		_, m := x.ToModel()
		ms = append(ms, m)
	}
	q := t.Rating.INSERT(t.Rating.MutableColumns.Except(t.Rating.CreatedAt, t.Rating.UpdatedAt)).
		MODELS(ms).
		RETURNING(t.Rating.AllColumns)

	var dest []model.Rating
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to batch create a rating: %w", err)
	}
	return dest, nil
}

func (s *PgRatingStore) Create(ctx context.Context, in dto.RatingChangeset) (model.Rating, error) {
	_, m := in.ToModel()

	q := t.Rating.INSERT(t.Rating.MutableColumns.Except(t.Rating.CreatedAt, t.Rating.UpdatedAt)).
		MODEL(m).
		RETURNING(t.Rating.AllColumns)

	var dest model.Rating
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a rating: %w", err)
	}
	return dest, nil
}

func (s *PgRatingStore) Update(ctx context.Context, in dto.RatingChangeset) (model.Rating, error) {
	cols, m := in.ToModel()

	q := t.Rating.UPDATE(cols).
		MODEL(m).
		RETURNING(t.Rating.AllColumns)

	var dest model.Rating
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a rating: %w", err)
	}
	return dest, nil
}

func (s *PgRatingStore) Delete(ctx context.Context, ratingID uuid.UUID) (uuid.UUID, error) {
	q := t.Rating.DELETE().WHERE(t.Rating.ID.EQ(p.UUID(ratingID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a rating: %w", err)
	}
	return ratingID, nil
}

func (s *PgRatingStore) Get(ctx context.Context, ratingID uuid.UUID) (model.Rating, error) {
	q := p.SELECT(t.Rating.AllColumns).
		FROM(t.Rating).
		WHERE(t.Rating.ID.EQ(p.UUID(ratingID)))

	var dest model.Rating
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a rating: %w", err)
	}
	return dest, nil
}

func (s *PgRatingStore) List(ctx context.Context) ([]model.Rating, error) {
	q := p.SELECT(t.Rating.AllColumns).
		FROM(t.Rating)

	var dest []model.Rating
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}
	return dest, nil
}
