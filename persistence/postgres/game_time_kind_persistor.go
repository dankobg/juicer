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

var _ persistence.GameTimeKindPersistor = (*PgGameTimeKindPersistor)(nil)

type PgGameTimeKindPersistor struct {
	*PgPersistor
}

func NewPgGameTimeKindPersistor(ps *PgPersistor) *PgGameTimeKindPersistor {
	return &PgGameTimeKindPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameTimeKindIntegrityViolation struct{ errIntegrityViolation }
	ErrGameTimeKindUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameTimeKindNotFound         = errors.New("gametimekind not found")
	errGameTimeKindIntegrity        = ErrGameTimeKindIntegrityViolation{}
	errGameTimeKindUniqueName       = ErrGameTimeKindUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameTimeKindUniqueIsoAlpha2  = ErrGameTimeKindUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameTimeKindUniqueIsoAlpha3  = ErrGameTimeKindUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameTimeKindUniqueIsoNumeric = ErrGameTimeKindUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameTimeKindPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gametimekind_name":
			return errGameTimeKindUniqueName
		case "uq_gametimekind_iso_alpha2":
			return errGameTimeKindUniqueIsoAlpha2
		case "uq_gametimekind_iso_alpha3":
			return errGameTimeKindUniqueIsoAlpha3
		case "uq_gametimekind_iso_numeric":
			return errGameTimeKindUniqueIsoNumeric
		}

		return errGameTimeKindIntegrity
	}

	return pgErr
}

func (pst *PgGameTimeKindPersistor) ListGameTimeKinds(ctx context.Context, filters dbtype.ListGameTimeKindsFilters) (dbtype.PagedResult[models.GameTimeKind], error) {
	q := psql.Select(
		sm.Columns(models.GameTimeKinds.Columns),
		sm.From(models.GameTimeKinds.Name()),
		sm.GroupBy(models.GameTimeKinds.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameTimeKinds.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListGameTimeKindsParams) {
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

	gametimekinds, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameTimeKindsRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameTimeKind]{}, fmt.Errorf("query gametimekinds")
	}

	result := dbtype.PagedResult[models.GameTimeKind]{
		Data: make([]models.GameTimeKind, len(gametimekinds)),
	}
	for i, row := range gametimekinds {
		result.Data[i] = row.GameTimeKind
	}

	if len(gametimekinds) > 0 {
		result.TotalCount = gametimekinds[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameTimeKindPersistor) GetGameTimeKindByID(ctx context.Context, gametimekindID int64) (models.GameTimeKind, error) {
	q := psql.Select(
		sm.Columns(models.GameTimeKinds.Columns),
		sm.From(models.GameTimeKinds.Name()),
		sm.Where(models.GameTimeKinds.Columns.ID.EQ(psql.Arg(gametimekindID))),
	)

	gametimekind, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeKind]())
	if err != nil {
		return models.GameTimeKind{}, fmt.Errorf("query gametimekind")
	}

	return gametimekind, nil
}

func (pst *PgGameTimeKindPersistor) CreateGameTimeKind(ctx context.Context, in models.GameTimeKindSetter) (models.GameTimeKind, error) {
	q := models.GameTimeKinds.Insert(&in, im.Returning(models.GameTimeKinds.Columns))

	gametimekind, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeKind]())
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

	gametimekind, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameTimeKind]())
	if err != nil {
		return models.GameTimeKind{}, fmt.Errorf("update gametimekind")
	}

	return gametimekind, nil
}

func (pst *PgGameTimeKindPersistor) DeleteGameTimeKindByID(ctx context.Context, gametimekindID int64) (int64, error) {
	q := models.GameTimeKinds.Delete(dm.Where(models.GameTimeKinds.Columns.ID.EQ(psql.Arg(gametimekindID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gametimekinds: %w", err)
	}

	return nil
}
