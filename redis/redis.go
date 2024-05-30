package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	host := "redis"
	port := "6379"
	db := 0
	pwd := ""

	addr := net.JoinHostPort(host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis server: %w", err)
	}

	return rdb, nil
}
