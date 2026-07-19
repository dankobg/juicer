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

var _ game.GameTimeCategoryPersistor = (*PgGameTimeCategoryPersistor)(nil)

type PgGameTimeCategoryPersistor struct {
	*postgres.PgPersistor
}

func NewPgGameTimeCategoryPersistor(pst *postgres.PgPersistor) *PgGameTimeCategoryPersistor {
	return &PgGameTimeCategoryPersistor{
		PgPersistor: pst,
	}
}

func convertGameTimeCategoryPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}
		return pgErr
	}

	return err
}

func (pst *PgGameTimeCategoryPersistor) ListGameTimeCategories(ctx context.Context, filters game.ListGameTimeCategoriesFilters) (pagination.WithTotal[models.GameTimeCategory], error) {
	q := psql.Select(
		sm.Columns(models.GameTimeCategories.Columns),
		sm.From(models.GameTimeCategories.Name()),
		sm.GroupBy(models.GameTimeCategories.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.GameTimeCategories.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

	if postgres.HasAnyLogicFilters(&filters.ListGameTimeCategoriesParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameTimeCategories.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameTimeCategories.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameTimeCategoriesRow struct {
		models.GameTimeCategory
		TotalCount int64
	}

	gameTimeCategoriesRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGameTimeCategoriesRow]())
	if err != nil {
		return pagination.WithTotal[models.GameTimeCategory]{}, fmt.Errorf("query gametimecategories")
	}

	gameTimeCategories := make([]models.GameTimeCategory, len(gameTimeCategoriesRows))
	for i, row := range gameTimeCategoriesRows {
		gameTimeCategories[i] = row.GameTimeCategory
	}

	var total int64
	if len(gameTimeCategoriesRows) > 0 {
		total = gameTimeCategoriesRows[0].TotalCount
	}

	out := pagination.NewWithTotal(gameTimeCategories, total)

	return out, nil
}

func (pst *PgGameTimeCategoryPersistor) GetGameTimeCategoryByID(ctx context.Context, gametimecategoryID int64) (models.GameTimeCategory, error) {
	q := psql.Select(
		sm.Columns(models.GameTimeCategories.Columns),
		sm.From(models.GameTimeCategories.Name()),
		sm.Where(models.GameTimeCategories.Columns.ID.EQ(psql.Arg(gametimecategoryID))),
	)

	gametimecategory, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeCategory]())
	if err != nil {
		return models.GameTimeCategory{}, fmt.Errorf("query gametimecategory")
	}

	return gametimecategory, nil
}

func (pst *PgGameTimeCategoryPersistor) CreateGameTimeCategory(ctx context.Context, in models.GameTimeCategorySetter) (models.GameTimeCategory, error) {
	q := models.GameTimeCategories.Insert(&in, im.Returning(models.GameTimeCategories.Columns))

	gametimecategory, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeCategory]())
	if err != nil {
		return models.GameTimeCategory{}, fmt.Errorf("insert gametimecategory")
	}

	return gametimecategory, nil
}

func (pst *PgGameTimeCategoryPersistor) UpdateGameTimeCategory(ctx context.Context, gametimecategoryID int64, in models.GameTimeCategorySetter) (models.GameTimeCategory, error) {
	q := models.GameTimeCategories.Update(
		in.UpdateMod(),
		um.Where(models.GameTimeCategories.Columns.ID.EQ(psql.Arg(gametimecategoryID))),
		um.Returning(models.GameTimeCategories.Columns),
	)

	gametimecategory, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.GameTimeCategory]())
	if err != nil {
		return models.GameTimeCategory{}, fmt.Errorf("update gametimecategory")
	}

	return gametimecategory, nil
}

func (pst *PgGameTimeCategoryPersistor) DeleteGameTimeCategoryByID(ctx context.Context, gametimecategoryID int64) (int64, error) {
	q := models.GameTimeCategories.Delete(dm.Where(models.GameTimeCategories.Columns.ID.EQ(psql.Arg(gametimecategoryID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return 0, fmt.Errorf("delete gametimecategory: %w", err)
	}

	return gametimecategoryID, nil
}

func (pst *PgGameTimeCategoryPersistor) BulkDeleteGameTimeCategories(ctx context.Context, ids []int64) error {
	gametimecategoryIDs := make([]any, len(ids))
	for i, id := range ids {
		gametimecategoryIDs[i] = id
	}

	q := models.GameTimeCategories.Delete(dm.Where(models.GameTimeCategories.Columns.ID.In(psql.Arg(gametimecategoryIDs...))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete gametimecategories: %w", err)
	}

	return nil
}
