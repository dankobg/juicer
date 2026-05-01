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

func (pst *RedisPresencePersistor) SetPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool, channel string) (persistence.PresenceChannelsDiff, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, channel, now.Unix(), expiration.Unix()}

	result, err := pst.setPresenceScript.Run(ctx, pst.rdb, keys, args...).Result()
	if err != nil {
		return persistence.PresenceChannelsDiff{}, err
	}

	arr, ok := result.([]any)
	if !ok {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid lua result: %T", result)
	}
	if len(arr) != 2 {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid lua result length: expected 2, got: %d", len(arr))
	}

	connJoinedArr, ok1 := arr[0].([]any)
	userJoinedArr, ok2 := arr[1].([]any)
	if !ok1 || !ok2 {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays")
	}

	presenceDiffInfo := persistence.PresenceChannelsDiff{
		ConnJoined: make([]string, len(connJoinedArr)),
		UserJoined: make([]string, len(userJoinedArr)),
	}

	for i, item := range connJoinedArr {
		v, ok := item.(string)
		if !ok {
			return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays not string")
		}
		presenceDiffInfo.ConnJoined[i] = v
	}

	for i, item := range userJoinedArr {
		v, ok := item.(string)
		if !ok {
			return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays not string")
		}
		presenceDiffInfo.UserJoined[i] = v
	}

	return presenceDiffInfo, nil
}

func (pst *RedisPresencePersistor) ClearPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) (persistence.PresenceChannelsDiff, error) {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, now.Unix(), expiration.Unix()}

	result, err := pst.clearPresenceScript.Run(ctx, pst.rdb, keys, args...).Result()
	if err != nil {
		return persistence.PresenceChannelsDiff{}, err
	}

	arr, ok := result.([]any)
	if !ok {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid lua result: %T", result)
	}
	if len(arr) != 2 {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid lua result length: expected 2, got: %d", len(arr))
	}

	connLeftArr, ok1 := arr[0].([]any)
	userLeftArr, ok2 := arr[1].([]any)
	if !ok1 || !ok2 {
		return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays")
	}

	presenceDiffInfo := persistence.PresenceChannelsDiff{
		ConnLeft: make([]string, len(connLeftArr)),
		UserLeft: make([]string, len(userLeftArr)),
	}

	for i, item := range connLeftArr {
		v, ok := item.(string)
		if !ok {
			return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays not string")
		}
		presenceDiffInfo.ConnLeft[i] = v
	}

	for i, item := range userLeftArr {
		v, ok := item.(string)
		if !ok {
			return persistence.PresenceChannelsDiff{}, fmt.Errorf("invalid result sub arrays not string")
		}
		presenceDiffInfo.UserLeft[i] = v
	}

	return presenceDiffInfo, nil
}

func (pst *RedisPresencePersistor) RefreshPresence(ctx context.Context, userID uuid.UUID, connID uuid.UUID, username string, guest bool) error {
	now := time.Now()
	expiration := now.Add(time.Minute * 2)

	keys := []string{}
	args := []any{userID.String(), connID.String(), username, guest, now.Unix(), expiration.Unix()}

	if _, err := pst.refreshPresenceScript.Run(ctx, pst.rdb, keys, args...).Result(); err != nil {
		return err
	}

	return nil
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
	userIDs, err := pst.rdb.SMembers(ctx, "presence:channel:users:"+channel).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list users in channel: %w", err)
	}

	cmds := make([]*redis.MapStringStringCmd, len(userIDs))

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		for i, userID := range userIDs {
			cmds[i] = p.HGetAll(ctx, "presence:user:meta:"+userID)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("pipeline failed: %w", err)
	}

	users := make([]persistence.UserPresenceInfo, 0, len(userIDs))

	for i, cmd := range cmds {
		meta, err := cmd.Result()
		if err != nil {
			return nil, err
		}

		users = append(users, persistence.UserPresenceInfo{
			ID:       userIDs[i],
			Username: meta["username"],
			Guest:    meta["auth_state"] == "guest",
		})
	}

	return users, nil
}

func (pst *RedisPresencePersistor) ListChannelsForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	channels, err := pst.rdb.SMembers(ctx, "presence:user:channels:"+userID.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user channels: %w", err)
	}

	return channels, nil
}

func (pst *RedisPresencePersistor) UsersCountInChannel(ctx context.Context, channel string) (int64, error) {
	count, err := pst.rdb.SCard(ctx, "presence:channel:users:"+channel).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count users in channel: %w", err)
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
