package redis

import (
	"net"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	// host := os.Getenv("REDIS_ADDRESS")
	// port := os.Getenv("REDIS_PASSWORD")
	// db := os.Getenv("REDIS_DB")
	// pwd := os.Getenv("REDIS_PASSWORD")

	os.Getenv("REDIS_DB")
	host := "redis"
	port := "6379"
	db := "0"
	pwd := ""

	addr := net.JoinHostPort(host, port)

	dbv, err := strconv.Atoi(db)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       dbv,
	})

	return rdb, nil
}
