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

var _ persistence.GameVariantPersistor = (*PgGameVariantPersistor)(nil)

type PgGameVariantPersistor struct {
	*PgPersistor
}

func NewPgGameVariantPersistor(ps *PgPersistor) *PgGameVariantPersistor {
	return &PgGameVariantPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameVariantIntegrityViolation struct{ errIntegrityViolation }
	ErrGameVariantUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameVariantNotFound         = errors.New("gamevariant not found")
	errGameVariantIntegrity        = ErrGameVariantIntegrityViolation{}
	errGameVariantUniqueName       = ErrGameVariantUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameVariantUniqueIsoAlpha2  = ErrGameVariantUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameVariantUniqueIsoAlpha3  = ErrGameVariantUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameVariantUniqueIsoNumeric = ErrGameVariantUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameVariantPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gamevariant_name":
			return errGameVariantUniqueName
		case "uq_gamevariant_iso_alpha2":
			return errGameVariantUniqueIsoAlpha2
		case "uq_gamevariant_iso_alpha3":
			return errGameVariantUniqueIsoAlpha3
		case "uq_gamevariant_iso_numeric":
			return errGameVariantUniqueIsoNumeric
		}

		return errGameVariantIntegrity
	}

	return pgErr
}

func (pst *PgGameVariantPersistor) ListGameVariants(ctx context.Context, filters dbtype.ListGameVariantsFilters) (dbtype.PagedResult[models.GameVariant], error) {
	q := psql.Select(
		sm.Columns(models.GameVariants.Columns),
		sm.From(models.GameVariants.Name()),
		sm.GroupBy(models.GameVariants.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameVariants.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListGameVariantsParams) {
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

	gamevariants, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameVariantsRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameVariant]{}, fmt.Errorf("query gamevariants")
	}

	result := dbtype.PagedResult[models.GameVariant]{
		Data: make([]models.GameVariant, len(gamevariants)),
	}
	for i, row := range gamevariants {
		result.Data[i] = row.GameVariant
	}

	if len(gamevariants) > 0 {
		result.TotalCount = gamevariants[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameVariantPersistor) GetGameVariantByID(ctx context.Context, gamevariantID int64) (models.GameVariant, error) {
	q := psql.Select(
		sm.Columns(models.GameVariants.Columns),
		sm.From(models.GameVariants.Name()),
		sm.Where(models.GameVariants.Columns.ID.EQ(psql.Arg(gamevariantID))),
	)

	gamevariant, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameVariant]())
	if err != nil {
		return models.GameVariant{}, fmt.Errorf("query gamevariant")
	}

	return gamevariant, nil
}

func (pst *PgGameVariantPersistor) CreateGameVariant(ctx context.Context, in models.GameVariantSetter) (models.GameVariant, error) {
	q := models.GameVariants.Insert(&in, im.Returning(models.GameVariants.Columns))

	gamevariant, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameVariant]())
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

	gamevariant, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameVariant]())
	if err != nil {
		return models.GameVariant{}, fmt.Errorf("update gamevariant")
	}

	return gamevariant, nil
}

func (pst *PgGameVariantPersistor) DeleteGameVariantByID(ctx context.Context, gamevariantID int64) (int64, error) {
	q := models.GameVariants.Delete(dm.Where(models.GameVariants.Columns.ID.EQ(psql.Arg(gamevariantID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gamevariants: %w", err)
	}

	return nil
}
