package resources

import "gorm.io/gorm"

type Resources struct {
	db  *gorm.DB
	bus *Bus
}

func (r *Resources) GetDB() *gorm.DB {
	return r.db
}

func InitResources() *Resources {
	db := initDB()
	bus := initBus()

	return &Resources{db, bus}
}
