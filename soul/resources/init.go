package resources

import "gorm.io/gorm"

type Resources struct {
	db *gorm.DB
}

func (r *Resources) GetDB() *gorm.DB {
	return r.db
}

func InitResources() *Resources {
	db := initDB()

	return &Resources{db}
}
