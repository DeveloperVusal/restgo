package redis

import (
	"github.com/redis/go-redis/v9"
)

func Conn(opt *redis.Options) *redis.Client {
	rdb := redis.NewClient(opt)

	return rdb
}
