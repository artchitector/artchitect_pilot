package resources

import (
	"github.com/go-redis/redis/v8"
)

func initRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr})
}
