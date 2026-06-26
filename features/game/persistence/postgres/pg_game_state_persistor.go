package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/pagination"
	"github.com/dankobg/juicer/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ game.GameStatePersistor = (*PgGameStatePersistor)(nil)

type PgGameStatePersistor struct {
	*postgres.PgPersistor
}

func NewPgGameStatePersistor(pst *postgres.PgPersistor) *PgGameStatePersistor {
	return &PgGameStatePersistor{
		PgPersistor: pst,
	}
}

func convertGameStatePgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameStatePersistor) ListGameStates(ctx context.Context, filters game.ListGameStatesFilters) (pagination.WithTotal[models.GameState], error) {
	q := psql.Select(
		sm.Columns(models.GameStates.Columns),
		sm.From(models.GameStates.Name()),
		sm.GroupBy(models.GameStates.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameStates.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameStates.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameStates.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameStatesRow struct {
		models.GameState
		TotalCount int64
	}

	gameStatesRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameStatesRow]())
	if err != nil {
		return pagination.WithTotal[models.GameState]{}, fmt.Errorf("query gamestates")
	}

	gameStates := make([]models.GameState, len(gameStatesRows))
	for i, row := range gameStatesRows {
		gameStates[i] = row.GameState
	}

	var total int64
	if len(gameStatesRows) > 0 {
		total = gameStatesRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameStates, total)

	return out, nil
}

func (pst *PgGameStatePersistor) GetGameStateByID(ctx context.Context, gamestateID int64) (models.GameState, error) {
	q := psql.Select(
		sm.Columns(models.GameStates.Columns),
		sm.From(models.GameStates.Name()),
		sm.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))),
	)

	gamestate, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("query gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) CreateGameState(ctx context.Context, in models.GameStateSetter) (models.GameState, error) {
	q := models.GameStates.Insert(&in, im.Returning(models.GameStates.Columns))

	gamestate, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("insert gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) UpdateGameState(ctx context.Context, gamestateID int64, in models.GameStateSetter) (models.GameState, error) {
	q := models.GameStates.Update(
		in.UpdateMod(),
		um.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))),
		um.Returning(models.GameStates.Columns),
	)

	gamestate, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("update gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) DeleteGameStateByID(ctx context.Context, gamestateID int64) (int64, error) {
	q := models.GameStates.Delete(dm.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gamestate: %w", err)
	}

	return gamestateID, nil
}

func (pst *PgGameStatePersistor) BulkDeleteGameStates(ctx context.Context, ids []int64) error {
	gamestateIDs := make([]any, len(ids))
	for i, id := range ids {
		gamestateIDs[i] = id
	}

	q := models.GameStates.Delete(dm.Where(models.GameStates.Columns.ID.In(psql.Arg(gamestateIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gamestates: %w", err)
	}

	return nil
}
