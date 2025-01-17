package storage

import (
	"something/pkg/model"

	"github.com/google/uuid"
)

func (p *Postgres) SaveUser(u *model.User) error {
	return p.db.Create(u).Error
}

func (p *Postgres) FindByID(id uuid.UUID) (*model.User, error) {
	var u *model.User
	if err := p.db.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (p *Postgres) FindByUsername(username string) (*model.User, error) {
	var u *model.User
	if err := p.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}