package postgres

import (
	"context"
	"errors"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
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
	"github.com/stephenafamo/bob/orm"
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

func (pst *PgGamePersistor) ListGames(ctx context.Context, filters dbtype.ListGamesFilters) (dbtype.PagedResult[dbtype.GameDetails], error) {
	q := psql.Select(
		sm.Columns(models.Games.Columns),
		sm.From(models.Games.Name()),
		sm.GroupBy(models.Games.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Games.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if filters.WithGameHashes {
		q.Apply(sm.Columns(
			psql.Raw(`(
  SELECT
    json_agg(row_to_json(game_hashes_row))
  FROM
    (
      SELECT
        game_move.id AS "id",
        game_move.game_id AS "gameID",
        game_move.hash AS "hash"
      FROM
        public.game_hash
      WHERE
        game_hash.game_id = game.id
    ) AS game_hashes_row
) AS "game_hashes"`),
		))
	}

	if hasAnyLogicFilters(&filters.ListGamesParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListGamesParamsEmbedMoves:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(game_moves_row))
  FROM
    (
      SELECT
        game_move.id AS "id",
        game_move.game_id AS "gameID",
        game_move.fen AS "fen",
        game_move.san AS "san",
        game_move.uci AS "uci",
        game_move.check AS "check",
        game_move.played_at AS "playedAt"
      FROM
        public.game_move
      WHERE
        game_move.game_id = game.id
    ) AS game_moves_row
) AS "game_moves"`),
					))
				}
			}
		}

		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.GameVariantID != nil {
			gameVariantIDs := make([]any, len(*filters.GameVariantID))
			for i, id := range *filters.GameVariantID {
				gameVariantIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameVariantID.In(psql.Arg(gameVariantIDs...))))
		}

		if filters.GameTimeKindID != nil {
			gameTimeKindIDs := make([]any, len(*filters.GameTimeKindID))
			for i, id := range *filters.GameTimeKindID {
				gameTimeKindIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameTimeKindID.In(psql.Arg(gameTimeKindIDs...))))
		}

		if filters.GameTimeCategoryID != nil {
			gameTimeCategoryIDs := make([]any, len(*filters.GameTimeCategoryID))
			for i, id := range *filters.GameTimeCategoryID {
				gameTimeCategoryIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameTimeCategoryID.In(psql.Arg(gameTimeCategoryIDs...))))
		}

		if filters.GameResultID != nil {
			gameResultIDs := make([]any, len(*filters.GameResultID))
			for i, id := range *filters.GameResultID {
				gameResultIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameResultID.In(psql.Arg(gameResultIDs...))))
		}

		if filters.GameResultStatusID != nil {
			gameResultStatusIDs := make([]any, len(*filters.GameResultStatusID))
			for i, id := range *filters.GameResultStatusID {
				gameResultStatusIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameResultStatusID.In(psql.Arg(gameResultStatusIDs...))))
		}

		if filters.GameStateID != nil {
			gameStateIDs := make([]any, len(*filters.GameStateID))
			for i, id := range *filters.GameStateID {
				gameStateIDs[i] = id
			}

			q.Apply(sm.Where(models.Games.Columns.GameStateID.In(psql.Arg(gameStateIDs...))))
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
		dbtype.GameDetails
		TotalCount int64
	}

	games, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListGamesRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.PagedResult[dbtype.GameDetails]{}, fmt.Errorf("query games")
	}

	result := dbtype.PagedResult[dbtype.GameDetails]{
		Data: make([]dbtype.GameDetails, len(games)),
	}
	for i, row := range games {
		result.Data[i] = row.GameDetails
	}

	if len(games) > 0 {
		result.TotalCount = games[0].TotalCount
	}

	return result, nil
}

func (pst *PgGamePersistor) GetGameByID(ctx context.Context, gameID int64, filters dbtype.GetGameByIDFilters) (dbtype.GameDetails, error) {
	q := psql.Select(
		sm.Columns(models.Games.Columns),
		sm.From(models.Games.Name()),
		sm.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))),
	)

	if filters.WithGameHashes {
		q.Apply(sm.Columns(
			psql.Raw(`(
  SELECT
    json_agg(row_to_json(game_hashes_row))
  FROM
    (
      SELECT
        game_history_hash.id AS "id",
        game_history_hash.game_id AS "gameID",
        game_history_hash.hash AS "hash"
      FROM
        public.game_history_hash
      WHERE
        game_history_hash.game_id = game.id
    ) AS game_hashes_row
) AS "game_history_hashes"`),
		))
	}

	if hasAnyLogicFilters(&filters.GetGameParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.GetGameParamsEmbedMoves:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(game_moves_row))
  FROM
    (
      SELECT
        game_move.id AS "id",
        game_move.game_id AS "gameID",
        game_move.fen AS "fen",
        game_move.san AS "san",
        game_move.uci AS "uci",
        game_move.check AS "check",
        game_move.played_at AS "playedAt"
      FROM
        public.game_move
      WHERE
        game_move.game_id = game.id
    ) AS game_moves_row
) AS "game_moves"`),
					))
				}
			}
		}
	}

	gameDetails, err := bob.One(ctx, pst.exec, q, scan.StructMapper[dbtype.GameDetails](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.GameDetails{}, fmt.Errorf("query game")
	}

	return gameDetails, nil
}

func (pst *PgGamePersistor) CreateGame(ctx context.Context, in models.GameSetter, inMoves []models.GameMoveSetter, inHashes []models.GameHistoryHashSetter) (models.Game, error) {
	q := models.Games.Insert(&in, im.Returning(models.Games.Columns))

	game, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("insert game: %w", err)
	}

	if len(inMoves) > 0 {
		moveSetters := make([]*models.GameMoveSetter, len(inMoves))
		for i, x := range inMoves {
			x.GameID.Set(game.ID)
			moveSetters[i] = &x
		}

		q2 := models.GameMoves.Insert(bob.ToMods(moveSetters...), im.Returning(models.GameMoves.Columns))
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("insert game moves")
		}
	}

	if len(inHashes) > 0 {
		hashSetters := make([]*models.GameHistoryHashSetter, len(inHashes))
		for i, x := range inHashes {
			x.GameID.Set(game.ID)
			hashSetters[i] = &x
		}

		q2 := models.GameHistoryHashes.Insert(bob.ToMods(hashSetters...), im.Returning(models.GameHistoryHashes.Columns))
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("insert game history hashes")
		}
	}

	return game, nil
}

func (pst *PgGamePersistor) UpdateGame(ctx context.Context, gameID int64, in models.GameSetter, newMove *models.GameMoveSetter, newHash *models.GameHistoryHashSetter) (models.Game, error) {
	q := models.Games.Update(
		in.UpdateMod(),
		um.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))),
		um.Returning(models.Games.Columns),
	)

	game, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("update game")
	}

	if newMove != nil {
		newMove.GameID.Set(gameID)

		q2 := models.GameMoves.Insert(newMove)
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("update game move")
		}
	}

	if newHash != nil {
		newHash.GameID.Set(gameID)

		q2 := models.GameHistoryHashes.Insert(newHash)
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("update game history hash")
		}
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
