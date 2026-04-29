package redis

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
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

func (pst *RedisPresencePersistor) UserLastSeen(ctx context.Context, userID uuid.UUID) (time.Time, error) {
	lastSeenPresenceKey := "presence:user:last-seen"

	ts, err := pst.rdb.ZScore(ctx, lastSeenPresenceKey, userID.String()).Result()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get last seen: %w", err)
	}

	return time.Unix(int64(ts), 0), nil
}

func (pst *RedisPresencePersistor) ListUsersInChannel(ctx context.Context, channel string) ([]persistence.UserPresenceInfo, error) {
	connIDs, err := pst.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		ByScore: true,
		Key:     "presence:channel:conns:" + channel,
		Start:   time.Now().Unix(),
		Stop:    "+inf",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list conns in channel: %w", err)
	}

	cmds := make([]*redis.MapStringStringCmd, len(connIDs))

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for i, connID := range connIDs {
			cmds[i] = p.HGetAll(ctx, "presence:conn:"+connID)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("ListUsersInChannel pipeline failed: %w", err)
	}

	userPresenceSet := make(map[string]persistence.UserPresenceInfo)

	for _, cmd := range cmds {
		meta, err := cmd.Result()
		if err != nil {
			return nil, err
		}

		userID := meta["user_id"]
		username := meta["username"]
		authState := meta["auth_state"]

		if _, exists := userPresenceSet[userID]; exists {
			continue
		}

		userPresenceSet[userID] = persistence.UserPresenceInfo{
			ID:       userID,
			Username: username,
			Guest:    authState == "guest",
		}
	}

	userPresenceInfos := make([]persistence.UserPresenceInfo, 0, len(userPresenceSet))
	for _, info := range userPresenceSet {
		userPresenceInfos = append(userPresenceInfos, info)
	}

	return userPresenceInfos, nil
}

func (pst *RedisPresencePersistor) ListChannelsForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	connIDs, err := pst.rdb.SMembers(ctx, "presence:user:conns:"+userID.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user conns: %w", err)
	}

	cmds := make([]*redis.StringSliceCmd, len(connIDs))

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for i, connID := range connIDs {
			cmds[i] = p.SMembers(ctx, "presence:conn:channels:"+connID)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("ListChannelsForUser pipeline failed: %w", err)
	}

	channelSet := make(map[string]struct{})

	for _, cmd := range cmds {
		channels, err := cmd.Result()
		if err != nil {
			continue
		}

		for _, ch := range channels {
			channelSet[ch] = struct{}{}
		}
	}

	channels := make([]string, 0, len(channelSet))
	for ch := range channelSet {
		channels = append(channels, ch)
	}

	return channels, nil
}

func (pst *RedisPresencePersistor) UsersCountInChannel(ctx context.Context, channel string) (int64, error) {
	connIDs, err := pst.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		ByScore: true,
		Key:     "presence:channel:conns:" + channel,
		Start:   time.Now().Unix(),
		Stop:    "+inf",
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to list conns in channel: %w", err)
	}

	cmds := make([]*redis.StringCmd, len(connIDs))

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for i, connID := range connIDs {
			cmds[i] = p.HGet(ctx, "presence:conn:"+connID, "user_id")
		}
		return nil
	}); err != nil {
		return 0, fmt.Errorf("UsersCountInChannel pipeline failed: %w", err)
	}

	userSet := make(map[string]struct{})

	for _, cmd := range cmds {
		userID, err := cmd.Result()
		if err != nil {
			continue
		}

		if _, exists := userSet[userID]; exists {
			continue
		}

		userSet[userID] = struct{}{}
	}

	return int64(len(userSet)), nil
}

func (pst *RedisPresencePersistor) ConnsCountInChannel(ctx context.Context, channel string) (int64, error) {
	count, err := pst.rdb.ZCount(ctx, "presence:channel:conns:"+channel, strconv.FormatInt(time.Now().Unix(), 10), "+inf").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get conns in channel count: %w", err)
	}

	return count, nil
}

func (pst *RedisPresencePersistor) TotalActiveConnsCount(ctx context.Context) (int64, error) {
	count, err := pst.rdb.ZCount(ctx, "presence:conns", strconv.FormatInt(time.Now().Unix(), 10), "+inf").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get total active conns count: %w", err)
	}

	return count, nil
}
