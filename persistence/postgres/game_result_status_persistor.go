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

var _ persistence.GameResultStatusPersistor = (*PgGameResultStatusPersistor)(nil)

type PgGameResultStatusPersistor struct {
	*PgPersistor
}

func NewPgGameResultStatusPersistor(ps *PgPersistor) *PgGameResultStatusPersistor {
	return &PgGameResultStatusPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameResultStatusIntegrityViolation struct{ errIntegrityViolation }
	ErrGameResultStatusUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameResultStatusNotFound         = errors.New("gameresultstatus not found")
	errGameResultStatusIntegrity        = ErrGameResultStatusIntegrityViolation{}
	errGameResultStatusUniqueName       = ErrGameResultStatusUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameResultStatusUniqueIsoAlpha2  = ErrGameResultStatusUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameResultStatusUniqueIsoAlpha3  = ErrGameResultStatusUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameResultStatusUniqueIsoNumeric = ErrGameResultStatusUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameResultStatusPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gameresultstatus_name":
			return errGameResultStatusUniqueName
		case "uq_gameresultstatus_iso_alpha2":
			return errGameResultStatusUniqueIsoAlpha2
		case "uq_gameresultstatus_iso_alpha3":
			return errGameResultStatusUniqueIsoAlpha3
		case "uq_gameresultstatus_iso_numeric":
			return errGameResultStatusUniqueIsoNumeric
		}

		return errGameResultStatusIntegrity
	}

	return pgErr
}

func (pst *PgGameResultStatusPersistor) ListGameResultStatuses(ctx context.Context, filters dbtype.ListGameResultStatusesFilters) (dbtype.PagedResult[models.GameResultStatus], error) {
	q := psql.Select(
		sm.Columns(models.GameResultStatuses.Columns),
		sm.From(models.GameResultStatuses.Name()),
		sm.GroupBy(models.GameResultStatuses.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameResultStatuses.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters) {
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

	gameresultstatuses, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameResultStatusesRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameResultStatus]{}, fmt.Errorf("query gameresultstatuses")
	}

	result := dbtype.PagedResult[models.GameResultStatus]{
		Data: make([]models.GameResultStatus, len(gameresultstatuses)),
	}
	for i, row := range gameresultstatuses {
		result.Data[i] = row.GameResultStatus
	}

	if len(gameresultstatuses) > 0 {
		result.TotalCount = gameresultstatuses[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameResultStatusPersistor) GetGameResultStatusByID(ctx context.Context, gameresultstatusID int64) (models.GameResultStatus, error) {
	q := psql.Select(
		sm.Columns(models.GameResultStatuses.Columns),
		sm.From(models.GameResultStatuses.Name()),
		sm.Where(models.GameResultStatuses.Columns.ID.EQ(psql.Arg(gameresultstatusID))),
	)

	gameresultstatus, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResultStatus]())
	if err != nil {
		return models.GameResultStatus{}, fmt.Errorf("query gameresultstatus")
	}

	return gameresultstatus, nil
}

func (pst *PgGameResultStatusPersistor) CreateGameResultStatus(ctx context.Context, in models.GameResultStatusSetter) (models.GameResultStatus, error) {
	q := models.GameResultStatuses.Insert(&in, im.Returning(models.GameResultStatuses.Columns))

	gameresultstatus, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResultStatus]())
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

	gameresultstatus, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResultStatus]())
	if err != nil {
		return models.GameResultStatus{}, fmt.Errorf("update gameresultstatus")
	}

	return gameresultstatus, nil
}

func (pst *PgGameResultStatusPersistor) DeleteGameResultStatusByID(ctx context.Context, gameresultstatusID int64) (int64, error) {
	q := models.GameResultStatuses.Delete(dm.Where(models.GameResultStatuses.Columns.ID.EQ(psql.Arg(gameresultstatusID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gameresultstatuses: %w", err)
	}

	return nil
}
