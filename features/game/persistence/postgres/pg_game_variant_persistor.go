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

var _ game.GameVariantPersistor = (*PgGameVariantPersistor)(nil)

type PgGameVariantPersistor struct {
	*postgres.PgPersistor
}

func NewPgGameVariantPersistor(pst *postgres.PgPersistor) *PgGameVariantPersistor {
	return &PgGameVariantPersistor{
		PgPersistor: pst,
	}
}

func convertGameVariantPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameVariantPersistor) ListGameVariants(ctx context.Context, filters game.ListGameVariantsFilters) (pagination.WithTotal[models.GameVariant], error) {
	q := psql.Select(
		sm.Columns(models.GameVariants.Columns),
		sm.From(models.GameVariants.Name()),
		sm.GroupBy(models.GameVariants.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameVariants.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters.ListGameVariantsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameVariants.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameVariants.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameVariantsRow struct {
		models.GameVariant
		TotalCount int64
	}

	gameVariantsRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameVariantsRow]())
	if err != nil {
		return pagination.WithTotal[models.GameVariant]{}, fmt.Errorf("query game variants")
	}

	gameVariants := make([]models.GameVariant, len(gameVariantsRows))
	for i, row := range gameVariantsRows {
		gameVariants[i] = row.GameVariant
	}

	var total int64
	if len(gameVariantsRows) > 0 {
		total = gameVariantsRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameVariants, total)

	return out, nil
}

func (pst *PgGameVariantPersistor) GetGameVariantByID(ctx context.Context, gamevariantID int64) (models.GameVariant, error) {
	q := psql.Select(
		sm.Columns(models.GameVariants.Columns),
		sm.From(models.GameVariants.Name()),
		sm.Where(models.GameVariants.Columns.ID.EQ(psql.Arg(gamevariantID))),
	)

	gamevariant, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameVariant]())
	if err != nil {
		return models.GameVariant{}, fmt.Errorf("query gamevariant")
	}

	return gamevariant, nil
}

func (pst *PgGameVariantPersistor) CreateGameVariant(ctx context.Context, in models.GameVariantSetter) (models.GameVariant, error) {
	q := models.GameVariants.Insert(&in, im.Returning(models.GameVariants.Columns))

	gamevariant, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameVariant]())
	if err != nil {
		return models.GameVariant{}, fmt.Errorf("insert gamevariant")
	}

	return gamevariant, nil
}

func (pst *PgGameVariantPersistor) UpdateGameVariant(ctx context.Context, gamevariantID int64, in models.GameVariantSetter) (models.GameVariant, error) {
	q := models.GameVariants.Update(
		in.UpdateMod(),
		um.Where(models.GameVariants.Columns.ID.EQ(psql.Arg(gamevariantID))),
		um.Returning(models.GameVariants.Columns),
	)

	gamevariant, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameVariant]())
	if err != nil {
		return models.GameVariant{}, fmt.Errorf("update gamevariant")
	}

	return gamevariant, nil
}

func (pst *PgGameVariantPersistor) DeleteGameVariantByID(ctx context.Context, gamevariantID int64) (int64, error) {
	q := models.GameVariants.Delete(dm.Where(models.GameVariants.Columns.ID.EQ(psql.Arg(gamevariantID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gamevariant: %w", err)
	}

	return gamevariantID, nil
}

func (pst *PgGameVariantPersistor) BulkDeleteGameVariants(ctx context.Context, ids []int64) error {
	gamevariantIDs := make([]any, len(ids))
	for i, id := range ids {
		gamevariantIDs[i] = id
	}

	q := models.GameVariants.Delete(dm.Where(models.GameVariants.Columns.ID.In(psql.Arg(gamevariantIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gamevariants: %w", err)
	}

	return nil
}
