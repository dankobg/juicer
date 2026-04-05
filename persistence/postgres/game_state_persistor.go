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

var _ persistence.GameStatePersistor = (*PgGameStatePersistor)(nil)

type PgGameStatePersistor struct {
	*PgPersistor
}

func NewPgGameStatePersistor(ps *PgPersistor) *PgGameStatePersistor {
	return &PgGameStatePersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameStateIntegrityViolation struct{ errIntegrityViolation }
	ErrGameStateUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameStateNotFound         = errors.New("gamestate not found")
	errGameStateIntegrity        = ErrGameStateIntegrityViolation{}
	errGameStateUniqueName       = ErrGameStateUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameStateUniqueIsoAlpha2  = ErrGameStateUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameStateUniqueIsoAlpha3  = ErrGameStateUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameStateUniqueIsoNumeric = ErrGameStateUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGameStatePgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_gamestate_name":
			return errGameStateUniqueName
		case "uq_gamestate_iso_alpha2":
			return errGameStateUniqueIsoAlpha2
		case "uq_gamestate_iso_alpha3":
			return errGameStateUniqueIsoAlpha3
		case "uq_gamestate_iso_numeric":
			return errGameStateUniqueIsoNumeric
		}

		return errGameStateIntegrity
	}

	return pgErr
}

func (pst *PgGameStatePersistor) ListGameStates(ctx context.Context, filters dbtype.ListGameStatesFilters) (dbtype.PagedResult[models.GameState], error) {
	q := psql.Select(
		sm.Columns(models.GameStates.Columns),
		sm.From(models.GameStates.Name()),
		sm.GroupBy(models.GameStates.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.GameStates.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.GameStates.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.GameStates.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListGameStatesRow struct {
		models.GameState
		TotalCount int64
	}

	gamestates, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGameStatesRow]())
	if err != nil {
		return dbtype.PagedResult[models.GameState]{}, fmt.Errorf("query gamestates")
	}

	result := dbtype.PagedResult[models.GameState]{
		Data: make([]models.GameState, len(gamestates)),
	}
	for i, row := range gamestates {
		result.Data[i] = row.GameState
	}

	if len(gamestates) > 0 {
		result.TotalCount = gamestates[0].TotalCount
	}

	return result, nil
}

func (pst *PgGameStatePersistor) GetGameStateByID(ctx context.Context, gamestateID int64) (models.GameState, error) {
	q := psql.Select(
		sm.Columns(models.GameStates.Columns),
		sm.From(models.GameStates.Name()),
		sm.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))),
	)

	gamestate, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("query gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) CreateGameState(ctx context.Context, in models.GameStateSetter) (models.GameState, error) {
	q := models.GameStates.Insert(&in, im.Returning(models.GameStates.Columns))

	gamestate, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("insert gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) UpdateGameState(ctx context.Context, gamestateID int64, in models.GameStateSetter) (models.GameState, error) {
	q := models.GameStates.Update(
		in.UpdateMod(),
		um.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))),
		um.Returning(models.GameStates.Columns),
	)

	gamestate, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.GameState]())
	if err != nil {
		return models.GameState{}, fmt.Errorf("update gamestate")
	}

	return gamestate, nil
}

func (pst *PgGameStatePersistor) DeleteGameStateByID(ctx context.Context, gamestateID int64) (int64, error) {
	q := models.GameStates.Delete(dm.Where(models.GameStates.Columns.ID.EQ(psql.Arg(gamestateID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete gamestate: %w", err)
	}

	return gamestateID, nil
}

func (pst *PgGameStatePersistor) BulkDeleteGameStates(ctx context.Context, ids []int64) error {
	gamestateIDs := make([]any, len(ids))
	for i, id := range ids {
		gamestateIDs[i] = id
	}

	q := models.GameStates.Delete(dm.Where(models.GameStates.Columns.ID.In(psql.Arg(gamestateIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete gamestates: %w", err)
	}

	return nil
}
