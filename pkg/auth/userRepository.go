package auth

import (
	"github.com/box1bs/ClockworkChronicle/pkg/database"
	"github.com/box1bs/ClockworkChronicle/pkg/model"

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