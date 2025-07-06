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

var _ store.GameResultStore = (*PgGameResultStore)(nil)

type PgGameResultStore struct {
	*PgStore
}

func NewPgGameResultStore(ps *PgStore) *PgGameResultStore {
	return &PgGameResultStore{
		PgStore: ps,
	}
}

func (s *PgGameResultStore) Create(ctx context.Context, in dto.GameResultChangeset) (model.GameResult, error) {
	_, m := in.ToModel()

	q := t.GameResult.INSERT(t.GameResult.MutableColumns.Except(t.GameResult.CreatedAt, t.GameResult.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameResult.AllColumns)

	var dest model.GameResult
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_result: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStore) Update(ctx context.Context, in dto.GameResultChangeset) (model.GameResult, error) {
	cols, m := in.ToModel()

	q := t.GameResult.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameResult.AllColumns)

	var dest model.GameResult
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_result: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStore) Delete(ctx context.Context, gameresultID uuid.UUID) (uuid.UUID, error) {
	q := t.GameResult.DELETE().WHERE(t.GameResult.ID.EQ(p.UUID(gameresultID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_result: %w", err)
	}
	return gameresultID, nil
}

func (s *PgGameResultStore) Get(ctx context.Context, gameresultID uuid.UUID) (model.GameResult, error) {
	q := p.SELECT(t.GameResult.AllColumns).
		FROM(t.GameResult).
		WHERE(t.GameResult.ID.EQ(p.UUID(gameresultID)))

	var dest model.GameResult
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_result: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStore) List(ctx context.Context) ([]model.GameResult, error) {
	q := p.SELECT(t.GameResult.AllColumns).
		FROM(t.GameResult)

	var dest []model.GameResult
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gameresults: %w", err)
	}
	return dest, nil
}
