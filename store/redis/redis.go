package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/store"
	"github.com/redis/go-redis/v9"
)

// var _ store.Store = (*RedisStore)(nil)

type RedisStore struct {
	rdb *redis.Client
}

func Connect(config config.RedisConfig) (*redis.Client, error) {
	addr := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis server: %w", err)
	}

	return rdb, nil
}

func New(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (ps *RedisStore) Presence() store.PresenceStore {
	return NewPgPresenceStore(ps)
}
