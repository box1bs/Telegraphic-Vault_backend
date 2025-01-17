package auth

import (
	"something/pkg/database"
	"something/pkg/model"

	"github.com/google/uuid"
)

type UserRepository struct {
	storage.Storage
}

func (r *UserRepository) FindByID(u uuid.UUID) (*model.User, error) {
	return r.Storage.FindByID(u)
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	return r.Storage.FindByUsername(username)
}