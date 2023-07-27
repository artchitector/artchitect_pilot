package resources

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Resources struct {
	env     *Env
	db      *gorm.DB
	redises map[string]*redis.Client
	webcam  *Webcam
}

func (r *Resources) GetDB() *gorm.DB {
	return r.db
}

func (r *Resources) GetEnv() *Env {
	return r.env
}

func (r *Resources) GetRedises() map[string]*redis.Client {
	return r.redises
}

func (r *Resources) GetWebcam() *Webcam {
	return r.webcam
}

func InitResources() *Resources {
	env := initEnv()
	db := initDB(env)
	redises := initRedises(env)

	return &Resources{env, db, redises, &Webcam{env.OriginURL}}
}
