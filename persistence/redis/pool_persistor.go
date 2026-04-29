package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/dankobg/juicer/persistence"
	"github.com/dankobg/juicer/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var _ persistence.PoolPersistor = (*RedisPoolPersistor)(nil)

type RedisPoolPersistor struct {
	*RedisPersistor
}

func NewRedisPoolPersistor(rs *RedisPersistor) *RedisPoolPersistor {
	return &RedisPoolPersistor{
		RedisPersistor: rs,
	}
}

func (pst *RedisPoolPersistor) JoinPool(ctx context.Context, userID uuid.UUID, pool dbtype.Pool) error {
	var (
		poolKey     = fmt.Sprintf("pool:%s", pool.Name())
		poolUserKey = fmt.Sprintf("pool-user:%s", userID.String())
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

		return nil
	}); err != nil {
		return fmt.Errorf("JoinPool pipeline failed: %w", err)
	}

	return nil
}

func (pst *RedisPoolPersistor) LeavePool(ctx context.Context, userID uuid.UUID) error {
	poolUserKey := fmt.Sprintf("pool-user:%s", userID.String())

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
