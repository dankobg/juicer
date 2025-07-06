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

var _ store.GameTimeCategoryStore = (*PgGameTimeCategoryStore)(nil)

type PgGameTimeCategoryStore struct {
	*PgStore
}

func NewPgGameTimeCategoryStore(ps *PgStore) *PgGameTimeCategoryStore {
	return &PgGameTimeCategoryStore{
		PgStore: ps,
	}
}

func (s *PgGameTimeCategoryStore) Create(ctx context.Context, in dto.GameTimeCategoryChangeset) (model.GameTimeCategory, error) {
	_, m := in.ToModel()

	q := t.GameTimeCategory.INSERT(t.GameTimeCategory.MutableColumns.Except(t.GameTimeCategory.CreatedAt, t.GameTimeCategory.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameTimeCategory.AllColumns)

	var dest model.GameTimeCategory
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_time_category: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeCategoryStore) Update(ctx context.Context, in dto.GameTimeCategoryChangeset) (model.GameTimeCategory, error) {
	cols, m := in.ToModel()

	q := t.GameTimeCategory.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameTimeCategory.AllColumns)

	var dest model.GameTimeCategory
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_time_category: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeCategoryStore) Delete(ctx context.Context, gametimecategoryID uuid.UUID) (uuid.UUID, error) {
	q := t.GameTimeCategory.DELETE().WHERE(t.GameTimeCategory.ID.EQ(p.UUID(gametimecategoryID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_time_category: %w", err)
	}
	return gametimecategoryID, nil
}

func (s *PgGameTimeCategoryStore) Get(ctx context.Context, gametimecategoryID uuid.UUID) (model.GameTimeCategory, error) {
	q := p.SELECT(t.GameTimeCategory.AllColumns).
		FROM(t.GameTimeCategory).
		WHERE(t.GameTimeCategory.ID.EQ(p.UUID(gametimecategoryID)))

	var dest model.GameTimeCategory
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_time_category: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeCategoryStore) List(ctx context.Context) ([]model.GameTimeCategory, error) {
	q := p.SELECT(t.GameTimeCategory.AllColumns).
		FROM(t.GameTimeCategory)

	var dest []model.GameTimeCategory
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gametimecategorys: %w", err)
	}
	return dest, nil
}
