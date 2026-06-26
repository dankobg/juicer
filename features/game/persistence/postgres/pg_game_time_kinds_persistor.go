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

var _ game.GameTimeKindPersistor = (*PgGameTimeKindPersistor)(nil)

type PgGameTimeKindPersistor struct {
	*postgres.PgPersistor
}

func NewPgGameTimeKindPersistor(pst *postgres.PgPersistor) *PgGameTimeKindPersistor {
	return &PgGameTimeKindPersistor{
		PgPersistor: pst,
	}
}

func convertGameTimeKindPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameTimeKindPersistor) ListGameTimeKinds(ctx context.Context, filters game.ListGameTimeKindsFilters) (pagination.WithTotal[models.GameTimeKind], error) {
	q := psql.Select(
		sm.Columns(models.GameTimeKinds.Columns),
		sm.From(models.GameTimeKinds.Name()),
		sm.GroupBy(models.GameTimeKinds.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameTimeKinds.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters.ListGameTimeKindsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameTimeKinds.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameTimeKinds.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameTimeKindsRow struct {
		models.GameTimeKind
		TotalCount int64
	}

	gameTimeKindsRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameTimeKindsRow]())
	if err != nil {
		return pagination.WithTotal[models.GameTimeKind]{}, fmt.Errorf("query game time kinds")
	}

	gameTimeKinds := make([]models.GameTimeKind, len(gameTimeKindsRows))
	for i, row := range gameTimeKindsRows {
		gameTimeKinds[i] = row.GameTimeKind
	}

	var total int64
	if len(gameTimeKindsRows) > 0 {
		total = gameTimeKindsRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameTimeKinds, total)

	return out, nil
}

func (pst *PgGameTimeKindPersistor) GetGameTimeKindByID(ctx context.Context, gametimekindID int64) (models.GameTimeKind, error) {
	q := psql.Select(
		sm.Columns(models.GameTimeKinds.Columns),
		sm.From(models.GameTimeKinds.Name()),
		sm.Where(models.GameTimeKinds.Columns.ID.EQ(psql.Arg(gametimekindID))),
	)

	gametimekind, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeKind]())
	if err != nil {
		return models.GameTimeKind{}, fmt.Errorf("query gametimekind")
	}

	return gametimekind, nil
}

func (pst *PgGameTimeKindPersistor) CreateGameTimeKind(ctx context.Context, in models.GameTimeKindSetter) (models.GameTimeKind, error) {
	q := models.GameTimeKinds.Insert(&in, im.Returning(models.GameTimeKinds.Columns))

	gametimekind, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeKind]())
	if err != nil {
		return models.GameTimeKind{}, fmt.Errorf("insert gametimekind")
	}

	return gametimekind, nil
}

func (pst *PgGameTimeKindPersistor) UpdateGameTimeKind(ctx context.Context, gametimekindID int64, in models.GameTimeKindSetter) (models.GameTimeKind, error) {
	q := models.GameTimeKinds.Update(
		in.UpdateMod(),
		um.Where(models.GameTimeKinds.Columns.ID.EQ(psql.Arg(gametimekindID))),
		um.Returning(models.GameTimeKinds.Columns),
	)

	gametimekind, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeKind]())
	if err != nil {
		return models.GameTimeKind{}, fmt.Errorf("update gametimekind")
	}

	return gametimekind, nil
}

func (pst *PgGameTimeKindPersistor) DeleteGameTimeKindByID(ctx context.Context, gametimekindID int64) (int64, error) {
	q := models.GameTimeKinds.Delete(dm.Where(models.GameTimeKinds.Columns.ID.EQ(psql.Arg(gametimekindID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gametimekind: %w", err)
	}

	return gametimekindID, nil
}

func (pst *PgGameTimeKindPersistor) BulkDeleteGameTimeKinds(ctx context.Context, ids []int64) error {
	gametimekindIDs := make([]any, len(ids))
	for i, id := range ids {
		gametimekindIDs[i] = id
	}

	q := models.GameTimeKinds.Delete(dm.Where(models.GameTimeKinds.Columns.ID.In(psql.Arg(gametimekindIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gametimekinds: %w", err)
	}

	return nil
}
