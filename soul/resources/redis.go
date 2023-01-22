package resources

import (
	"github.com/go-redis/redis/v8"
)

func initRedis(env *Env) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: env.RedisHost, Password: env.RedisPassword, DB: 0})
}
