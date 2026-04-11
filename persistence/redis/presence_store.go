package redis

import (
	"context"

	"github.com/dankobg/juicer/persistence"
	"github.com/google/uuid"
)

var _ persistence.PresencePersistor = (*RedisPresencePersistor)(nil)

type RedisPresencePersistor struct {
	*RedisPersistor
}

func NewPgPresenceStore(rs *RedisPersistor) *RedisPresencePersistor {
	return &RedisPresencePersistor{
		RedisPersistor: rs,
	}
}

func (pst *RedisPresencePersistor) SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channel string) ([]string, []string, error) {
	panic("SetPresence IMPLEMENT")
}

func (pst *RedisPresencePersistor) ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, []string, error) {
	panic("ClearPresence IMPLEMENT")
}

func (pst *RedisPresencePersistor) GetPresence(ctx context.Context, userID uuid.UUID) ([]string, error) {
	panic("GetPresence IMPLEMENT")
}

func (pst *RedisPresencePersistor) RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, error) {
	panic("RefreshPresence IMPLEMENT")
}

func (pst *RedisPresencePersistor) CountInChannel(ctx context.Context, channel string) (int64, error) {
	panic("CountInChannel IMPLEMENT")
}

func (pst *RedisPresencePersistor) LastSeen(ctx context.Context, userID uuid.UUID) (int64, error) {
	panic("LastSeen IMPLEMENT")
}

func (pst *RedisPresencePersistor) GetChannelsForUser(userID uuid.UUID) ([]string, error) {
	panic("GetChannelsForUser IMPLEMENT")
}

func (pst *RedisPresencePersistor) GetInChannel(ctx context.Context, channel string) ([]string, error) {
	panic("* IMPLEMENT")
}

// func (s *RedisPresencePersistor) SetActiveGame(ctx context.Context, in models.Game) (models.Game, error) {
// 	bb, err := json.Marshal(&in)
// 	if err != nil {
// 		return models.Game{}, fmt.Errorf("failed to set active game presence, marshal error: %w", err)
// 	}

// 	if err := s.rdb.Set(ctx, fmt.Sprintf("game.%d", in.ID), bb, 0).Err(); err != nil {
// 		return models.Game{}, fmt.Errorf("failed to set active game presence: %w", err)
// 	}

// 	return in, nil
// }

// func (s *RedisPresencePersistor) GetActiveGame(ctx context.Context, gameID int64) (models.Game, error) {
// 	val, err := s.rdb.Get(ctx, fmt.Sprintf("game.%d", gameID)).Result()
// 	if err != nil {
// 		return models.Game{}, fmt.Errorf("failed to get active game presence: %w", err)
// 	}

// 	var game models.Game
// 	if err := json.Unmarshal([]byte(val), &game); err != nil {
// 		return models.Game{}, fmt.Errorf("failed to get active game presence, unmarshal error: %w", err)
// 	}

// 	return game, nil
// }

// func (s *RedisPresencePersistor) RemoveActiveGame(ctx context.Context, gameID int64) error {
// 	if err := s.rdb.Del(ctx, fmt.Sprintf("game.%d", gameID)).Err(); err != nil {
// 		return fmt.Errorf("failed to remove active game presence: %w", err)
// 	}

// 	return nil
// }

// func (s *RedisPresencePersistor) SetPlayerGameID(ctx context.Context, clientID uuid.UUID, gameID int64) error {
// 	if err := s.rdb.Set(ctx, fmt.Sprintf("client.%s", clientID.String()), gameID, 0).Err(); err != nil {
// 		return fmt.Errorf("failed to set player game id presence: %w", err)
// 	}

// 	return nil
// }

// func (s *RedisPresencePersistor) GetPlayerGameID(ctx context.Context, clientID uuid.UUID) (int64, error) {
// 	gameID, err := s.rdb.Get(ctx, fmt.Sprintf("client.%s", clientID.String())).Result()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get player game id presence: %w", err)
// 	}

// 	n, err := strconv.Atoi(gameID)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to convert gameID to number")
// 	}

// 	return int64(n), nil
// }

// func (s *RedisPresencePersistor) DelPlayerGameID(ctx context.Context, clientID uuid.UUID) error {
// 	if err := s.rdb.Del(ctx, fmt.Sprintf("client.%s", clientID.String())).Err(); err != nil {
// 		return fmt.Errorf("failed to remove player game id presence: %w", err)
// 	}

// 	return nil
// }
