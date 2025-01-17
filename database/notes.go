package storage

import (
	"context"
	"something/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NoteFilter struct {
	UserID uuid.UUID 	`json:"user_id"`
	Tag    string 		`json:"tag"`
}

func (p *Postgres) CreateNote(ctx context.Context, note model.Note) error {
	return p.db.WithContext(ctx).Create(&note).Error
}

func (p *Postgres) GetNote(ctx context.Context, userID uuid.UUID, title string) (*model.Note, error) {
	var note model.Note
	err := p.db.WithContext(ctx).Where("user_id = ? AND title = ?", userID, title).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (p *Postgres) UpdateNote(ctx context.Context, userID uuid.UUID, currentTitle, newTitle, content string, newTagNames []string) (*model.Note, error) {
	note, err := p.GetNote(ctx, userID, currentTitle)
	if err != nil {
		return nil, err
	}

	if currentTitle != newTitle {
		var exist model.Bookmark
		if err := p.db.WithContext(ctx).Model(&model.Note{}).Where("user_id = ? AND title = ?", userID, newTitle).First(&exist).Error; err == nil {
			return nil, model.ErrAlreadyExists
		}
	}

	note.Title = newTitle
	note.Content = content

	if newTagNames != nil {
		p.AddTagToNote(ctx, note, newTagNames)
	}

	err = p.db.WithContext(ctx).Save(note).Error
	return note, err
}

func (p *Postgres) DeleteNote(ctx context.Context, userID uuid.UUID, title string) error {
	result := p.db.WithContext(ctx).Where("user_id = ? AND title = ?", userID, title).Delete(&model.Note{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *Postgres) ListNotes(ctx context.Context, filter NoteFilter) ([]*model.Note, error) {
	var notes []*model.Note
	query := p.db.WithContext(ctx)

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Tag != "" {
		query = query.Joins("JOIN note_tags ON note.id = note_tags.note_id").
			Joins("JOIN tags ON tags.id = note_tags.tag_id").
			Where("tags.name = ?", filter.Tag)
	}

	err := query.Preload("Tags").Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}