package storage

import (
	"context"
	"something/model"

	"github.com/google/uuid"
)

type NoteFilter struct {
	UserID uuid.UUID 	`json:"user_id"`
	Tag    string 		`json:"tag"`
}

func (p *Postgres) CreateNote(ctx context.Context, note model.Note) error {
	return p.db.WithContext(ctx).Create(&note).Error
}

func (p *Postgres) GetNote(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note *model.Note
	err := p.db.WithContext(ctx).Where("id = ?", id).First(note).Error
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (p *Postgres) UpdateNote(ctx context.Context, note *model.Note) error {
	return p.db.WithContext(ctx).Save(note).Error
}

func (p *Postgres) DeleteNote(ctx context.Context, id uuid.UUID) error {
	return p.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Note{}).Error
}

func (p *Postgres) ListNotes(ctx context.Context, filter NoteFilter) ([]*model.Note, error) {
	var notes []*model.Note
	query := p.db.WithContext(ctx)

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Tag != "" {
		query = query.Where("? = ANY(tags)", filter.Tag)
	}

	err := query.Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}