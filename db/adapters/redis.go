package adapters

import (
	"github.com/go-redis/redis/v8"
)

func CreateRedisClient(addr string) *redis.Client {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	return redis
}
