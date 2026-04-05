package redis

import (
	"github.com/dankobg/juicer/persistence"
	"github.com/redis/go-redis/v9"
)

// var _ store.Store = (*RedisStore)(nil)

type RedisPersistor struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *RedisPersistor {
	return &RedisPersistor{rdb: rdb}
}

func (ps *RedisPersistor) Presence() persistence.PresencePersistor {
	return NewPgPresenceStore(ps)
}
