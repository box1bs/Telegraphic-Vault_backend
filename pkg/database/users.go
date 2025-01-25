package storage

import (
	"github.com/box1bs/ClockworkChronicle/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (p *Postgres) LastLoginUpdate(u *model.User) error {
	return p.db.Model(u).
	Where("id = ?", u.ID).
	UpdateColumn("last_login_at", gorm.Expr("CURRENT_TIMESTAMP")).
	Error
}