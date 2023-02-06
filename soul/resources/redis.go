package resources

import (
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func initRedises(env *Env) map[string]*redis.Client {
	mp := make(map[string]*redis.Client)
	if env.RedisHostRU != "" {
		log.Info().Msgf("[redis] activate redis %s", "ru")
		mp["ru"] = redis.NewClient(&redis.Options{Addr: env.RedisHostRU, Password: env.RedisPassword, DB: 0})
	}
	if env.RedisHostEU != "" {
		log.Info().Msgf("[redis] activate redis %s", "eu")
		mp["eu"] = redis.NewClient(&redis.Options{Addr: env.RedisHostEU, Password: env.RedisPassword, DB: 0})
	}
	return mp
}
