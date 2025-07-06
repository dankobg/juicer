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

var _ store.GameResultStatusStore = (*PgGameResultStatusStore)(nil)

type PgGameResultStatusStore struct {
	*PgStore
}

func NewPgGameResultStatusStore(ps *PgStore) *PgGameResultStatusStore {
	return &PgGameResultStatusStore{
		PgStore: ps,
	}
}

func (s *PgGameResultStatusStore) Create(ctx context.Context, in dto.GameResultStatusChangeset) (model.GameResultStatus, error) {
	_, m := in.ToModel()

	q := t.GameResultStatus.INSERT(t.GameResultStatus.MutableColumns.Except(t.GameResultStatus.CreatedAt, t.GameResultStatus.UpdatedAt)).
		MODEL(m).
		RETURNING(t.GameResultStatus.AllColumns)

	var dest model.GameResultStatus
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create a game_result_status: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStatusStore) Update(ctx context.Context, in dto.GameResultStatusChangeset) (model.GameResultStatus, error) {
	cols, m := in.ToModel()

	q := t.GameResultStatus.UPDATE(cols).
		MODEL(m).
		RETURNING(t.GameResultStatus.AllColumns)

	var dest model.GameResultStatus
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game_result_status: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStatusStore) Delete(ctx context.Context, gameresultstatusID uuid.UUID) (uuid.UUID, error) {
	q := t.GameResultStatus.DELETE().WHERE(t.GameResultStatus.ID.EQ(p.UUID(gameresultstatusID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game_result_status: %w", err)
	}
	return gameresultstatusID, nil
}

func (s *PgGameResultStatusStore) Get(ctx context.Context, gameresultstatusID uuid.UUID) (model.GameResultStatus, error) {
	q := p.SELECT(t.GameResultStatus.AllColumns).
		FROM(t.GameResultStatus).
		WHERE(t.GameResultStatus.ID.EQ(p.UUID(gameresultstatusID)))

	var dest model.GameResultStatus
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game_result_status: %w", err)
	}
	return dest, nil
}

func (s *PgGameResultStatusStore) List(ctx context.Context) ([]model.GameResultStatus, error) {
	q := p.SELECT(t.GameResultStatus.AllColumns).
		FROM(t.GameResultStatus)

	var dest []model.GameResultStatus
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list gameresultstatuss: %w", err)
	}
	return dest, nil
}
