package resources

import (
	"gorm.io/gorm"
)

type Resources struct {
	env *Env
	db  *gorm.DB
}

func (r *Resources) GetDB() *gorm.DB {
	return r.db
}

func (r *Resources) GetEnv() *Env {
	return r.env
}

func InitResources() *Resources {
	env := initEnv()
	db := initDB(env)

	return &Resources{env, db}
}
