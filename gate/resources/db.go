package resources

import (
	"github.com/artchitector/artchitect.git/gate/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database connection
func initDB(env *Env) *gorm.DB {
	db, err := gorm.Open(postgres.Open(env.DbDSN), &gorm.Config{})

	if err := db.AutoMigrate(&model.Painting{}); err != nil {
		log.Fatal().Err(errors.Wrap(err, "failed to auto-migrate"))
	}

	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "database connection failed"))
	}

	return db
}
