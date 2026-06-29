package postgres

import (
	"context"
	"errors"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/pagination"
	"github.com/dankobg/juicer/postgres"
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

var _ game.GamePersistor = (*PgGamePersistor)(nil)

type PgGamePersistor struct {
	*postgres.PgPersistor
}

func NewPgGamePersistor(pst *postgres.PgPersistor) *PgGamePersistor {
	return &PgGamePersistor{
		PgPersistor: pst,
	}
}

func convertGamePgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		_ = pgerrcode.SuccessfulCompletion
		// 	switch pgErr.Code {
		// 	case pgerrcode.UniqueViolation:
		// 		return game.ErrGameAlreadyExists
		// 	}

		return pgErr
	}

	return err
}

func (pst *PgGamePersistor) ListGames(ctx context.Context, filters game.ListGamesFilters) (pagination.WithTotal[game.GameDetails], error) {
	q := psql.Select(
		sm.Columns(models.Games.Columns),
		sm.From(models.Games.Name()),
		sm.GroupBy(models.Games.Columns.ID),
	)
	postgres.AddOrderBy(&q, filters.Sort, models.Games.Columns.Names())
	postgres.AddPagination(&q, filters.Page, filters.PageSize)

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

	if postgres.HasAnyLogicFilters(&filters.ListGamesParams) {
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
		game.GameDetails
		TotalCount int64
	}

	gamesRows, err := bob.All(ctx, pst.Exec, q, scan.StructMapper[ListGamesRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return pagination.WithTotal[game.GameDetails]{}, fmt.Errorf("query games")
	}

	games := make([]game.GameDetails, len(gamesRows))
	for i, row := range gamesRows {
		games[i] = row.GameDetails
	}

	var total int64
	if len(gamesRows) > 0 {
		total = gamesRows[0].TotalCount
	}

	out := pagination.NewWithTotal(games, total)

	return out, nil
}

func (pst *PgGamePersistor) GetGameByID(ctx context.Context, gameID int64, filters game.GetGameByIDFilters) (game.GameDetails, error) {
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

	if postgres.HasAnyLogicFilters(&filters.GetGameParams) {
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

	gameDetails, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[game.GameDetails](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return game.GameDetails{}, fmt.Errorf("query game")
	}

	return gameDetails, nil
}

func (pst *PgGamePersistor) CreateGame(ctx context.Context, in models.GameSetter, inMoves []models.GameMoveSetter, inHashes []models.GameHistoryHashSetter) (models.Game, error) {
	q := models.Games.Insert(&in, im.Returning(models.Games.Columns))

	game, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.Game]())
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
		if _, err := bob.Exec(ctx, pst.Exec, q2); err != nil {
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
		if _, err := bob.Exec(ctx, pst.Exec, q2); err != nil {
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

	game, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[models.Game]())
	if err != nil {
		return models.Game{}, fmt.Errorf("update game")
	}

	if newMove != nil {
		newMove.GameID.Set(gameID)

		q2 := models.GameMoves.Insert(newMove)
		if _, err := bob.Exec(ctx, pst.Exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("update game move")
		}
	}

	if newHash != nil {
		newHash.GameID.Set(gameID)

		q2 := models.GameHistoryHashes.Insert(newHash)
		if _, err := bob.Exec(ctx, pst.Exec, q2); err != nil {
			return models.Game{}, fmt.Errorf("update game history hash")
		}
	}

	return game, nil
}

func (pst *PgGamePersistor) DeleteGameByID(ctx context.Context, gameID int64) (int64, error) {
	q := models.Games.Delete(dm.Where(models.Games.Columns.ID.EQ(psql.Arg(gameID))))
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
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
	if _, err := bob.Exec(ctx, pst.Exec, q); err != nil {
		return fmt.Errorf("delete games: %w", err)
	}

	return nil
}

func (pst *PgGamePersistor) GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *game.GameStatsForUserFilters) (game.GameStats, error) {
	raw := `WITH stats AS (
    SELECT
        gtc.name AS category,
        CASE
            WHEN gr.name = 'interrupted' THEN 'interrupted'
            WHEN gr.name = 'draw' THEN 'draw'
            WHEN (
                (gr.name = 'white-won' AND g.white_id = ?)
                OR
                (gr.name = 'black-won' AND g.black_id = ?)
            ) THEN 'win'
            ELSE 'loss'
        END AS result
    FROM game g
    JOIN game_time_category gtc
        ON g.game_time_category_id = gtc.id
    JOIN game_state gs
        ON g.game_state_id = gs.id
    LEFT JOIN game_result gr
        ON g.game_result_id = gr.id
    WHERE
        (g.white_id = ? OR g.black_id = ?)
        AND gs.name = 'finished'
),
counts AS (
    SELECT
        category,
        COUNT(*) FILTER (WHERE result = 'win') AS win,
        COUNT(*) FILTER (WHERE result = 'loss') AS loss,
        COUNT(*) FILTER (WHERE result = 'draw') AS draw,
        COUNT(*) FILTER (WHERE result = 'interrupted') AS interrupted,
        COUNT(*) AS total
    FROM stats
    GROUP BY category
),
all_counts AS (
    SELECT
        COUNT(*) FILTER (WHERE result = 'win') AS win,
        COUNT(*) FILTER (WHERE result = 'loss') AS loss,
        COUNT(*) FILTER (WHERE result = 'draw') AS draw,
        COUNT(*) FILTER (WHERE result = 'interrupted') AS interrupted,
        COUNT(*) AS total
    FROM stats
)
SELECT json_build_object(
    'all',
    (
        SELECT row_to_json(a)
        FROM all_counts a
    ),
    'hyperbullet',
    COALESCE(
        (SELECT row_to_json(c) FROM counts c WHERE category = 'hyperbullet'),
        json_build_object('win',0,'loss',0,'draw',0,'interrupted',0,'total',0)
    ),
    'bullet',
    COALESCE(
        (SELECT row_to_json(c) FROM counts c WHERE category = 'bullet'),
        json_build_object('win',0,'loss',0,'draw',0,'interrupted',0,'total',0)
    ),
    'blitz',
    COALESCE(
        (SELECT row_to_json(c) FROM counts c WHERE category = 'blitz'),
        json_build_object('win',0,'loss',0,'draw',0,'interrupted',0,'total',0)
    ),
    'rapid',
    COALESCE(
        (SELECT row_to_json(c) FROM counts c WHERE category = 'rapid'),
        json_build_object('win',0,'loss',0,'draw',0,'interrupted',0,'total',0)
    ),
    'classical',
    COALESCE(
        (SELECT row_to_json(c) FROM counts c WHERE category = 'classical'),
        json_build_object('win',0,'loss',0,'draw',0,'interrupted',0,'total',0)
    )
) AS stats;`

	q := psql.RawQuery(raw, userID.String(), userID.String(), userID.String(), userID.String())
	gameStats, err := bob.One(ctx, pst.Exec, q, scan.StructMapper[game.GameStatsJSON](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return game.GameStats{}, fmt.Errorf("query game stats")
	}

	return gameStats.Stats.Val, nil
}
