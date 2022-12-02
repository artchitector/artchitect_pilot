package resources

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database connection
func initDB() *gorm.DB {
	// TODO Make .env file for config
	dsn := "host=localhost user=artchitect password=1234 port=5432 sslmode=disable TimeZone=Europe/Moscow database=artchitect"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "database connection failed"))
	}

	return db
}
