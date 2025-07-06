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

var _ store.GameVariantStore = (*PgGameVariantStore)(nil)

type PgGameVariantStore struct {
	*PgStore
}

func NewPgGameVariantStore(ps *PgStore) *PgGameVariantStore {
	return &PgGameVariantStore{
		PgStore: ps,
	}
}

func (s *PgGameVariantStore) Create(ctx context.Context, in dto.GameVariantChangeset) (model.GameVariant, error) {
	_, m := in.ToModel()

	q := t.GameVariant.INSERT(t.GameVariant.MutableColumns.Except(t.GameVariant.CreatedAt, t.GameVariant.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameVariant.AllColumns)

	var dest model.GameVariant
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_variant: %w", err)
	}
	return dest, nil
}

func (s *PgGameVariantStore) Update(ctx context.Context, in dto.GameVariantChangeset) (model.GameVariant, error) {
	cols, m := in.ToModel()

	q := t.GameVariant.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameVariant.AllColumns)

	var dest model.GameVariant
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_variant: %w", err)
	}
	return dest, nil
}

func (s *PgGameVariantStore) Delete(ctx context.Context, gamevariantID uuid.UUID) (uuid.UUID, error) {
	q := t.GameVariant.DELETE().WHERE(t.GameVariant.ID.EQ(p.UUID(gamevariantID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_variant: %w", err)
	}
	return gamevariantID, nil
}

func (s *PgGameVariantStore) Get(ctx context.Context, gamevariantID uuid.UUID) (model.GameVariant, error) {
	q := p.SELECT(t.GameVariant.AllColumns).
		FROM(t.GameVariant).
		WHERE(t.GameVariant.ID.EQ(p.UUID(gamevariantID)))

	var dest model.GameVariant
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_variant: %w", err)
	}
	return dest, nil
}

func (s *PgGameVariantStore) List(ctx context.Context) ([]model.GameVariant, error) {
	q := p.SELECT(t.GameVariant.AllColumns).
		FROM(t.GameVariant)

	var dest []model.GameVariant
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gamevariants: %w", err)
	}
	return dest, nil
}
