package config

import (
	"something/model"
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
	FindByID(uuid.UUID) (*model.User, error)
	FindByUsername(string) (*model.User, error)
	SaveRefreshToken(uuid.UUID, string, time.Time) error
	RevokeRefreshToken(uuid.UUID, string) error
}

type AuthConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}