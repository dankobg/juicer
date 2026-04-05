package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

var _ persistence.GameTimeCategoryPersistor = (*PgGameTimeCategoryPersistor)(nil)

type PgGameTimeCategoryPersistor struct {
	*PgPersistor
}

func NewPgGameTimeCategoryPersistor(ps *PgPersistor) *PgGameTimeCategoryPersistor {
	return &PgGameTimeCategoryPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameTimeCategoryIntegrityViolation struct{ errIntegrityViolation }
	ErrGameTimeCategoryUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameTimeCategoryNotFound         = errors.New("gametimecategory not found")
	errGameTimeCategoryIntegrity        = ErrGameTimeCategoryIntegrityViolation{}
	errGameTimeCategoryUniqueName       = ErrGameTimeCategoryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameTimeCategoryUniqueIsoAlpha2  = ErrGameTimeCategoryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameTimeCategoryUniqueIsoAlpha3  = ErrGameTimeCategoryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameTimeCategoryUniqueIsoNumeric = ErrGameTimeCategoryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameTimeCategoryPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gametimecategory_name":
			return errGameTimeCategoryUniqueName
		case "uq_gametimecategory_iso_alpha2":
			return errGameTimeCategoryUniqueIsoAlpha2
		case "uq_gametimecategory_iso_alpha3":
			return errGameTimeCategoryUniqueIsoAlpha3
		case "uq_gametimecategory_iso_numeric":
			return errGameTimeCategoryUniqueIsoNumeric
		}

		return errGameTimeCategoryIntegrity
	}

	return pgErr
}

func (pst *PgGameTimeCategoryPersistor) ListGameTimeCategories(ctx context.Context, filters dbtype.ListGameTimeCategoriesFilters) (dbtype.PagedResult[models.GameTimeCategory], error) {
	q := psql.Select(
		sm.Columns(models.GameTimeCategories.Columns),
		sm.From(models.GameTimeCategories.Name()),
		sm.GroupBy(models.GameTimeCategories.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameTimeCategories.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListGameTimeCategoriesParams) {
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

	gametimecategories, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameTimeCategoriesRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameTimeCategory]{}, fmt.Errorf("query gametimecategories")
	}

	result := dbtype.PagedResult[models.GameTimeCategory]{
		Data: make([]models.GameTimeCategory, len(gametimecategories)),
	}
	for i, row := range gametimecategories {
		result.Data[i] = row.GameTimeCategory
	}

	if len(gametimecategories) > 0 {
		result.TotalCount = gametimecategories[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameTimeCategoryPersistor) GetGameTimeCategoryByID(ctx context.Context, gametimecategoryID int64) (models.GameTimeCategory, error) {
	q := psql.Select(
		sm.Columns(models.GameTimeCategories.Columns),
		sm.From(models.GameTimeCategories.Name()),
		sm.Where(models.GameTimeCategories.Columns.ID.EQ(psql.Arg(gametimecategoryID))),
	)

	gametimecategory, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeCategory]())
	if err != nil {
		return models.GameTimeCategory{}, fmt.Errorf("query gametimecategory")
	}

	return gametimecategory, nil
}

func (pst *PgGameTimeCategoryPersistor) CreateGameTimeCategory(ctx context.Context, in models.GameTimeCategorySetter) (models.GameTimeCategory, error) {
	q := models.GameTimeCategories.Insert(&in, im.Returning(models.GameTimeCategories.Columns))

	gametimecategory, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeCategory]())
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

	gametimecategory, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeCategory]())
	if err != nil {
		return models.GameTimeCategory{}, fmt.Errorf("update gametimecategory")
	}

	return gametimecategory, nil
}

func (pst *PgGameTimeCategoryPersistor) DeleteGameTimeCategoryByID(ctx context.Context, gametimecategoryID int64) (int64, error) {
	q := models.GameTimeCategories.Delete(dm.Where(models.GameTimeCategories.Columns.ID.EQ(psql.Arg(gametimecategoryID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gametimecategories: %w", err)
	}

	return nil
}
