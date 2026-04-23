package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
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

var _ persistence.GamePersistor = (*PgGamePersistor)(nil)

type PgGamePersistor struct {
	*PgPersistor
}

func NewPgGamePersistor(ps *PgPersistor) *PgGamePersistor {
	return &PgGamePersistor{
		PgPersistor: ps,
	}
}

type (
	ErrGameIntegrityViolation struct{ errIntegrityViolation }
	ErrGameUniqueViolation    struct{ errUniqueViolation }
)

var (
	ErrGameNotFound         = errors.New("game not found")
	errGameIntegrity        = ErrGameIntegrityViolation{}
	errGameUniqueName       = ErrGameUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errGameUniqueIsoAlpha2  = ErrGameUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errGameUniqueIsoAlpha3  = ErrGameUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errGameUniqueIsoNumeric = ErrGameUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertGamePgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_game_name":
			return errGameUniqueName
		case "uq_game_iso_alpha2":
			return errGameUniqueIsoAlpha2
		case "uq_game_iso_alpha3":
			return errGameUniqueIsoAlpha3
		case "uq_game_iso_numeric":
			return errGameUniqueIsoNumeric
		}

		return errGameIntegrity
	}

	return pgErr
}

func (pst *PgGamePersistor) ListGames(ctx context.Context, filters dbtype.ListGamesFilters) (dbtype.PagedResult[models.Game], error) {
	q := psql.Select(
		sm.Columns(models.Games.Columns),
		sm.From(models.Games.Name()),
		sm.GroupBy(models.Games.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Games.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListGamesParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.ID.In(psql.Arg(ids...))))
		}
		if filters.VariantID != nil {
			variantIDs := make([]any, len(*filters.VariantID))
			for i, id := range *filters.VariantID {
				variantIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.VariantID.In(psql.Arg(variantIDs...))))
		}
		if filters.TimeKindID != nil {
			timeKindIDs := make([]any, len(*filters.TimeKindID))
			for i, id := range *filters.TimeKindID {
				timeKindIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.TimeKindID.In(psql.Arg(timeKindIDs...))))
		}
		if filters.TimeCategoryID != nil {
			timeCategoryIDs := make([]any, len(*filters.TimeCategoryID))
			for i, id := range *filters.TimeCategoryID {
				timeCategoryIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.TimeCategoryID.In(psql.Arg(timeCategoryIDs...))))
		}
		if filters.ResultID != nil {
			resultIDs := make([]any, len(*filters.ResultID))
			for i, id := range *filters.ResultID {
				resultIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.ResultID.In(psql.Arg(resultIDs...))))
		}
		if filters.ResultStatusID != nil {
			resultStatusIDs := make([]any, len(*filters.ResultStatusID))
			for i, id := range *filters.ResultStatusID {
				resultStatusIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.ResultStatusID.In(psql.Arg(resultStatusIDs...))))
		}
		if filters.StateID != nil {
			stateIDs := make([]any, len(*filters.StateID))
			for i, id := range *filters.StateID {
				stateIDs[i] = id
			}
			q.Apply(sm.Where(models.Games.Columns.StateID.In(psql.Arg(stateIDs...))))
		}
		if filters.Rated != nil {
			q.Apply(sm.Where(models.Games.Columns.Rated.EQ(psql.Arg(*filters.Rated))))
		}
		if filters.CreatedAtFrom != nil {
			q.Apply(sm.Where(models.Games.Columns.CreatedAt.GTE(psql.Arg(*filters.CreatedAtFrom))))
		}
		if filters.CreatedAtTo != nil {
			q.Apply(sm.Where(models.Games.Columns.CreatedAt.LTE(psql.Arg(*filters.CreatedAtTo))))
		}
	}

	type ListGamesRow struct {
		models.Game
		TotalCount int64
	}

	games, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGamesRow]())
	if err != nil {
		return dbtype.PagedResult[models.Game]{}, fmt.Errorf("query games")
	}

	result := dbtype.PagedResult[models.Game]{
		Data: make([]models.Game, len(games)),
	}
	for i, row := range games {
		result.Data[i] = row.Game
	}

	if len(games) > 0 {
		result.TotalCount = games[0].TotalCount
	}

	return result, nil
}

func (pst *PgGamePersistor) GetGameByID(ctx context.Context, gameID int64) (models.Game, error) {
	q := psql.Select(
		sm.Columns(models.Games.Columns),
		sm.From(models.Games.Name()),
		sm.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))),
	)

	game, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("query game")
	}

	return game, nil
}

func (pst *PgGamePersistor) CreateGame(ctx context.Context, in models.GameSetter, inMoves []models.GameMoveSetter) (models.Game, error) {
	q := models.Games.Insert(&in, im.Returning(models.Games.Columns))

	game, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("insert game")
	}

	moveSetters := make([]*models.GameMoveSetter, len(inMoves))
	for i, x := range inMoves {
		moveSetters[i] = &x
	}

	q2 := models.GameMoves.Insert(bob.ToMods(moveSetters...), im.Returning(models.GameMoves.Columns))
	if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
		return models.Game{}, fmt.Errorf("insert game moves")
	}

	return game, nil
}

func (pst *PgGamePersistor) UpdateGame(ctx context.Context, gameID int64, in models.GameSetter) (models.Game, error) {
	q := models.Games.Update(
		in.UpdateMod(),
		um.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))),
		um.Returning(models.Games.Columns),
	)

	game, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("update game")
	}

	return game, nil
}

func (pst *PgGamePersistor) DeleteGameByID(ctx context.Context, gameID int64) (int64, error) {
	q := models.Games.Delete(dm.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete game: %w", err)
	}

	return gameID, nil
}

func (pst *PgGamePersistor) BulkDeleteGames(ctx context.Context, ids []int64) error {
	gameIDs := make([]any, len(ids))
	for i, id := range ids {
		gameIDs[i] = id
	}

	q := models.Games.Delete(dm.Where(models.Games.Columns.ID.In(psql.Arg(gameIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete games: %w", err)
	}

	return nil
}

func (pst *PgGamePersistor) GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *dbtype.GameStatsForUserFilters) (dbtype.GameStats, error) {
	panic("")
}
