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

var _ store.GameTimeKindStore = (*PgGameTimeKindStore)(nil)

type PgGameTimeKindStore struct {
	*PgStore
}

func NewPgGameTimeKindStore(ps *PgStore) *PgGameTimeKindStore {
	return &PgGameTimeKindStore{
		PgStore: ps,
	}
}

func (s *PgGameTimeKindStore) Create(ctx context.Context, in dto.GameTimeKindChangeset) (model.GameTimeKind, error) {
	_, m := in.ToModel()

	q := t.GameTimeKind.INSERT(t.GameTimeKind.MutableColumns.Except(t.GameTimeKind.CreatedAt, t.GameTimeKind.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameTimeKind.AllColumns)

	var dest model.GameTimeKind
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_time_kind: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeKindStore) Update(ctx context.Context, in dto.GameTimeKindChangeset) (model.GameTimeKind, error) {
	cols, m := in.ToModel()

	q := t.GameTimeKind.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameTimeKind.AllColumns)

	var dest model.GameTimeKind
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_time_kind: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeKindStore) Delete(ctx context.Context, gametimekindID uuid.UUID) (uuid.UUID, error) {
	q := t.GameTimeKind.DELETE().WHERE(t.GameTimeKind.ID.EQ(p.UUID(gametimekindID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_time_kind: %w", err)
	}
	return gametimekindID, nil
}

func (s *PgGameTimeKindStore) Get(ctx context.Context, gametimekindID uuid.UUID) (model.GameTimeKind, error) {
	q := p.SELECT(t.GameTimeKind.AllColumns).
		FROM(t.GameTimeKind).
		WHERE(t.GameTimeKind.ID.EQ(p.UUID(gametimekindID)))

	var dest model.GameTimeKind
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_time_kind: %w", err)
	}
	return dest, nil
}

func (s *PgGameTimeKindStore) List(ctx context.Context) ([]model.GameTimeKind, error) {
	q := p.SELECT(t.GameTimeKind.AllColumns).
		FROM(t.GameTimeKind)

	var dest []model.GameTimeKind
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gametimekinds: %w", err)
	}
	return dest, nil
}
