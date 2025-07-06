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

var _ store.GameStateStore = (*PgGameStateStore)(nil)

type PgGameStateStore struct {
	*PgStore
}

func NewPgGameStateStore(ps *PgStore) *PgGameStateStore {
	return &PgGameStateStore{
		PgStore: ps,
	}
}

func (s *PgGameStateStore) Create(ctx context.Context, in dto.GameStateChangeset) (model.GameState, error) {
	_, m := in.ToModel()

	q := t.GameState.INSERT(t.GameState.MutableColumns.Except(t.GameState.CreatedAt, t.GameState.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameState.AllColumns)

	var dest model.GameState
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_state: %w", err)
	}
	return dest, nil
}

func (s *PgGameStateStore) Update(ctx context.Context, in dto.GameStateChangeset) (model.GameState, error) {
	cols, m := in.ToModel()

	q := t.GameState.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameState.AllColumns)

	var dest model.GameState
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_state: %w", err)
	}
	return dest, nil
}

func (s *PgGameStateStore) Delete(ctx context.Context, gamestateID uuid.UUID) (uuid.UUID, error) {
	q := t.GameState.DELETE().WHERE(t.GameState.ID.EQ(p.UUID(gamestateID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_state: %w", err)
	}
	return gamestateID, nil
}

func (s *PgGameStateStore) Get(ctx context.Context, gamestateID uuid.UUID) (model.GameState, error) {
	q := p.SELECT(t.GameState.AllColumns).
		FROM(t.GameState).
		WHERE(t.GameState.ID.EQ(p.UUID(gamestateID)))

	var dest model.GameState
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_state: %w", err)
	}
	return dest, nil
}

func (s *PgGameStateStore) List(ctx context.Context) ([]model.GameState, error) {
	q := p.SELECT(t.GameState.AllColumns).
		FROM(t.GameState)

	var dest []model.GameState
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gamestates: %w", err)
	}
	return dest, nil
}
