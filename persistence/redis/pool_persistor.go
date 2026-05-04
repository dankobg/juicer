package redis

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:embed lua-scripts/pool.lua
var poolLuaScript string

var _ persistence.PoolPersistor = (*RedisPoolPersistor)(nil)

type RedisPoolPersistor struct {
	*RedisPersistor

	poolScript *redis.Script
}

func NewRedisPoolPersistor(rs *RedisPersistor) *RedisPoolPersistor {
	return &RedisPoolPersistor{
		RedisPersistor: rs,
		poolScript:     redis.NewScript(poolLuaScript),
	}
}

func (pst *RedisPoolPersistor) MatchPair(ctx context.Context, pool dbtype.Pool) ([2]string, error) {
	poolKey := fmt.Sprintf("pool:%s", pool.Name())

	res, err := pst.poolScript.Run(ctx, pst.rdb, []string{poolKey}).Result()
	if err != nil {
		return [2]string{}, err
	}

	arr, ok := res.([]any)
	if !ok {
		return [2]string{}, fmt.Errorf("invalid result type: %T", res)
	}

	if len(arr) != 2 {
		return [2]string{}, fmt.Errorf("not pair of 2: %d", len(arr))
	}

	out := [2]string{}

	for i, item := range arr {
		userID, ok := item.(string)
		if !ok {
			return [2]string{}, fmt.Errorf("invalid result sub item: %T", item)
		}

		out[i] = userID
	}

	return out, nil
}

func (pst *RedisPoolPersistor) JoinPool(ctx context.Context, userID uuid.UUID, pool dbtype.Pool) error {
	var (
		poolKey         = fmt.Sprintf("pool:%s", pool.Name())
		poolUserKey     = fmt.Sprintf("pool-user:%s", userID.String())
		poolUserMetaKey = fmt.Sprintf("pool-user:meta:%s", userID.String())
	)

	if _, err := pst.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZAdd(ctx, poolKey, redis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: userID.String(),
		}).Err(); err != nil {
			return fmt.Errorf("JoinPool zadd poolKey: %w", err)
		}

		if err := p.Set(ctx, poolUserKey, poolKey, time.Minute*10).Err(); err != nil {
			return fmt.Errorf("JoinPool set user: %w", err)
		}

		if err := p.HSet(ctx, poolUserMetaKey, "rating", 1500).Err(); err != nil {
			return fmt.Errorf("JoinPool set user meta: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("JoinPool pipeline failed: %w", err)
	}

	return nil
}

func (pst *RedisPoolPersistor) LeavePool(ctx context.Context, userID uuid.UUID) error {
	poolUserKey := fmt.Sprintf("pool-user:%s", userID.String())
	poolUserMetaKey := fmt.Sprintf("pool-user:meta:%s", userID.String())

	poolKey, err := pst.rdb.Get(ctx, poolUserKey).Result()
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return err
	}

	if _, err = pst.rdb.TxPipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, poolKey, userID.String()).Err(); err != nil {
			return fmt.Errorf("LeavePool zrem poolKey")
		}

		if err := p.Del(ctx, poolUserKey).Err(); err != nil {
			return fmt.Errorf("LeavePool del poolUserKey")
		}

		if err := p.Del(ctx, poolUserMetaKey).Err(); err != nil {
			return fmt.Errorf("LeavePool del poolUserMetaKey")
		}

		return nil
	}); err != nil {
		return fmt.Errorf("LeavePool failed: %w", err)
	}

	return nil
}

func (pst *RedisPoolPersistor) ListPoolPlayers(ctx context.Context, pool dbtype.Pool) ([]string, error) {
	poolKey := fmt.Sprintf("pool:%s", pool.Name())

	res, err := pst.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{Key: poolKey, ByScore: true, Start: "-inf", Stop: "+inf"}).Result()
	if err != nil {
		return nil, fmt.Errorf("ListPoolPlayers failed: %w", err)
	}

	return res, nil
}

func (pst *RedisPoolPersistor) GetPoolUserMeta(ctx context.Context, userID uuid.UUID) (map[string]any, error) {
	poolUserMetaKey := fmt.Sprintf("pool-user:meta:%s", userID.String())

	res, err := pst.rdb.HGetAll(ctx, poolUserMetaKey).Result()
	if err != nil {
		return nil, fmt.Errorf("GetPoolUserMeta failed: %w", err)
	}

	out := make(map[string]any)

	for k, v := range res {
		out[k] = v
	}

	return out, nil
}
