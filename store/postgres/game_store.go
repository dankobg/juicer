package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/db/gen/test/public/model"
	t "github.com/dankobg/juicer/db/gen/test/public/table"
	"github.com/dankobg/juicer/store"
	"github.com/google/uuid"

	p "github.com/go-jet/jet/v2/postgres"
)

var _ store.GameStore = (*PgGameStore)(nil)

type PgGameStore struct {
	*PgStore
}

func NewPgGameStore(ps *PgStore) *PgGameStore {
	return &PgGameStore{
		PgStore: ps,
	}
}

func (s *PgGameStore) Create(ctx context.Context, in dto.GameChangeset, inMoves []dto.GameMoveChangeset) (model.Game, error) {
	_, m := in.ToModel()

	q := t.Game.INSERT(t.Game.ID, t.Game.MutableColumns.Except(t.Game.CreatedAt, t.Game.UpdatedAt)).
		MODEL(m).
		RETURNING(t.Game.AllColumns)
	var gameDest model.Game
	var gameMovesDest []model.GameMove
	gms := make([]model.GameMove, 0, len(inMoves))
	for _, x := range inMoves {
		_, m := x.ToModel()
		gms = append(gms, m)
	}

	q2 := t.GameMove.INSERT(t.GameMove.MutableColumns).
		MODELS(gms).
		RETURNING(t.GameMove.AllColumns)

	if err := WithTx(context.Background(), s.db, func(tx *sql.Tx) error {
		if err := q.QueryContext(ctx, tx, &gameDest); err != nil {
			return err
		}
		if len(gms) > 0 {
			if err := q2.QueryContext(ctx, tx, &gameMovesDest); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return gameDest, fmt.Errorf("failed to create a game: %w", err)
	}
	return gameDest, nil
}

func (s *PgGameStore) Update(ctx context.Context, gameID uuid.UUID, in dto.GameChangeset) (model.Game, error) {
	cols, m := in.ToModel()

	q := t.Game.UPDATE(cols).
		MODEL(m).
		WHERE(t.Game.ID.EQ(p.UUID(gameID))).
		RETURNING(t.Game.AllColumns)

	var dest model.Game
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update a game: %w", err)
	}
	return dest, nil
}

func (s *PgGameStore) Delete(ctx context.Context, gameID uuid.UUID) (uuid.UUID, error) {
	q := t.Game.DELETE().WHERE(t.Game.ID.EQ(p.UUID(gameID)))

	if _, err := q.ExecContext(ctx, s.db); err != nil {
		return uuid.Nil, fmt.Errorf("failed to delete a game: %w", err)
	}
	return gameID, nil
}

func (s *PgGameStore) Get(ctx context.Context, gameID uuid.UUID) (model.Game, error) {
	q := p.SELECT(t.Game.AllColumns).
		FROM(t.Game).
		WHERE(t.Game.ID.EQ(p.UUID(gameID)))

	var dest model.Game
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to get a game: %w", err)
	}
	return dest, nil
}

func (s *PgGameStore) List(ctx context.Context) ([]model.Game, error) {
	q := p.SELECT(t.Game.AllColumns).
		FROM(t.Game)

	var dest []model.Game
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}
	return dest, nil
}

func (s *PgGameStore) InsertMove(ctx context.Context, in dto.GameMoveChangeset) (model.GameMove, error) {
	_, m := in.ToModel()

	q := t.Game.INSERT(t.GameMove.ID, t.GameMove.MutableColumns).
		MODEL(m).
		RETURNING(t.GameMove.AllColumns)

	var dest model.GameMove
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to insert a game move: %w", err)
	}
	return dest, nil
}

func (s *PgGameStore) GetGameStatsForUser(ctx context.Context, userID uuid.UUID, filters *store.GameStatsForUserFilters) (api.GameStats, error) {
	var condition string
	rawArgs := p.RawArgs{"$user_id": userID.String()}
	if filters != nil {
		if filters.From != nil {
			condition = "AND g.start_time >= $from_time"
			rawArgs["$from_time"] = *filters.From
		}
		if filters.To != nil {
			condition = "AND g.start_time <= $to_time"
			rawArgs["$to_time"] = *filters.To
		}
	}
	raw := fmt.Sprintf(`WITH time_categories AS (
  SELECT name FROM game_time_category
),
user_games AS (
  SELECT
    gtc.name AS time_category,
    COUNT(*) FILTER (
      WHERE (u.id = g.white_id AND gr.name = 'white-won')
         OR (u.id = g.black_id AND gr.name = 'black-won')
    ) AS win,
    COUNT(*) FILTER (
      WHERE (u.id = g.white_id AND gr.name = 'black-won')
         OR (u.id = g.black_id AND gr.name = 'white-won')
    ) AS loss,
    COUNT(*) FILTER (WHERE gr.name = 'draw') AS draw,
    COUNT(*) FILTER (WHERE gr.name = 'interrupted') AS interrupted,
    COUNT(*) AS total
  FROM game g
  JOIN "user" u ON u.id = g.white_id OR u.id = g.black_id
  JOIN game_time_category gtc ON g.time_category_id = gtc.id
  JOIN game_result gr ON g.result_id = gr.id
  WHERE u.id = $user_id %s
  GROUP BY gtc.name
)
SELECT json_object_agg(key, val) as game_stats
FROM (
  -- per time control stats
  SELECT tc.name AS key, json_build_object(
    'win', COALESCE(ug.win, 0),
    'loss', COALESCE(ug.loss, 0),
    'draw', COALESCE(ug.draw, 0),
    'interrupted', COALESCE(ug.interrupted, 0),
    'total', COALESCE(ug.total, 0)
  ) AS val
  FROM time_categories tc
  LEFT JOIN user_games ug ON tc.name = ug.time_category
  UNION ALL
  -- all time controls together stats
  SELECT 'all', json_build_object(
    'win', COALESCE(SUM(ug.win), 0),
    'loss', COALESCE(SUM(ug.loss), 0),
    'draw', COALESCE(SUM(ug.draw), 0),
    'interrupted', COALESCE(SUM(ug.interrupted), 0),
    'total', COALESCE(SUM(ug.total), 0)
  )
  FROM user_games ug
) AS combined`, condition)
	q := p.RawStatement(raw, rawArgs)
	var dest struct {
		GameStats []byte
	}
	var gameStats api.GameStats
	if err := q.QueryContext(ctx, s.db, &dest); err != nil {
		return gameStats, fmt.Errorf("failed to get game stats: %w", err)
	}
	if err := json.Unmarshal(dest.GameStats, &gameStats); err != nil {
		return gameStats, fmt.Errorf("failed to unmarshal db json: %w", err)
	}
	return gameStats, nil
}
