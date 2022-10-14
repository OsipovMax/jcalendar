package pkg

import (
	"context"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	postgresDSNKey = "POSTGRES"
)

func NewDB(_ context.Context) (*gorm.DB, error) {
	var (
		dsn = os.Getenv(postgresDSNKey)
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
