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

func NewRedisPresencePersistor(rs *RedisPersistor) *RedisPresencePersistor {
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

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("setpresence pipeline: %w", err)
	}

	activeUserGamesPresences, err := pst.rdb.ZRange(ctx, presenceActiveGamesKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch active games presences")
	}

	afterUserPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch after presences")
	}

	oldChannels := make([]string, 0)
	newChannels := make([]string, 0)

	for _, before := range beforeUserPresences {
		channel := strings.Split(before, "#")[1]
		oldChannels = append(oldChannels, channel)
	}
	for _, after := range afterUserPresences {
		channel := strings.Split(after, "#")[1]
		newChannels = append(newChannels, channel)
	}

	for _, x := range activeUserGamesPresences {
		gameIDStr := strings.Split(x, "#")[1]
		activeGameKey := "activegame:" + gameIDStr
		oldChannels = append(oldChannels, activeGameKey)
		newChannels = append(newChannels, activeGameKey)
	}

	slices.Sort(oldChannels)
	slices.Sort(newChannels)

	return oldChannels, newChannels, nil
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
		presenceLastSeenKey    = "presence:last-seen:" + userID.String()
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

	oldChannels := make([]string, 0)
	newChannels := make([]string, 0)
	removedChannels := make([]string, 0)

	for _, before := range beforePresences {
		channel := strings.Split(before, "#")[1]
		oldChannels = append(oldChannels, channel)
	}
	for _, after := range afterPresences {
		channel := strings.Split(after, "#")[1]
		newChannels = append(newChannels, channel)
	}

	slices.Sort(oldChannels)
	slices.Sort(newChannels)

	var i, j int

	for i < len(oldChannels) && j < len(newChannels) {
		if oldChannels[i] == newChannels[j] {
			// present in both -> unchanged
			i++
			j++
		} else if oldChannels[i] < newChannels[j] {
			// in old but not in new -> removed
			removedChannels = append(removedChannels, oldChannels[i])
			i++
		} else {
			// in new but not in old -> added (optional)
			// addedChannels = append(addedChannels, newChannels[j])
			j++
		}
	}

	// anything left in old -> removed
	for i < len(oldChannels) {
		removedChannels = append(removedChannels, oldChannels[i])
		i++
	}

	// anything left in new -> added (optional)
	for j < len(newChannels) {
		// addedChannels = append(addedChannels, newChannels[j])
		j++
	}

	slices.Sort(removedChannels)

	return oldChannels, newChannels, removedChannels, nil
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
		presenceLastSeenKey    = "presence:last-seen:" + userID.String()
	)

	now := time.Now()
	expireDur := time.Minute * 2
	expiry := now.Add(expireDur)

	beforePresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch before presences")
	}

	oldChannels := make([]string, 0)
	newChannels := make([]string, 0)

	beforeActiveUserGamesPresences, err := pst.rdb.ZRange(ctx, presenceActiveGamesKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch active games presences")
	}

	for _, x := range beforeActiveUserGamesPresences {
		gameIDStr := strings.Split(x, "#")[1]
		activeGameKey := "activegame:" + gameIDStr
		oldChannels = append(oldChannels, activeGameKey)
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
					presenceActiveGamesKey,
				}

				for _, key := range removeExpired {
					if err := p.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", float64(now.UnixMilli()))).Err(); err != nil {
						return fmt.Errorf("refreshpresence zremrangebyscore expired presences: %w", err)
					}
				}
			}
		}

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("refreshpresence pipeline: %w", err)
	}

	afterPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch after presences")
	}

	for _, before := range beforePresences {
		channel := strings.Split(before, "#")[1]
		oldChannels = append(oldChannels, channel)
	}
	for _, after := range afterPresences {
		channel := strings.Split(after, "#")[1]
		newChannels = append(newChannels, channel)
	}

	afterActiveUserGamesPresences, err := pst.rdb.ZRange(ctx, presenceActiveGamesKey, 0, -1).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("fetch active games presences")
	}

	for _, x := range afterActiveUserGamesPresences {
		gameIDStr := strings.Split(x, "#")[1]
		activeGameKey := "activegame:" + gameIDStr
		newChannels = append(newChannels, activeGameKey)
	}

	slices.Sort(oldChannels)
	slices.Sort(newChannels)

	return beforePresences, afterPresences, nil
}

func (pst *RedisPresencePersistor) GetPresence(ctx context.Context, userID uuid.UUID) ([]string, error) {
	presenceUserKey := "presence:user:" + userID.String()

	userPresence, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("get user presence: %w", err)
	}

	if userPresence == nil {
		return []string{}, nil
	}

	channels := make([]string, len(userPresence))
	for i, connChannelKey := range userPresence {
		channel := strings.Split(connChannelKey, "#")[1]
		channels[i] = channel
	}

	return channels, nil
}

func (pst *RedisPresencePersistor) CountInChannel(ctx context.Context, channel string) (int64, error) {
	presenceChannelKey := "presence:channel:" + channel
	count, err := pst.rdb.ZCard(ctx, presenceChannelKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get count in channel: %w", err)
	}

	return count, nil
}

func (pst *RedisPresencePersistor) LastSeen(ctx context.Context, userID uuid.UUID) (int64, error) {
	presenceLastSeenKey := "presence:last-seen:" + userID.String()

	ts, err := pst.rdb.ZScore(ctx, presenceLastSeenKey, userID.String()).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get last seen: %w", err)
	}

	return int64(ts), nil
}

func (pst *RedisPresencePersistor) GetChannelsForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var (
		presenceUserKey        = "presence:user:" + userID.String()
		presenceActiveGamesKey = "presence:active-games:" + userID.String()
	)

	userPresences, err := pst.rdb.ZRange(ctx, presenceUserKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("fetch user presences: %w", err)
	}

	activeUserGamesPresences, err := pst.rdb.ZRange(ctx, presenceActiveGamesKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("fetch active user game presences: %w", err)
	}

	channels := make([]string, 0)

	for _, x := range userPresences {
		channel := strings.Split(x, "#")[1]
		channels = append(channels, channel)
	}

	for _, x := range activeUserGamesPresences {
		gameIDStr := strings.Split(x, "#")[1]
		activeGameKey := "activegame:" + gameIDStr
		channels = append(channels, activeGameKey)
	}

	slices.Sort(channels)

	return channels, nil
}

func (pst *RedisPresencePersistor) GetUsersInChannel(ctx context.Context, channel string) ([]persistence.UserPresenceInfo, error) {
	presenceChannelKey := "presence:channel:" + channel

	presenceChannels, err := pst.rdb.ZRange(ctx, presenceChannelKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("fetch users in channels: %w", err)
	}

	userPresenceInfos := make([]persistence.UserPresenceInfo, len(presenceChannels))

	for i, userKey := range presenceChannels {
		parts := strings.Split(userKey, "#")

		guest := true
		if parts[3] == "auth" {
			guest = false
		}
		userPresenceInfos[i] = persistence.UserPresenceInfo{
			ID:       parts[0],
			Username: parts[2],
			Guest:    guest,
		}
	}

	return userPresenceInfos, nil
}
