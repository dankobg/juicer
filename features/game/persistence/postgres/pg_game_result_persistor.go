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

var _ game.GameResultPersistor = (*PgGameResultPersistor)(nil)

type PgGameResultPersistor struct {
	*postgres.PgPersistor
}

func NewPgGameResultPersistor(pst *postgres.PgPersistor) *PgGameResultPersistor {
	return &PgGameResultPersistor{
		PgPersistor: pst,
	}
}

func convertGameResultPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameResultPersistor) ListGameResults(ctx context.Context, filters game.ListGameResultsFilters) (pagination.WithTotal[models.GameResult], error) {
	q := psql.Select(
		sm.Columns(models.GameResults.Columns),
		sm.From(models.GameResults.Name()),
		sm.GroupBy(models.GameResults.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameResults.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameResults.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameResults.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameResultsRow struct {
		models.GameResult
		TotalCount int64
	}

	gameResultsRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameResultsRow]())
	if err != nil {
		return pagination.WithTotal[models.GameResult]{}, fmt.Errorf("query gameresults")
	}

	gameResults := make([]models.GameResult, len(gameResultsRows))
	for i, row := range gameResultsRows {
		gameResults[i] = row.GameResult
	}

	var total int64
	if len(gameResultsRows) > 0 {
		total = gameResultsRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameResults, total)

	return out, nil
}

func (pst *PgGameResultPersistor) GetGameResultByID(ctx context.Context, gameresultID int64) (models.GameResult, error) {
	q := psql.Select(
		sm.Columns(models.GameResults.Columns),
		sm.From(models.GameResults.Name()),
		sm.Where(models.GameResults.Columns.ID.EQ(psql.Arg(gameresultID))),
	)

	gameresult, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResult]())
	if err != nil {
		return models.GameResult{}, fmt.Errorf("query gameresult")
	}

	return gameresult, nil
}

func (pst *PgGameResultPersistor) CreateGameResult(ctx context.Context, in models.GameResultSetter) (models.GameResult, error) {
	q := models.GameResults.Insert(&in, im.Returning(models.GameResults.Columns))

	gameresult, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResult]())
	if err != nil {
		return models.GameResult{}, fmt.Errorf("insert gameresult")
	}

	return gameresult, nil
}

func (pst *PgGameResultPersistor) UpdateGameResult(ctx context.Context, gameresultID int64, in models.GameResultSetter) (models.GameResult, error) {
	q := models.GameResults.Update(
		in.UpdateMod(),
		um.Where(models.GameResults.Columns.ID.EQ(psql.Arg(gameresultID))),
		um.Returning(models.GameResults.Columns),
	)

	gameresult, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResult]())
	if err != nil {
		return models.GameResult{}, fmt.Errorf("update gameresult")
	}

	return gameresult, nil
}

func (pst *PgGameResultPersistor) DeleteGameResultByID(ctx context.Context, gameresultID int64) (int64, error) {
	q := models.GameResults.Delete(dm.Where(models.GameResults.Columns.ID.EQ(psql.Arg(gameresultID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gameresult: %w", err)
	}

	return gameresultID, nil
}

func (pst *PgGameResultPersistor) BulkDeleteGameResults(ctx context.Context, ids []int64) error {
	gameresultIDs := make([]any, len(ids))
	for i, id := range ids {
		gameresultIDs[i] = id
	}

	q := models.GameResults.Delete(dm.Where(models.GameResults.Columns.ID.In(psql.Arg(gameresultIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gameresults: %w", err)
	}

	return nil
}
