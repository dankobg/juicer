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

var _ game.GameResultStatusPersistor = (*PgGameResultStatusPersistor)(nil)

type PgGameResultStatusPersistor struct {
	*postgres.PgPersistor
}

func NewPgGameResultStatusPersistor(pst *postgres.PgPersistor) *PgGameResultStatusPersistor {
	return &PgGameResultStatusPersistor{
		PgPersistor: pst,
	}
}

func convertGameResultStatusPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameResultStatusPersistor) ListGameResultStatuses(ctx context.Context, filters game.ListGameResultStatusesFilters) (pagination.WithTotal[models.GameResultStatus], error) {
	q := psql.Select(
		sm.Columns(models.GameResultStatuses.Columns),
		sm.From(models.GameResultStatuses.Name()),
		sm.GroupBy(models.GameResultStatuses.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameResultStatuses.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameResultStatuses.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameResultStatuses.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameResultStatusesRow struct {
		models.GameResultStatus
		TotalCount int64
	}

	gameResultStatusesRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameResultStatusesRow]())
	if err != nil {
		return pagination.WithTotal[models.GameResultStatus]{}, fmt.Errorf("query gameresultstatuses")
	}

	gameResultStatuses := make([]models.GameResultStatus, len(gameResultStatusesRows))
	for i, row := range gameResultStatusesRows {
		gameResultStatuses[i] = row.GameResultStatus
	}

	var total int64
	if len(gameResultStatusesRows) > 0 {
		total = gameResultStatusesRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameResultStatuses, total)

	return out, nil
}

func (pst *PgGameResultStatusPersistor) GetGameResultStatusByID(ctx context.Context, gameresultstatusID int64) (models.GameResultStatus, error) {
	q := psql.Select(
		sm.Columns(models.GameResultStatuses.Columns),
		sm.From(models.GameResultStatuses.Name()),
		sm.Where(models.GameResultStatuses.Columns.ID.EQ(psql.Arg(gameresultstatusID))),
	)

	gameresultstatus, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResultStatus]())
	if err != nil {
		return models.GameResultStatus{}, fmt.Errorf("query gameresultstatus")
	}

	return gameresultstatus, nil
}

func (pst *PgGameResultStatusPersistor) CreateGameResultStatus(ctx context.Context, in models.GameResultStatusSetter) (models.GameResultStatus, error) {
	q := models.GameResultStatuses.Insert(&in, im.Returning(models.GameResultStatuses.Columns))

	gameresultstatus, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResultStatus]())
	if err != nil {
		return models.GameResultStatus{}, fmt.Errorf("insert gameresultstatus")
	}

	return gameresultstatus, nil
}

func (pst *PgGameResultStatusPersistor) UpdateGameResultStatus(ctx context.Context, gameresultstatusID int64, in models.GameResultStatusSetter) (models.GameResultStatus, error) {
	q := models.GameResultStatuses.Update(
		in.UpdateMod(),
		um.Where(models.GameResultStatuses.Columns.ID.EQ(psql.Arg(gameresultstatusID))),
		um.Returning(models.GameResultStatuses.Columns),
	)

	gameresultstatus, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameResultStatus]())
	if err != nil {
		return models.GameResultStatus{}, fmt.Errorf("update gameresultstatus")
	}

	return gameresultstatus, nil
}

func (pst *PgGameResultStatusPersistor) DeleteGameResultStatusByID(ctx context.Context, gameresultstatusID int64) (int64, error) {
	q := models.GameResultStatuses.Delete(dm.Where(models.GameResultStatuses.Columns.ID.EQ(psql.Arg(gameresultstatusID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gameresultstatus: %w", err)
	}

	return gameresultstatusID, nil
}

func (pst *PgGameResultStatusPersistor) BulkDeleteGameResultStatuses(ctx context.Context, ids []int64) error {
	gameresultstatusIDs := make([]any, len(ids))
	for i, id := range ids {
		gameresultstatusIDs[i] = id
	}

	q := models.GameResultStatuses.Delete(dm.Where(models.GameResultStatuses.Columns.ID.In(psql.Arg(gameresultstatusIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gameresultstatuses: %w", err)
	}

	return nil
}
