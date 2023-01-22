package resources

import (
	"github.com/go-redis/redis/v8"
)

func initRedis(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})
}
