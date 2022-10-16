package pkg

import (
	"context"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	postgresDSNKey = "POSTGRES"
)

func NewDB(_ context.Context) (*gorm.DB, error) {
	var (
		dsn = os.Getenv(postgresDSNKey)
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
