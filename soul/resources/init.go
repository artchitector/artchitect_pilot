package resources

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Resources struct {
	env   *Env
	db    *gorm.DB
	redis *redis.Client
}

func (r *Resources) GetDB() *gorm.DB {
	return r.db
}

func (r *Resources) GetEnv() *Env {
	return r.env
}

func (r *Resources) GetRedis() *redis.Client {
	return r.redis
}

func InitResources() *Resources {
	env := initEnv()
	db := initDB(env)
	red := initRedis(env)

	return &Resources{env, db, red}
}
