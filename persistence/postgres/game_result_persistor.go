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

var _ persistence.GameResultPersistor = (*PgGameResultPersistor)(nil)

type PgGameResultPersistor struct {
	*PgPersistor
}

func NewPgGameResultPersistor(ps *PgPersistor) *PgGameResultPersistor {
	return &PgGameResultPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameResultIntegrityViolation struct{ errIntegrityViolation }
	ErrGameResultUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameResultNotFound         = errors.New("gameresult not found")
	errGameResultIntegrity        = ErrGameResultIntegrityViolation{}
	errGameResultUniqueName       = ErrGameResultUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameResultUniqueIsoAlpha2  = ErrGameResultUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameResultUniqueIsoAlpha3  = ErrGameResultUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameResultUniqueIsoNumeric = ErrGameResultUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameResultPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gameresult_name":
			return errGameResultUniqueName
		case "uq_gameresult_iso_alpha2":
			return errGameResultUniqueIsoAlpha2
		case "uq_gameresult_iso_alpha3":
			return errGameResultUniqueIsoAlpha3
		case "uq_gameresult_iso_numeric":
			return errGameResultUniqueIsoNumeric
		}

		return errGameResultIntegrity
	}

	return pgErr
}

func (pst *PgGameResultPersistor) ListGameResults(ctx context.Context, filters dbtype.ListGameResultsFilters) (dbtype.PagedResult[models.GameResult], error) {
	q := psql.Select(
		sm.Columns(models.GameResults.Columns),
		sm.From(models.GameResults.Name()),
		sm.GroupBy(models.GameResults.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameResults.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters) {
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

	gameresults, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameResultsRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameResult]{}, fmt.Errorf("query gameresults")
	}

	result := dbtype.PagedResult[models.GameResult]{
		Data: make([]models.GameResult, len(gameresults)),
	}
	for i, row := range gameresults {
		result.Data[i] = row.GameResult
	}

	if len(gameresults) > 0 {
		result.TotalCount = gameresults[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameResultPersistor) GetGameResultByID(ctx context.Context, gameresultID int64) (models.GameResult, error) {
	q := psql.Select(
		sm.Columns(models.GameResults.Columns),
		sm.From(models.GameResults.Name()),
		sm.Where(models.GameResults.Columns.ID.EQ(psql.Arg(gameresultID))),
	)

	gameresult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResult]())
	if err != nil {
		return models.GameResult{}, fmt.Errorf("query gameresult")
	}

	return gameresult, nil
}

func (pst *PgGameResultPersistor) CreateGameResult(ctx context.Context, in models.GameResultSetter) (models.GameResult, error) {
	q := models.GameResults.Insert(&in, im.Returning(models.GameResults.Columns))

	gameresult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResult]())
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

	gameresult, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameResult]())
	if err != nil {
		return models.GameResult{}, fmt.Errorf("update gameresult")
	}

	return gameresult, nil
}

func (pst *PgGameResultPersistor) DeleteGameResultByID(ctx context.Context, gameresultID int64) (int64, error) {
	q := models.GameResults.Delete(dm.Where(models.GameResults.Columns.ID.EQ(psql.Arg(gameresultID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gameresults: %w", err)
	}

	return nil
}
