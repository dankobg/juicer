package redis

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dankobg/juicer/persistence"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	authStr := "auth"
	if guest {
		authStr = "guest"
	}

	var (
		presenceChannelKey     = "presence:channel:" + channel
		presenceUsersKey       = "presence:users"
		presenceUserKey        = "presence:user:" + userID.String()
		presenceActiveGamesKey = "presence:active-games:" + userID.String()
		userKey                = userID.String() + "#" + connID.String() + "#" + username + "#" + authStr
		connChannelKey         = connID.String() + "#" + channel
		userkeyChannel         = userKey + "#" + channel
	)

	now := time.Now()
	expireDur := time.Minute * 2
	expiry := now.Add(expireDur)

	beforeUserPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch before presences")
	}

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZAdd(ctx, presenceChannelKey, redis.Z{
			Score:  float64(expiry.UnixMilli()),
			Member: userKey,
		}).Err(); err != nil {
			return fmt.Errorf("setpresence zadd presenceChannelKey: %w", err)
		}

		if err := p.Expire(ctx, presenceChannelKey, expireDur).Err(); err != nil {
			return fmt.Errorf("setpresence expire presenceChannelKey: %w", err)
		}

		if err := p.ZAdd(ctx, presenceUserKey, redis.Z{
			Score:  float64(expiry.UnixMilli()),
			Member: connChannelKey,
		}).Err(); err != nil {
			return fmt.Errorf("setpresence zadd presenceUserKey: %w", err)
		}

		if err := p.Expire(ctx, presenceUserKey, expireDur).Err(); err != nil {
			return fmt.Errorf("setpresence expire presenceUserKey: %w", err)
		}

		if err := p.ZAdd(ctx, presenceUsersKey, redis.Z{
			Score:  float64(expiry.UnixMilli()),
			Member: userkeyChannel,
		}).Err(); err != nil {
			return fmt.Errorf("setpresence zadd userpresences: %w", err)
		}

		_ = []any{presenceActiveGamesKey}

		// @TODO: HANDLE ACTIVE GAMES

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("setpresence pipeline: %w", err)
	}

	afterUserPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch after presences")
	}

	slices.Sort(beforeUserPresences)
	slices.Sort(afterUserPresences)

	return beforeUserPresences, afterUserPresences, nil
}

func (pst *RedisPresencePersistor) ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, []string, error) {
	authStr := "auth"
	if guest {
		authStr = "guest"
	}

	var (
		presenceUsersKey       = "presence:users"
		presenceUserKey        = "presence:user:" + userID.String()
		presenceActiveGamesKey = "presence:active-games:" + userID.String()
		userKey                = userID.String() + "#" + connID.String() + "#" + username + "#" + authStr
		presenceLastSeenKey    = "presence:last-seen" + userID.String()
	)

	now := time.Now()
	expireDur := time.Minute * 2
	expiry := now.Add(expireDur)

	beforePresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fetch before presences")
	}

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, connChannelKey := range beforePresences {
			parts := strings.Split(connChannelKey, "#")
			cid, channel := parts[0], parts[1]
			if connID.String() == cid {
				var (
					presenceChannelKey = "presence:channel:" + channel
					userkeyChannel     = userKey + "#" + channel
				)

				if err := p.ZRem(ctx, presenceChannelKey, userKey).Err(); err != nil {
					return fmt.Errorf("clearpresence zrem presenceChannelKey: %w", err)
				}

				if err := p.ZRem(ctx, presenceUserKey, connChannelKey).Err(); err != nil {
					return fmt.Errorf("clearpresence zrem presenceUserKey: %w", err)
				}

				if err := p.ZRem(ctx, presenceUsersKey, userkeyChannel).Err(); err != nil {
					return fmt.Errorf("clearpresence zrem presenceUserKey: %w", err)
				}

				if err := p.ZAdd(ctx, presenceLastSeenKey, redis.Z{
					Score:  float64(now.UnixMilli()),
					Member: userID.String(),
				}).Err(); err != nil {
					return fmt.Errorf("clearpresence zadd presenceLastSeenKey: %w", err)
				}

			}
		}

		_ = []any{presenceActiveGamesKey, expiry}

		// @TODO: HANDLE ACTIVE GAMES

		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("refreshpresence pipeline: %w", err)
	}

	afterPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fetch after presences")
	}

	slices.Sort(beforePresences)
	slices.Sort(afterPresences)

	return beforePresences, afterPresences, nil, nil
}

func (pst *RedisPresencePersistor) GetPresence(ctx context.Context, userID uuid.UUID) ([]string, error) {
	panic("GetPresence IMPLEMENT")
}

func (pst *RedisPresencePersistor) RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, error) {
	authStr := "auth"
	if guest {
		authStr = "guest"
	}

	var (
		presenceUsersKey       = "presence:users"
		presenceUserKey        = "presence:user:" + userID.String()
		presenceActiveGamesKey = "presence:active-games:" + userID.String()
		userKey                = userID.String() + "#" + connID.String() + "#" + username + "#" + authStr
		presenceLastSeenKey    = "presence:last-seen" + userID.String()
	)

	now := time.Now()
	expireDur := time.Minute * 2
	expiry := now.Add(expireDur)

	beforePresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch before presences")
	}

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, connChannelKey := range beforePresences {
			parts := strings.Split(connChannelKey, "#")
			cid, channel := parts[0], parts[1]
			if connID.String() == cid {
				var (
					presenceChannelKey = "presence:channel:" + channel
					userkeyChannel     = userKey + "#" + channel
				)

				if err := p.ZAdd(ctx, presenceChannelKey, redis.Z{
					Score:  float64(expiry.UnixMilli()),
					Member: userKey,
				}).Err(); err != nil {
					return fmt.Errorf("refreshpresence zadd presenceChannelKey: %w", err)
				}

				if err := p.Expire(ctx, presenceChannelKey, expireDur).Err(); err != nil {
					return fmt.Errorf("refreshpresence expire presenceChannelKey: %w", err)
				}

				if err := p.ZAdd(ctx, presenceUserKey, redis.Z{
					Score:  float64(expiry.UnixMilli()),
					Member: connChannelKey,
				}).Err(); err != nil {
					return fmt.Errorf("refreshpresence zadd presenceUserKey: %w", err)
				}

				if err := p.Expire(ctx, presenceUserKey, expireDur).Err(); err != nil {
					return fmt.Errorf("refreshpresence expire presenceUserKey: %w", err)
				}

				if err := p.ZAdd(ctx, presenceUsersKey, redis.Z{
					Score:  float64(expiry.UnixMilli()),
					Member: userkeyChannel,
				}).Err(); err != nil {
					return fmt.Errorf("refreshpresence zadd presenceUsersKey: %w", err)
				}

				if err := p.ZAdd(ctx, presenceLastSeenKey, redis.Z{
					Score:  float64(now.UnixMilli()),
					Member: userID.String(),
				}).Err(); err != nil {
					return fmt.Errorf("refreshpresence zadd presenceUsersKey: %w", err)
				}

				removeExpired := []string{
					presenceChannelKey,
					presenceUserKey,
					presenceUsersKey,
				}

				for _, key := range removeExpired {
					if err := p.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", float64(now.UnixMilli()))).Err(); err != nil {
						return fmt.Errorf("refreshpresence zremrangebyscore expired presences: %w", err)
					}
				}
			}
		}

		_ = []any{presenceActiveGamesKey}

		// @TODO: HANDLE ACTIVE GAMES

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("refreshpresence pipeline: %w", err)
	}

	afterPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch after presences")
	}

	slices.Sort(beforePresences)
	slices.Sort(afterPresences)

	return beforePresences, afterPresences, nil
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
