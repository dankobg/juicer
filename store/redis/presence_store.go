package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dankobg/juicer/db/gen/test/public/model"
	"github.com/dankobg/juicer/store"
	"github.com/google/uuid"
)

var _ store.PresenceStore = (*RedisPresenceStore)(nil)

type RedisPresenceStore struct {
	*RedisStore
}

func NewPgPresenceStore(rs *RedisStore) *RedisPresenceStore {
	return &RedisPresenceStore{
		RedisStore: rs,
	}
}

func (s *RedisPresenceStore) SetActiveGame(ctx context.Context, in model.Game) (model.Game, error) {
	bb, err := json.Marshal(&in)
	if err != nil {
		return model.Game{}, fmt.Errorf("failed to set active game presence, marshal error: %w", err)
	}
	if err := s.rdb.Set(ctx, fmt.Sprintf("game.%s", in.ID.String()), bb, 0).Err(); err != nil {
		return model.Game{}, fmt.Errorf("failed to set active game presence: %w", err)
	}
	return in, nil
}

func (s *RedisPresenceStore) GetActiveGame(ctx context.Context, gameID uuid.UUID) (model.Game, error) {
	val, err := s.rdb.Get(ctx, fmt.Sprintf("game.%s", gameID.String())).Result()
	if err != nil {
		return model.Game{}, fmt.Errorf("failed to get active game presence: %w", err)
	}
	var game model.Game
	if err := json.Unmarshal([]byte(val), &game); err != nil {
		return model.Game{}, fmt.Errorf("failed to get active game presence, unmarshal error: %w", err)
	}
	return game, nil
}

func (s *RedisPresenceStore) RemoveActiveGame(ctx context.Context, gameID uuid.UUID) error {
	if err := s.rdb.Del(ctx, fmt.Sprintf("game.%s", gameID.String())).Err(); err != nil {
		return fmt.Errorf("failed to remove active game presence: %w", err)
	}
	return nil
}

func (s *RedisPresenceStore) SetPlayerGameID(ctx context.Context, clientID uuid.UUID, gameID uuid.UUID) error {
	if err := s.rdb.Set(ctx, fmt.Sprintf("client.%s", clientID.String()), gameID.String(), 0).Err(); err != nil {
		return fmt.Errorf("failed to set player game id presence: %w", err)
	}
	return nil
}

func (s *RedisPresenceStore) GetPlayerGameID(ctx context.Context, clientID uuid.UUID) (uuid.UUID, error) {
	gameID, err := s.rdb.Get(ctx, fmt.Sprintf("client.%s", clientID.String())).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get player game id presence: %w", err)
	}
	return uuid.MustParse(gameID), nil
}

func (s *RedisPresenceStore) DelPlayerGameID(ctx context.Context, clientID uuid.UUID) error {
	if err := s.rdb.Del(ctx, fmt.Sprintf("client.%s", clientID.String())).Err(); err != nil {
		return fmt.Errorf("failed to remove player game id presence: %w", err)
	}
	return nil
}
