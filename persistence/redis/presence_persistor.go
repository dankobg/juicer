package redis

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dankobg/juicer/persistence"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:embed lua-scripts/set_presence.lua
var setPresenceLuaScript string

//go:embed lua-scripts/refresh_presence.lua
var refreshPresenceLuaScript string

//go:embed lua-scripts/clear_presence.lua
var clearPresenceLuaScript string

var _ persistence.PresencePersistor = (*RedisPresencePersistor)(nil)

type RedisPresencePersistor struct {
	*RedisPersistor

	setPresenceScript     *redis.Script
	refreshPresenceScript *redis.Script
	clearPresenceScript   *redis.Script
}

func NewRedisPresencePersistor(rs *RedisPersistor) *RedisPresencePersistor {
	return &RedisPresencePersistor{
		RedisPersistor:        rs,
		setPresenceScript:     redis.NewScript(setPresenceLuaScript),
		refreshPresenceScript: redis.NewScript(refreshPresenceLuaScript),
		clearPresenceScript:   redis.NewScript(clearPresenceLuaScript),
	}
}

func (pst *RedisPresencePersistor) SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channel string) ([]string, []string, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, channel, now.Unix(), expiration.Unix()}

	_, err := pst.setPresenceScript.Run(ctx, pst.rdb, keys, args...).Result()
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func (pst *RedisPresencePersistor) ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, []string, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, now.Unix(), expiration.Unix()}

	_, err := pst.clearPresenceScript.Run(ctx, pst.rdb, keys, args...).Result()
	if err != nil {
		return nil, nil, nil, err
	}

	return nil, nil, nil, nil
}

func (pst *RedisPresencePersistor) RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) ([]string, []string, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, now.Unix(), expiration.Unix()}

	_, err := pst.refreshPresenceScript.Run(ctx, pst.rdb, keys, args...).Result()
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
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
