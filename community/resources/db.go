package resources

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	log2 "log"
	"os"
	"time"
)

// Database connection
func initDB(env *Env) *gorm.DB {
	pg := postgres.Open(env.DbDSN)

	gormLogger := logger.New(
		log2.New(os.Stdout, "\r\n", log2.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Warn,
		},
	)
	db, err := gorm.Open(pg, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "failed to connect to postgres"))
	}

	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "database connection failed"))
	}

	return db
}
