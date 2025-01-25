package storage

import (
	"fmt"

	"github.com/box1bs/TelegraphicVault/pkg/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgresDB(dsn string) (*Postgres, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    err = db.AutoMigrate(&model.Bookmark{}, &model.Note{}, &model.Tag{}, &model.User{})
    if err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    return &Postgres{db: db}, nil
}