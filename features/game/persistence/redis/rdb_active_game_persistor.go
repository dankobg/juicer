package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dankobg/juicer/db/gen/models"
	"github.com/dankobg/juicer/features/game"
	"github.com/dankobg/juicer/gameplay"
	"github.com/dankobg/juicer/pagination"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var _ game.ActiveGamePersistor = (*RedisActiveGamePersistor)(nil)

type RedisActiveGamePersistor struct {
	rdb *redis.Client
}

func NewRedisActiveGamePersistor(rdb *redis.Client) *RedisActiveGamePersistor {
	return &RedisActiveGamePersistor{rdb: rdb}
}

func (pst *RedisActiveGamePersistor) GetActiveGameByID(ctx context.Context, gameID int64, filters game.GetActiveGameFilters) (game.GameDetails, error) {
	panic("")
}

func (pst *RedisActiveGamePersistor) ListActiveGames(ctx context.Context, filters game.ListActiveGameFilters) (pagination.WithTotal[game.GameDetails], error) {
	panic("")
}

func (pst *RedisActiveGamePersistor) ListUserActiveGames(ctx context.Context, userID uuid.UUID, filters game.ListActiveGameFilters) (pagination.WithTotal[game.GameDetails], error) {
	panic("")
}

func (pst *RedisActiveGamePersistor) IsGameActive(ctx context.Context, gameID int64) (bool, error) {
	activeGameKey := fmt.Sprintf("active-game:%d", gameID)

	n, err := pst.rdb.Exists(ctx, activeGameKey).Result()
	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func (pst *RedisActiveGamePersistor) IsUserInActiveGame(ctx context.Context, userID uuid.UUID, gameID int64) (bool, error) {
	activeGamesUserKey := "active-games:user:" + userID.String()

	ok, err := pst.rdb.SIsMember(ctx, activeGamesUserKey, gameID).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (pst *RedisActiveGamePersistor) ListActiveGameUsers(ctx context.Context, gameID int64) ([2]uuid.UUID, error) {
	panic("")
}

func (pst *RedisActiveGamePersistor) CreateActiveGame(ctx context.Context, gs *gameplay.GameState) error {
	activeGameKey := fmt.Sprintf("active-game:%d", gs.GameID)

	moves := make([]game.ActiveGameMove, len(gs.GameMoves))
	hashes := make([]game.ActiveGameHistoryHash, len(gs.Chess.HistoryHashes))

	for i, move := range gs.GameMoves {
		agm := game.ActiveGameMove{
			ID:  int64(i + 1),
			Fen: move.Fen,
		}
		if move.Uci != nil {
			agm.Uci = *move.Uci
		}

		if move.San != nil {
			agm.San = *move.San
		}

		if move.PlayedAt != nil {
			agm.PlayedAt = new(move.PlayedAt.AsTime())
		}

		moves[i] = agm
	}

	for i, hash := range gs.Chess.HistoryHashes {
		hashes[i] = game.ActiveGameHistoryHash{
			ID:   int64(i + 1),
			Hash: int64(hash),
		}
	}

	whiteSecs := int64(gs.WhiteRemainingGameTime / time.Second)
	whiteNS := int64(gs.WhiteRemainingGameTime % time.Second)
	blackSecs := int64(gs.BlackRemainingGameTime / time.Second)
	blackNS := int64(gs.BlackRemainingGameTime % time.Second)

	activeGame := game.ActiveGame{
		GameID:                 gs.GameID,
		WhiteID:                gs.White.ID.String(),
		BlackID:                gs.Black.ID.String(),
		GameVariant:            gs.GameVariant,
		GameTimeKind:           gs.GameTimeKind,
		GameTimeCategory:       gs.GameTimeCategory,
		TimeControlClockMs:     gs.GameTimeControl.GetClockMs(),
		TimeControlIncrementMs: gs.GameTimeControl.GetIncrementMs(),
		GameResult:             gs.GameResult,
		GameResultStatus:       gs.GameResultStatus,
		GameState:              gs.GameState,
		ReconnectTimeoutMs:     int32(gs.ReconnectTimeout.Milliseconds()),
		FirstMoveTimeoutMs:     int32(gs.FirstMoveTimeout.Milliseconds()),
		LastMove:               gs.LastMove,
		StartTime:              gs.StartTime,
		EndTime:                gs.EndTime,
		Rated:                  gs.Rated,
		GameMoves:              moves,
		GameHistoryHashes:      hashes,
		Version:                int32(gs.Version),
		WhiteGameRemainingSecs: int32(whiteSecs),
		WhiteGameRemainingNS:   whiteNS,
		BlackGameRemainingSecs: int32(blackSecs),
		BlackGameRemainingNS:   blackNS,
		PendingDrawOffers:      gs.PendingDrawOffers,
	}

	activeGameBytes, err := json.Marshal(activeGame)
	if err != nil {
		return err
	}

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.JSONSet(ctx, activeGameKey, "$", activeGameBytes).Err(); err != nil {
			return fmt.Errorf("CreateActiveGame jsonset: %w", err)
		}

		userIDs := [2]string{activeGame.WhiteID, activeGame.BlackID}

		for _, userID := range userIDs {
			activeGamesUserKey := "active-games:user:" + userID

			if err := p.SAdd(ctx, activeGamesUserKey, activeGame.GameID).Err(); err != nil {
				return fmt.Errorf("CreateActiveGame sadd: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("CreateActiveGame pipeline failed: %w", err)
	}

	return nil
}

func (pst *RedisActiveGamePersistor) UpdateActiveGame(ctx context.Context, gameID int64, in models.GameSetter) error {
	panic("")
}

func (pst *RedisActiveGamePersistor) DeleteActiveGameByID(ctx context.Context, gameID int64) error {
	activeGameKey := fmt.Sprintf("active-game:%d", gameID)

	var activeGame game.ActiveGame

	activeGameJSON, err := pst.rdb.JSONGet(ctx, activeGameKey, "$").Result()
	if err != nil {
		return fmt.Errorf("DeleteActiveGameByID jsonget: %w", err)
	}

	if activeGameJSON == "" {
		return redis.Nil
	}

	var wrapped []game.ActiveGame

	if err := json.Unmarshal([]byte(activeGameJSON), &wrapped); err != nil {
		return fmt.Errorf("DeleteActiveGameByID unmarshal: %w", err)
	}

	if len(wrapped) == 0 {
		return redis.Nil
	}

	activeGame = wrapped[0]

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		userIDs := [2]string{activeGame.WhiteID, activeGame.BlackID}

		for _, userID := range userIDs {
			activeGamesUserKey := "active-games:user:" + userID

			if err := p.SRem(ctx, activeGamesUserKey, gameID).Err(); err != nil {
				return fmt.Errorf("DeleteActiveGameByID srem: %w", err)
			}
		}

		if err := p.Del(ctx, activeGameKey).Err(); err != nil {
			return fmt.Errorf("DeleteActiveGameByID del: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("DeleteActiveGameByID pipeline failed: %w", err)
	}

	return nil
}
